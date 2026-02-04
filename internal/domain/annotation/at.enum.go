package annotation

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AtEnumQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(enum_annotation
		"@enum"
		.
		(identifier) @enum.name
	) @enum`,
}

type AtEnum struct {
	text string
	name string
}

func (ae *AtEnum) Text() string {
	return ae.text
}

func (ae *AtEnum) Name() string {
	return ae.name
}

func NewAtEnum(captures nvim.TsQueryMatch) *AtEnum {
	atEnum := AtEnum{}

	for _, capture := range captures {
		switch capture.Id {
		case "enum":
			atEnum.text = capture.Node.Text
		case "enum.name":
			atEnum.name = capture.Node.Text
		}
	}

	return &atEnum
}
