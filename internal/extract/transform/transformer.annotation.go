package transform

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

type TypeAnnotation struct {
	Name string
}

type AtTypeAnnotation struct {
	Types         []TypeAnnotation
	Documentation []string
}

var atTypeAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(type_annotation
			"@type" . (%s) @type.type
			.
			("," . (%s) @type.type)*
			.
			(comment)? @type.documentation
			.
		) @type
	)`, anyTypeAnnotationQuery, anyTypeAnnotationQuery),
}

func (tr *Transformer) getAtTypeAnnotations(buffer *nvim.ScratchBuffer) (*AtTypeAnnotation, error) {
	captures, err := buffer.TsQueryOne(atTypeAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if captures == nil {
		return nil, nil
	}

	atType := AtTypeAnnotation{}

	for _, capture := range *captures {
		switch capture.Id {
		case "type.type":
			atType.Types = append(atType.Types, TypeAnnotation{capture.Node.Text})
		case "type.documentation":
			atType.Documentation = strings.Split(capture.Node.Text, "\n")
		}
	}

	return &atType, nil
}

type AtOverloadAnnotation struct {
	Type          TypeAnnotation
	Documentation []string
}

var atOverloadAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation 
		(overload_annotation 
			"@overload"
			.
			(%s) @type 
			.
			(comment)? @documentation
			.
		)
	)
`, typeAnnotationQueries["function_type"]),
}

func (tr *Transformer) getAtOverloadAnnotations(buffer *nvim.ScratchBuffer) ([]AtOverloadAnnotation, error) {
	atOverloads := []AtOverloadAnnotation{}
	matches, err := buffer.SafeTsQueryAll(atOverloadAnnotationQuery)

	// fmt.Printf("\n Overload matches: %+v\n", matches)

	if err != nil {
		return atOverloads, err
	}

	if matches == nil {
		return atOverloads, nil
	}

	for _, match := range *matches {
		if match.HasError {
			tr.context.Logger().Warn(fmt.Sprintf("Skipping overload annotation in '%s' because it contains syntax errors", tr.context.Target().Name()))
			continue
		}

		overload := AtOverloadAnnotation{}

		for _, capture := range match.Captures {
			switch capture.Id {
			case "documentation":
				overload.Documentation = []string{capture.Node.Text}
			case "type":
				overload.Type = TypeAnnotation{capture.Node.Text}
			}
		}

		atOverloads = append(atOverloads, overload)
	}

	return atOverloads, nil
}

type AtGenericAnnotation struct {
	Name  string
	Types []TypeAnnotation
}

var atGenericAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(generic_annotation
			"@generic"
			.
			(identifier) @generic.name
			.
			(":"
				.
				parent_type:
					(%s) @generic.type
			)?
			.
			(","
				.
				(identifier) @generic.name
				.
				(":"
					.
					parent_type:
						(%s) @generic.type
				)?
			)*
		) @generic
	)`, anyTypeAnnotationQuery, anyTypeAnnotationQuery),
}

func (tr *Transformer) getAtGenericAnnotations(buffer *nvim.ScratchBuffer) ([]AtGenericAnnotation, error) {
	atGenericAnnotations := []AtGenericAnnotation{}
	matches, err := buffer.TsQueryAll(atGenericAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if matches == nil {
		return atGenericAnnotations, nil
	}

	for _, matchCaptures := range *matches {
		lexedGeneric := AtGenericAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "generic.name":
				lexedGeneric.Name = capture.Node.Text
			case "generic.type":
				lexedGeneric.Types = append(lexedGeneric.Types, TypeAnnotation{capture.Node.Text})
			}
		}

		if lexedGeneric.Name == "" {
			return atGenericAnnotations, fmt.Errorf("Could not retrieve generic name for generic annotation '%v'", lexedGeneric)
		}

		atGenericAnnotations = append(atGenericAnnotations, lexedGeneric)
	}

	return atGenericAnnotations, nil
}

type AtParamAnnotation struct {
	Name          string
	Type          TypeAnnotation
	Optional      bool
	Documentation []string
}

var atParamAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(param_annotation
			"@param"
			.
			([("...") (identifier)]) @name
			.
			"?"? @optional
			.
			(%s) @type
			.
			(comment)? @documentation
			.
		) @param
	)`, anyTypeAnnotationQuery),
}

func (tr *Transformer) getAtParamAnnotations(buffer *nvim.ScratchBuffer) (map[string]AtParamAnnotation, error) {
	params := map[string]AtParamAnnotation{}
	matches, err := buffer.TsQueryAll(atParamAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if matches == nil {
		return params, nil
	}

	for _, matchCaptures := range *matches {
		lexedParam := AtParamAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "name":
				lexedParam.Name = capture.Node.Text
			case "optional":
				lexedParam.Optional = true
			case "documentation":
				lexedParam.Documentation = []string{capture.Node.Text}
			case "type":
				lexedParam.Type = TypeAnnotation{capture.Node.Text}
			}
		}

		if lexedParam.Name == "" {
			return params, fmt.Errorf("Could not retrieve param name for param annotation")
		}

		params[lexedParam.Name] = lexedParam
	}

	return params, nil
}

var typeAnnotationQueries = map[string]string{
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

var anyTypeAnnotationQuery = fmt.Sprintf(`[%s]`, strings.Join(mapx.Values(typeAnnotationQueries), " "))
