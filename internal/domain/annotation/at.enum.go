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
	name string
}

func (ae *AtEnum) Name() string {
	return ae.name
}

func NewAtEnum(captures nvim.TsQueryMatch) *AtEnum {
	atEnum := AtEnum{}

	for _, capture := range captures {
		switch capture.Id {
		case "enum.name":
			atEnum.name = capture.Node.Text
		}
	}

	return &atEnum
}
