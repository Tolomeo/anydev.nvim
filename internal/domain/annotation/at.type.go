package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AtTypeQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
		(type_annotation
			"@type" . (%s) @type.type
			.
			("," . (%s) @type.type)*
			.
			(comment)? @type.documentation
			.
		) @type`, AnyTypeQuery, AnyTypeQuery),
}

type AtType struct {
	types []string
}

func (at *AtType) Types() []string {
	return at.types
}

func NewAtType(captures nvim.TsQueryMatch) *AtType {
	atType := AtType{}

	for _, capture := range captures {
		switch capture.Id {
		case "type.type":
			atType.types = append(atType.types, capture.Node.Text)
		}
	}

	return &atType
}
