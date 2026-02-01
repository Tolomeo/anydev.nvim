package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AtOverloadQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
		(overload_annotation 
			"@overload"
			.
			(%s) @type 
			.
			(comment)? @documentation
			.
		)`, TypeQueries["function_type"]),
}

type AtOverload struct {
	type_         string
	documentation []string
}

func (ao *AtOverload) Type() string {
	return ao.type_
}

func (ao *AtOverload) Documentation() []string {
	return ao.documentation
}

func NewOverload(captures nvim.TsQueryMatch) *AtOverload {
	atOverload := AtOverload{
		documentation: []string{},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "type":
			atOverload.type_ = capture.Node.Text
		case "documentation":
			atOverload.documentation = append(atOverload.documentation, capture.Node.Text)
		}
	}

	return &atOverload
}
