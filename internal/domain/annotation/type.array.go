package annotation

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var ArrayQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(array_type 
		(_) @array.itemstype 
	) @array`,
}

type Array struct {
	itemstype string
}

func (a *Array) ItemsType() string {
	return a.itemstype
}

func NewArray(captures nvim.TsQueryMatch) *Array {
	array := Array{}

	for _, capture := range captures {
		switch capture.Id {
		case "array.itemstype":
			array.itemstype = capture.Node.Text
		}
	}

	return &array
}
