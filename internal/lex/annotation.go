package lex

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type lexedParamAnnotation struct {
	Name          string
	Type          lexed.Symbol
	Optional      bool
	Documentation []string
}

type lexedOverloadAnnotation struct {
	Type          lexed.Function
	Documentation []string
}

type lexedAnnotations struct {
	params    map[string]lexedParamAnnotation
	overloads []lexedOverloadAnnotation
}

var overloadAnnotationQuery string = `
	(documentation 
		(overload_annotation 
			(function_type) @type 
			(comment)? @documentation
		)
	) 
`

func (l *lexer) lexOverloadAnnotations(dockblock []string, annotations *lexedAnnotations) error {
	for _, docLine := range dockblock {
		err := l.scratch([]string{docLine})

		if err != nil {
			return err
		}
		captures, err := l.nvim.TsQuery("luadoc", overloadAnnotationQuery)

		switch {
		case errors.Is(nvim.ErrNotFound, err):
			continue
		case err != nil:
			return err
		}

		lexedOverload := lexedOverloadAnnotation{}

		for _, capture := range captures {
			switch capture.Id {
			case "documentation":
				lexedOverload.Documentation = capture.Node.Text
			case "type":
				lexedOverloadType, err := l.lexType(strings.Join(capture.Node.Text, ""))

				if err != nil {
					return fmt.Errorf("Error lexing annotation line '%s': %w", docLine, err)
				}

				lexedOverloadFunctionType, isFunctionType := lexedOverloadType.(lexed.Function)

				if !isFunctionType {
					return fmt.Errorf("Error lexing overload annotation '%s': type is not function", docLine)
				}

				fmt.Printf("%v", lexedOverloadFunctionType)
				// lexedOverload.Type = lexedOverloadFunctionType
			}
		}

		annotations.overloads = append(annotations.overloads, lexedOverload)
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
			))`,
}

func (l *lexer) lexParamAnnotations(dockblock []string, annotations *lexedAnnotations) error {
	for _, docLine := range dockblock {
		for _, paramAnnotationQuery := range paramAnnotationQueries {
			err := l.scratch([]string{docLine})

			if err != nil {
				return err
			}

			captures, err := l.nvim.TsQuery("luadoc", paramAnnotationQuery)

			switch {
			case errors.Is(nvim.ErrNotFound, err):
				continue
			case err != nil:
				return err
			}

			lexedParam := lexedParamAnnotation{}

			for _, capture := range captures {
				switch capture.Id {
				case "name":
					lexedParam.Name = strings.Join(capture.Node.Text, "")
				case "optional":
					lexedParam.Optional = true
				case "documentation":
					lexedParam.Documentation = capture.Node.Text
				case "type":
					lexedParamType, err := l.lexType(strings.Join(capture.Node.Text, ""))

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
		params: make(map[string]lexedParamAnnotation),
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
