package annotation

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var OptionalQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(optional_type
		(_) @optional.type
	) @optional`,
}

type Optional struct {
	type_ string
}

func (o *Optional) Type() string {
	return o.type_
}

func NewOptional(captures nvim.TsQueryMatch) *Optional {
	optional := Optional{}

	for _, capture := range captures {
		switch capture.Id {
		case "optional.type":
			optional.type_ = capture.Node.Text
		}
	}

	return &optional
}
