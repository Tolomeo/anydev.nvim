package lex

import (
	"errors"
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/ts"
	"github.com/Tolomeo/anydev.nvim/internal/utils/cache"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

var lexedAnnotationsCache = cache.NewCache[*lexedAnnotations]()

type lexedAnnotations struct {
	Type      lexed.Symbol
	private   bool
	protected bool
	params    map[string]lexed.FunctionArg
	overloads []lexed.FunctionOverload
	generics  []lexed.FunctionGeneric
	returns   []lexed.FunctionReturn
}

var typeAnnotationQuery string = fmt.Sprintf(`
	(documentation
		(type_annotation
			%s @type
		)
	)
`, anyTypeQuery)

func (l *Lexer) lexTypeAnnotation(annotationLine string, annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: typeAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	typeValue, found := slicesx.FindFunc(*captures, func(capture ts.Capture) bool {
		return capture.Id == "type"
	})

	if !found {
		return false, fmt.Errorf("Error reading type value from type annotation in line '%s'", annotationLine)
	}

	lexedType, err := l.lexType(typeValue.Node.Text)

	if err != nil {
		return false, err
	}

	annotations.Type = lexedType
	return true, nil
}

var overloadAnnotationQuery string = fmt.Sprintf(`
	(documentation 
		(overload_annotation 
			%s @type 
			(comment)? @documentation
		) @overload
	) 
`, typeQueries["function_type"])

func (l *Lexer) lexOverloadAnnotation(annotationLine string, annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: overloadAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	overload := lexed.FunctionOverload{}

	for _, capture := range *captures {
		switch capture.Id {
		case "documentation":
			overload.Documentation = []string{capture.Node.Text}
		case "type":
			overloadType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return false, fmt.Errorf("Error lexing annotation line '%s': %w", annotationLine, err)
			}

			overloadFunction, isFunction := overloadType.(*lexed.Function)

			if !isFunction {
				return false, fmt.Errorf("Error lexing overload annotation '%s': type is not function", annotationLine)
			}

			overload.Generics = overloadFunction.Generics
			overload.Args = overloadFunction.Args
			overload.Documentation = overloadFunction.Documentation
			overload.Return = overloadFunction.Return
		}
	}

	annotations.overloads = append(annotations.overloads, overload)

	return true, nil
}

var genericAnnotationQuery string = fmt.Sprintf(`
	(documentation 
		(generic_annotation
			(identifier) @generic.name
			parent_type: 
				%s? @generic.type
		) @generic
	)
`, anyTypeQuery)

func (l *Lexer) lexGenericAnnotation(docLine string, annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: genericAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	lexedGeneric := lexed.FunctionGeneric{}

	for _, capture := range *captures {
		switch capture.Id {
		case "generic.name":
			lexedGeneric.Name = capture.Node.Text
		case "generic.type":
			genericType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return false, err
			}

			lexedGeneric.Types = append(lexedGeneric.Types, genericType)
		}
	}

	if lexedGeneric.Name == "" {
		return false, fmt.Errorf("Could not retrieve generic name for generic annotation '%s'", docLine)
	}

	annotations.generics = append(annotations.generics, lexedGeneric)

	return true, nil
}

var paramAnnotationQueries = map[string]string{
	"arg": fmt.Sprintf(`
		(documentation
			(param_annotation
				(identifier) @name
				"?"? @optional
				%s @type
				(comment)? @documentation
			)
		)
	`, anyTypeQuery),
	"vararg": fmt.Sprintf(`
		(documentation
			(param_annotation 
				"..." @name
				%s @type
				(comment)? @documentation
			)
		)
	`, anyTypeQuery),
}

func (l *Lexer) lexParamAnnotation(docLine string, annotations *lexedAnnotations) (bool, error) {
	captures, hasCaptures, err := slicesx.MapFindFunc(
		mapx.Values(paramAnnotationQueries),
		func(paramAnnotationQuery string) (*[]ts.Capture, bool, error) {
			captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: paramAnnotationQuery})

			switch {
			case err != nil:
				return captures, false, err
			case captures == nil:
				return nil, false, nil
			}

			return captures, true, nil
		},
	)

	switch {
	case err != nil:
		return false, err
	case !hasCaptures:
		return false, nil
	}

	lexedParam := lexed.FunctionArg{}

	for _, capture := range *captures {
		switch capture.Id {
		case "name":
			lexedParam.Name = capture.Node.Text
		case "optional":
			lexedParam.Optional = true
		case "documentation":
			lexedParam.Documentation = []string{capture.Node.Text}
		case "type":
			if generic, isGeneric := slicesx.FindFunc(annotations.generics, func(generic lexed.FunctionGeneric) bool {
				return generic.Name == capture.Node.Text
			}); isGeneric {
				lexedParam.Type = newReferenceType(generic.Name)
				continue
			}

			lexedParamType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return false, fmt.Errorf("Error lexing annotation line '%s': %w", docLine, err)
			}

			lexedParam.Type = lexedParamType
		}
	}

	if lexedParam.Name == "" {
		return false, fmt.Errorf("Could not retrieve param name for param annotation '%s'", docLine)
	}

	annotations.params[lexedParam.Name] = lexedParam

	return true, nil
}

