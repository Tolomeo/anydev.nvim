package transform

import (
	"fmt"
	"strings"

	// "github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	// "github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

type TypeAnnotation struct {
	Name string
}

type AtGenericAnnotation struct {
	Name string
	Type *TypeAnnotation
}

func (tr *Transformer) getAtGenericAnnotations(buffer *nvim.ScratchBuffer) ([]AtGenericAnnotation, error) {
	atGenericAnnotations := []AtGenericAnnotation{}
	matches, err := buffer.TsQueryAll(annotation.AtGenericsQuery)

	if err != nil {
		return nil, err
	}

	if matches == nil {
		return atGenericAnnotations, nil
	}

	for _, matchCaptures := range *matches {
		atGenenericsAnnotation := annotation.NewGenerics(matchCaptures)

		for _, generic := range atGenenericsAnnotation.Generics() {
			lexedGeneric := AtGenericAnnotation{
				Name: generic.Name(),
			}

			if generic.Type() != nil {
				lexedGeneric.Type = &TypeAnnotation{*generic.Type()}
			}

			atGenericAnnotations = append(atGenericAnnotations, lexedGeneric)
		}

	}

	return atGenericAnnotations, nil
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
