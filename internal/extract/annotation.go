package extract

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type lexedParamAnnotation struct {
	Name          string
	Type          string
	Optional      bool
	Documentation []string
}

type lexedAnnotations struct {
	params map[string]lexedParamAnnotation
}

var paramAnnotationQueries = map[string]string{
	"arg": `
		(documentation
			(param_annotation
				(identifier) @name
				"?"? @optional
				(builtin_type) @type
				(comment) @documentation
			)
		)
	`,
	"vararg": `
		(documentation
			(param_annotation 
					"..." @name
					(builtin_type) @type
					(comment) @documentation
			))`,
}

func (e *extractor) lexParamAnnotations(dockblock []string, annotations *lexedAnnotations) error {
	for _, documentationLine := range dockblock {
		for _, paramAnnotationQuery := range paramAnnotationQueries {
			err := e.scratch([]string{documentationLine})

			if err != nil {
				return err
			}

			captures, err := e.nvim.TsQuery("luadoc", paramAnnotationQuery)

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
				case "type":
					lexedParam.Type = strings.Join(capture.Node.Text, "")
				case "optional":
					lexedParam.Optional = true
				case "documentation":
					lexedParam.Documentation = capture.Node.Text
				}
			}

			if lexedParam.Name == "" {
				return fmt.Errorf("Could not retrieve param name for param annotation '%s'", documentationLine)
			}

			annotations.params[lexedParam.Name] = lexedParam
			break
		}
	}

	return nil
}

func (e *extractor) lexAnnotations(dockblock []string) (*lexedAnnotations, error) {
	annotations := lexedAnnotations{
		params: make(map[string]lexedParamAnnotation),
	}

	err := e.lexParamAnnotations(dockblock, &annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing annotations: %w", err)
	}

	return &annotations, nil
}