var returnAnnotationQuery string = fmt.Sprintf(`
	(documentation
		(return_annotation
			%s @return.type
			(comment)? @return.documentation
		) @return
	)
`, anyTypeQuery)

func (l *Lexer) lexReturnAnnotation(_ string, annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: returnAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	lexedFunctionReturn := lexed.FunctionReturn{}

	for _, capture := range *captures {
		switch capture.Id {
		case "return.name":
			lexedFunctionReturn.Name = &capture.Node.Text
		case "return.documentation":
			lexedFunctionReturn.Documentation = []string{capture.Node.Text}
		case "return.type":
			if generic, isGeneric := slicesx.FindFunc(annotations.generics, func(generic lexed.FunctionGeneric) bool {
				return generic.Name == capture.Node.Text
			}); isGeneric {
				lexedFunctionReturn.Type = newReferenceType(generic.Name)
				continue
			}

			lexedReturnType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return false, err
			}

			lexedFunctionReturn.Type = lexedReturnType
		}
	}

	annotations.returns = append(annotations.returns, lexedFunctionReturn)

	return true, nil
}

var privateAnnotationQuery string = `
	(documentation 
		(qualifier_annotation) @qualifier
		(#match? @qualifier "\\@private")
	)
`

func (l *Lexer) lexPrivateAnnotation(_ string, annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: privateAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	annotations.private = true

	return true, nil
}

var protectedAnnotationQuery string = `
	(documentation 
		(qualifier_annotation) @qualifier
		(#match? @qualifier "\\@protected")
	)
`

func (l *Lexer) lexProtectedAnnotation(_ string, annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: protectedAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	annotations.protected = true

	return true, nil
}

var annotationQuery string = `(documentation) @annotation`

func (l *Lexer) lexAnnotations(dockblock []string) (*lexedAnnotations, error) {
	if cachedAnnotations, cached := lexedAnnotationsCache.Get(dockblock...); cached {
		return cachedAnnotations, nil
	}

	annotationLines, err := slicesx.FilterFunc(dockblock, func(docLine string) (bool, error) {
		err := l.scratch([]string{docLine})

		if err != nil {
			return false, err
		}

		captures, err := l.context.nvim.SafeTsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: annotationQuery})

		switch {
		case errors.Is(nvim.ErrSafeTSQueryNoMatch, err):
			l.context.logger.Error(fmt.Sprintf("Skipping annotation line '%s' containing syntax errors", docLine))
			return false, nil
		case err != nil:
			return false, err
		case captures == nil:
			return false, nil
		}

		return true, nil
	})

	if err != nil {
		return nil, err
	}

	annotations := lexedAnnotations{
		params:    make(map[string]lexed.FunctionArg),
		generics:  []lexed.FunctionGeneric{},
		overloads: []lexed.FunctionOverload{},
		returns:   []lexed.FunctionReturn{},
	}

	// Generics are lexed ahead of other annotations, which could make use of them
	for _, annotationLine := range annotationLines {
		err := l.scratch([]string{annotationLine})

		if err != nil {
			return nil, err
		}

		_, err = l.lexGenericAnnotation(annotationLine, &annotations)

		if err != nil {
			return nil, fmt.Errorf("Error lexing generic annotations: %w", err)
		}
	}

	for _, annotationLine := range annotationLines {
		err := l.scratch([]string{annotationLine})

		if err != nil {
			return nil, err
		}

		matched, err := l.lexPrivateAnnotation(annotationLine, &annotations)

		switch {
		case err != nil:
			return nil, fmt.Errorf("Error lexing private annotation: %w", err)
		case matched:
			continue
		}

		matched, err = l.lexProtectedAnnotation(annotationLine, &annotations)

		switch {
		case err != nil:
			return nil, fmt.Errorf("Error lexing private annotation: %w", err)
		case matched:
			continue
		}

		matched, err = l.lexParamAnnotation(annotationLine, &annotations)

		switch {
		case err != nil:
			return nil, fmt.Errorf("Error lexing param annotation: %w", err)
		case matched:
			continue
		}

		matched, err = l.lexOverloadAnnotation(annotationLine, &annotations)

		switch {
		case err != nil:
			return nil, fmt.Errorf("Error lexing overload annotation: %w", err)
		case matched:
			continue
		}

		matched, err = l.lexReturnAnnotation(annotationLine, &annotations)

		switch {
		case err != nil:
			return nil, fmt.Errorf("Error lexing return annotation: %w", err)
		case matched:
			continue
		}

		matched, err = l.lexTypeAnnotation(annotationLine, &annotations)

		switch {
		case err != nil:
			return nil, fmt.Errorf("Error lexing type annotation: %w", err)
		case matched:
			continue
		}
	}

	lexedAnnotationsCache.Set(&annotations, dockblock...)
	return &annotations, nil
}
