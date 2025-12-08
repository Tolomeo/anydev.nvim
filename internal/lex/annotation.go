package lex

import (
	"errors"
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type lexedAnnotations struct {
	params    map[string]lexed.FunctionArg
	overloads []lexed.FunctionOverload
}

var overloadAnnotationQuery string = `
	(documentation 
		(overload_annotation 
			(function_type) @type 
			(comment)? @documentation
		) @overload
	) 
`

func (l *lexer) lexOverloadAnnotations(dockblock []string, annotations *lexedAnnotations) error {
	for _, docLine := range dockblock {
		err := l.scratch([]string{docLine})

		if err != nil {
			return err
		}

		captures, err := l.nvim.SafeTsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: overloadAnnotationQuery})

		switch {
		case errors.Is(nvim.ErrSafeTSQueryNoMatch, err):
			// TODO: TRACE
			fmt.Printf("Skipping overload annotation '%s' containing syntax errors\n", docLine)
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

var paramAnnotationQueries = map[string]string{
	"arg": `
		(documentation
			(param_annotation
				(identifier) @name
				"?"? @optional
				[
				 (builtin_type)
				 (function_type)
				 (member_type)
				] @type
				(comment)? @documentation
			)
		)
	`,
	"vararg": `
		(documentation
			(param_annotation 
				"..." @name
				[
				 (builtin_type)
				 (function_type)
				 (member_type)
				] @type
				(comment)? @documentation
			)
		)
	`,
}

func (l *lexer) lexParamAnnotations(dockblock []string, annotations *lexedAnnotations) error {
	for _, docLine := range dockblock {
		for _, paramAnnotationQuery := range paramAnnotationQueries {
			err := l.scratch([]string{docLine})

			if err != nil {
				return err
			}

			captures, err := l.nvim.SafeTsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: paramAnnotationQuery})

			switch {
			case errors.Is(nvim.ErrSafeTSQueryNoMatch, err):
				// TODO: TRACE
				fmt.Printf("Skipping param annotation '%s' containing syntax errors\n", docLine)
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

func (l *lexer) lexAnnotations(dockblock []string) (*lexedAnnotations, error) {
	annotations := lexedAnnotations{
		params: make(map[string]lexed.FunctionArg),
	}

	err := l.lexParamAnnotations(dockblock, &annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing annotations: %w", err)
	}

	err = l.lexOverloadAnnotations(dockblock, &annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing annotations: %w", err)
	}

	return &annotations, nil
}
