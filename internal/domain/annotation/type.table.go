package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// TODO: check if it is possible to mark value as optional
var TableQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(table_type
		key: (%s) @key
		value: (%s) @value
	) @table`, AnyTypeQuery, AnyTypeQuery),
}

type Table struct {
	key   string
	value string
}

func (t *Table) Key() string {
	return t.key
}

func (t *Table) Value() string {
	return t.value
}

func NewTable(captures nvim.TsQueryMatch) *Table {
	table := Table{}

	for _, capture := range captures {
		switch capture.Id {
		case "key":
			table.key = capture.Node.Text
		case "value":
			table.value = capture.Node.Text
		}
	}

	return &table
}
