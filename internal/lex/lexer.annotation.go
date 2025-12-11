package lex

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/ts"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type lexedAnnotations struct {
	private   bool
	protected bool
	params    map[string]lexed.FunctionArg
	overloads []lexed.FunctionOverload
	generics  []lexed.FunctionGeneric
	returns   []lexed.FunctionReturn
}

var typeQueries = map[string]string{
	"builtin_type":         "(builtin_type)",
	"identifier":           "(identifier)",
	"array_type":           "(array_type)",
	"table_type":           "(table_type)",
	"table_literal_type":   "(table_literal_type)",
	"union_type":           "(union_type)",
	"parenthesized_type":   "(parenthesized_type)",
	"tuple_type":           "(tuple_type)",
	"function_type":        "(function_type)",
	"member_type":          "(member_type)",
	"optional_type":        "(optional_type)",
	"literal_type":         "(literal_type)",
	"numeric_literal_type": "(numeric_literal_type)",
	"custom_type":          "(custom_type)",
}

var anyTypeQuery = fmt.Sprintf(`[%s]`, strings.Join(mapx.Values(typeQueries), " "))

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
	case errors.Is(nvim.ErrTSQueryNoMatch, err):
		return false, nil
	case err != nil:
		return false, err
	}

	overload := lexed.FunctionOverload{}

	for _, capture := range captures {
		switch capture.Id {
		case "documentation":
			overload.Documentation = []string{capture.Node.Text}
		case "type":
			overloadType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return false, fmt.Errorf("Error lexing annotation line '%s': %w", annotationLine, err)
			}

			overloadFunction, isFunction := overloadType.(lexed.Function)

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
	case errors.Is(nvim.ErrTSQueryNoMatch, err):
		return false, nil
	case err != nil:
		return false, err
	}

	lexedGeneric := lexed.FunctionGeneric{}

	for _, capture := range captures {
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
		func(paramAnnotationQuery string) ([]ts.Capture, bool, error) {
			captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: paramAnnotationQuery})

			switch {
			case errors.Is(nvim.ErrTSQueryNoMatch, err):
				return captures, false, nil
			case err != nil:
				return captures, false, err
			}

			return captures, true, nil
		},
	)

	if err != nil {
		return false, err
	}

	if !hasCaptures {
		return false, nil
	}

	lexedParam := lexed.FunctionArg{}

	for _, capture := range captures {
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
	case errors.Is(nvim.ErrTSQueryNoMatch, err):
		return false, nil
	case err != nil:
		return false, err
	}

	lexedFunctionReturn := lexed.FunctionReturn{}

	for _, capture := range captures {
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
	_, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: privateAnnotationQuery})

	switch {
	case errors.Is(nvim.ErrTSQueryNoMatch, err):
		return false, nil
	case err != nil:
		return false, err
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
	_, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: protectedAnnotationQuery})

	switch {
	case errors.Is(nvim.ErrTSQueryNoMatch, err):
		return false, nil
	case err != nil:
		return false, err
	}

	annotations.protected = true

	return true, nil
}

var annotationQuery string = `(documentation) @annotation`

func (l *Lexer) lexAnnotations(dockblock []string) (*lexedAnnotations, error) {
	annotationLines, err := slicesx.FilterFunc(dockblock, func(docLine string) (bool, error) {
		err := l.scratch([]string{docLine})

		if err != nil {
			return false, err
		}

		_, err = l.context.nvim.SafeTsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: annotationQuery})

		switch {
		case errors.Is(nvim.ErrSafeTSQueryNoMatch, err):
			l.context.logger.Error(fmt.Sprintf("Skipping annotation line '%s' containing syntax errors", docLine))
			return false, nil
		case errors.Is(nvim.ErrTSQueryNoMatch, err):
			return false, nil
		case err != nil:
			return false, err
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
	}

	return &annotations, nil
}
