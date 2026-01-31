package annotation

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var ParenthesizedQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(parenthesized_type 
		(_) @group.type
	)`,
}

type Parenthesized struct {
	type_ string
}

func (p *Parenthesized) Type() string {
	return p.type_
}

func NewParenthesized(captures nvim.TsQueryMatch) *Parenthesized {
	parenthesized := Parenthesized{}

	for _, matchCapture := range captures {
		switch matchCapture.Id {
		case "group.type":
			parenthesized.type_ = matchCapture.Node.Text
		}
	}

	return &parenthesized
}
