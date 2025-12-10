package lex

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type lexedAnnotations struct {
	params    map[string]lexed.FunctionArg
	overloads []lexed.FunctionOverload
	generics  []lexed.FunctionGeneric
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

func (l *Lexer) lexOverloadAnnotations(dockblock []string, annotations *lexedAnnotations) error {
	for _, docLine := range dockblock {
		err := l.scratch([]string{docLine})

		if err != nil {
			return err
		}

		captures, err := l.context.nvim.SafeTsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: overloadAnnotationQuery})

		switch {
		case errors.Is(nvim.ErrSafeTSQueryNoMatch, err):
			l.context.logger.Warn(fmt.Sprintf("Skipping overload annotation '%s' containing syntax errors", docLine))
			continue
		case errors.Is(nvim.ErrTSQueryNoMatch, err):
			continue
		case err != nil:
			return err
		}

		overload := lexed.FunctionOverload{}

		for _, capture := range captures {
			switch capture.Id {
			case "documentation":
				overload.Documentation = []string{capture.Node.Text}
			case "type":
				overloadType, err := l.lexType(capture.Node.Text)

				if err != nil {
					return fmt.Errorf("Error lexing annotation line '%s': %w", docLine, err)
				}

				overloadFunction, isFunction := overloadType.(lexed.Function)

				if !isFunction {
					return fmt.Errorf("Error lexing overload annotation '%s': type is not function", docLine)
				}

				overload.Generics = overloadFunction.Generics
				overload.Args = overloadFunction.Args
				overload.Documentation = overloadFunction.Documentation
				overload.Return = overloadFunction.Return
			}
		}

		annotations.overloads = append(annotations.overloads, overload)
	}

	return nil
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

func (l *Lexer) lexGenericAnnotations(dockblock []string, annotations *lexedAnnotations) error {
	for _, docLine := range dockblock {
		err := l.scratch([]string{docLine})

		captures, err := l.context.nvim.SafeTsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: genericAnnotationQuery})

		switch {
		case errors.Is(nvim.ErrSafeTSQueryNoMatch, err):
			l.context.logger.Warn(fmt.Sprintf("Skipping generic annotation '%s' containing syntax errors", docLine))
			continue
		case errors.Is(nvim.ErrTSQueryNoMatch, err):
			continue
		case err != nil:
			return err
		}

		lexedGeneric := lexed.FunctionGeneric{}

		for _, capture := range captures {
			switch capture.Id {
			case "generic.name":
				lexedGeneric.Name = capture.Node.Text
			case "generic.type":
				genericType, err := l.lexType(capture.Node.Text)

				if err != nil {
					return err
				}

				lexedGeneric.Types = append(lexedGeneric.Types, genericType)
			}
		}

		if lexedGeneric.Name == "" {
			return fmt.Errorf("Could not retrieve generic name for generic annotation '%s'", docLine)
		}

		annotations.generics = append(annotations.generics, lexedGeneric)
	}

	return nil
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

func (l *Lexer) lexParamAnnotations(dockblock []string, annotations *lexedAnnotations) error {
	for _, docLine := range dockblock {

		err := l.scratch([]string{docLine})

		if err != nil {
			return err
		}

		for _, paramAnnotationQuery := range paramAnnotationQueries {
			captures, err := l.context.nvim.SafeTsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: paramAnnotationQuery})

			switch {
			case errors.Is(nvim.ErrSafeTSQueryNoMatch, err):
				l.context.logger.Warn(fmt.Sprintf("Skipping param annotation '%s' containing syntax errors", docLine))
				continue
			case errors.Is(nvim.ErrTSQueryNoMatch, err):
				continue
			case err != nil:
				return err
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
						return fmt.Errorf("Error lexing annotation line '%s': %w", docLine, err)
					}

					lexedParam.Type = lexedParamType
				}
			}

			if lexedParam.Name == "" {
				return fmt.Errorf("Could not retrieve param name for param annotation '%s'", docLine)
			}

			annotations.params[lexedParam.Name] = lexedParam
			break
		}
	}

	return nil
}

func (l *Lexer) lexAnnotations(dockblock []string) (*lexedAnnotations, error) {
	annotations := lexedAnnotations{
		params:   make(map[string]lexed.FunctionArg),
		generics: []lexed.FunctionGeneric{},
	}

	// The order matters: generics need to be lexed first because other annotations rely on them

	err := l.lexGenericAnnotations(dockblock, &annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing annotations: %w", err)
	}

	err = l.lexParamAnnotations(dockblock, &annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing annotations: %w", err)
	}

	err = l.lexOverloadAnnotations(dockblock, &annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing annotations: %w", err)
	}

	return &annotations, nil
}
