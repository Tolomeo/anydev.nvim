package lex

import (
	"fmt"

	// "github.com/Tolomeo/anydev.nvim/internal/crawl"
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
	) @typeannotation
`, anyTypeQuery)

func (l *Lexer) lexTypeAnnotations(annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: typeAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	types := []string{}

	for _, capture := range *captures {
		switch capture.Id {
		case "type":
			types = append(types, capture.Node.Text)
		}
	}

	switch {
	case len(types) > 1:
		return false, fmt.Errorf("Error reading type annotation: too many annotations received (%d), expected 1", len(types))
	case len(types) < 1:
		return false, fmt.Errorf("Error reading type annotation: type value not found")
	}

	lexedType, err := l.lexType(types[0])

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
		)
	) @overload
`, typeQueries["function_type"])

func (l *Lexer) lexOverloadAnnotations(annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: overloadAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	overloads := []lexed.FunctionOverload{}

	for _, capture := range *captures {
		switch capture.Id {
		case "overload":
			overloads = append(overloads, lexed.FunctionOverload{})
		case "documentation":
			overloads[len(overloads)-1].Documentation = []string{capture.Node.Text}
		case "type":
			overloadType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return false, fmt.Errorf("Error overload annotation type: %w", err)
			}

			overloadFunction, isFunction := overloadType.(*lexed.Function)

			if !isFunction {
				return false, fmt.Errorf("Error lexing overload annotation type: lexed type '%+v' is not a function", overloadFunction)
			}

			overloads[len(overloads)-1].Generics = overloadFunction.Generics
			overloads[len(overloads)-1].Args = overloadFunction.Args
			overloads[len(overloads)-1].Documentation = overloadFunction.Documentation
			overloads[len(overloads)-1].Return = overloadFunction.Return
		}
	}

	annotations.overloads = overloads

	return true, nil
}

var genericAnnotationQuery string = fmt.Sprintf(`
	(documentation 
		(generic_annotation
			(identifier) @generic.name
			parent_type: 
				%s? @generic.type
		)
	) @generic
`, anyTypeQuery)

func (l *Lexer) lexGenericAnnotations(annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: genericAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	lexedGenerics := []lexed.FunctionGeneric{}

	for _, capture := range *captures {
		switch capture.Id {
		case "generic":
			lexedGenerics = append(lexedGenerics, lexed.FunctionGeneric{})
		case "generic.name":
			lexedGenerics[len(lexedGenerics)-1].Name = capture.Node.Text
		case "generic.type":
			genericType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return false, err
			}

			lexedGenerics[len(lexedGenerics)-1].Types = append(lexedGenerics[len(lexedGenerics)-1].Types, genericType)
		}
	}

	for _, lexedGeneric := range lexedGenerics {
		if lexedGeneric.Name == "" {
			return false, fmt.Errorf("Could not retrieve generic name for generic annotation '%v'", lexedGeneric)
		}

		annotations.generics = append(annotations.generics, lexedGeneric)
	}

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
		) @param
	`, anyTypeQuery),
	"vararg": fmt.Sprintf(`
		(documentation
			(param_annotation 
				"..." @name
				%s @type
				(comment)? @documentation
			)
		) @param
	`, anyTypeQuery),
}

func (l *Lexer) lexParamAnnotations(annotations *lexedAnnotations) (bool, error) {
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

	lexedParams := []lexed.FunctionArg{}

	for _, capture := range *captures {
		switch capture.Id {
		case "param":
			lexedParams = append(lexedParams, lexed.FunctionArg{})
		case "name":
			lexedParams[len(lexedParams)-1].Name = capture.Node.Text
		case "optional":
			lexedParams[len(lexedParams)-1].Optional = true
		case "documentation":
			lexedParams[len(lexedParams)-1].Documentation = []string{capture.Node.Text}
		case "type":
			if generic, isGeneric := slicesx.FindFunc(annotations.generics, func(generic lexed.FunctionGeneric) bool {
				return generic.Name == capture.Node.Text
			}); isGeneric {
				lexedParams[len(lexedParams)-1].Type = newReferenceType(generic.Name)
				continue
			}

			lexedParamType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return false, fmt.Errorf("Error lexing type annotations : %w", err)
			}

			lexedParams[len(lexedParams)-1].Type = lexedParamType
		}
	}

	for _, lexedParam := range lexedParams {
		if lexedParam.Name == "" {
			return false, fmt.Errorf("Could not retrieve param name for param annotation")
		}

		annotations.params[lexedParam.Name] = lexedParam
	}

	return true, nil
}

var returnAnnotationQuery string = fmt.Sprintf(`
	(documentation
		(return_annotation
			%s @return.type
			(comment)? @return.documentation
		)
	) @return
