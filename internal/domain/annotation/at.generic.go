package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AtGenericsQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
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
	) @generic`, AnyTypeQuery, AnyTypeQuery),
}

type atGeneric struct {
	name  string
	type_ *string
}

func (ag *atGeneric) Name() string {
	return ag.name
}

func (ag *atGeneric) Type() *string {
	return ag.type_
}

type AtGenerics struct {
	generics []atGeneric
}

func (ag *AtGenerics) Generics() []atGeneric {
	return ag.generics
}

func NewGenerics(captures nvim.TsQueryMatch) *AtGenerics {
	atGenerics := AtGenerics{}

	for _, capture := range captures {
		switch capture.Id {
		case "generic.name":
			atGenerics.generics = append(atGenerics.generics, atGeneric{name: capture.Node.Text})
		case "generic.type":
			atGenerics.generics[len(atGenerics.generics)-1].type_ = &capture.Node.Text
		}
	}

	return &atGenerics
}
