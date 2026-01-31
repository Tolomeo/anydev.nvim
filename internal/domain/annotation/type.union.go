package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var UnionQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(union_type
		(%s) @union.type
		(%s) @union.type
	) @union`, AnyTypeQuery, AnyTypeQuery),
}

type Union struct {
	types []string
}

func (u *Union) Types() []string {
	return u.types
}

func NewUnion(captures nvim.TsQueryMatch) *Union {
	union := Union{
		types: []string{},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "union.type":
			union.types = append(union.types, capture.Node.Text)
		}
	}

	return &union
}