`, anyTypeQuery)

func (l *Lexer) lexReturnAnnotations(annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: returnAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	lexedFunctionReturns := []lexed.FunctionReturn{}

	for _, capture := range *captures {
		switch capture.Id {
		case "return":
			lexedFunctionReturns = append(lexedFunctionReturns, lexed.FunctionReturn{})
		case "return.name":
			lexedFunctionReturns[len(lexedFunctionReturns)-1].Name = &capture.Node.Text
		case "return.documentation":
			lexedFunctionReturns[len(lexedFunctionReturns)-1].Documentation = []string{capture.Node.Text}
		case "return.type":
			if generic, isGeneric := slicesx.FindFunc(annotations.generics, func(generic lexed.FunctionGeneric) bool {
				return generic.Name == capture.Node.Text
			}); isGeneric {
				lexedFunctionReturns[len(lexedFunctionReturns)-1].Type = newReferenceType(generic.Name)
				continue
			}

			lexedReturnType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return false, err
			}

			lexedFunctionReturns[len(lexedFunctionReturns)-1].Type = lexedReturnType
		}
	}

	annotations.returns = lexedFunctionReturns

	return true, nil
}

var privateAnnotationQuery string = `
	(documentation 
		(qualifier_annotation) @qualifier
		(#match? @qualifier "\\@private")
	)
`

func (l *Lexer) lexPrivateAnnotation(annotations *lexedAnnotations) (bool, error) {
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

func (l *Lexer) lexProtectedAnnotation(annotations *lexedAnnotations) (bool, error) {
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

func (l *Lexer) lexAnnotations(dockblock []string) (*lexedAnnotations, error) {
	if cachedAnnotations, cached := lexedAnnotationsCache.Get(dockblock...); cached {
		return cachedAnnotations, nil
	}

	/* annotationLines, err := slicesx.FilterFunc(dockblock, func(docLine string) (bool, error) {
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
	} */

	annotations := lexedAnnotations{
		params:    make(map[string]lexed.FunctionArg),
		generics:  []lexed.FunctionGeneric{},
		overloads: []lexed.FunctionOverload{},
		returns:   []lexed.FunctionReturn{},
	}

	err := l.scratch(dockblock)

	if err != nil {
		return nil, err
	}

	// Generics are lexed ahead of other annotations, which could read them
	_, err = l.lexGenericAnnotations(&annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing generic annotations: %w", err)
	}

	_, err = l.lexPrivateAnnotation(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing private annotation: %w", err)
	}

	_, err = l.lexProtectedAnnotation(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing private annotation: %w", err)
	}

	_, err = l.lexParamAnnotations(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing param annotation: %w", err)
	}

	_, err = l.lexOverloadAnnotations(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing overload annotation: %w", err)
	}

	_, err = l.lexReturnAnnotations(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing return annotation: %w", err)
	}

	_, err = l.lexTypeAnnotations(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type annotation: %w", err)
	}

	lexedAnnotationsCache.Set(&annotations, dockblock...)
	return &annotations, nil
}
