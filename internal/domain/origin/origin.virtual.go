package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// local F = ...
var MetaVariableAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: [
				(vararg_expression) @assignment.right
			] 
		)
	) @meta`,
}

type MetaOrigin struct {
	origin
}

func (mo *MetaOrigin) Definition() []string {
	root, _ := mo.definition.Match.Find("meta")
	return strings.Split(root.Node.Text, "\n")
}

func NewMetaOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *MetaOrigin {
	return &MetaOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}
