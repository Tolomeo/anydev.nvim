package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// local T = {}
var TableVariableDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(variable_declaration
		(assignment_statement
			(variable_list
				name: (identifier)
			) @table.name
			(expression_list
				value: (table_constructor)
			) @table.value
		) 
	) @table`,
}

// F.T = {}
var TableFieldAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (dot_index_expression
				table: (_)
				field: (identifier) @table.name
			)
		)
		(expression_list
			value: (table_constructor) @table.value
		)
	) @table`,
}

// F['T'] = {}
var TableFieldIndexAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (bracket_index_expression
				table: (_)
				field: (string
					content: (string_content) @table.name
				)
			)
		)
		(expression_list
			value: (table_constructor) @table.value
		)
	) @table`,
}

type TableOrigin struct {
	origin
}

func (to *TableOrigin) Definition() []string {
	root, _ := to.definition.Match.Find("table")
	return strings.Split(root.Node.Text, "\n")
}

func NewTableOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *TableOrigin {
	return &TableOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}
