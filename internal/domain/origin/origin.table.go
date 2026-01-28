package origin

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
	definition string
	name       string
}

func (to *TableOrigin) Definition() []string {
	return strings.Split(to.definition, "\n")
}

func (to *TableOrigin) Name() string {
	return to.name
}

func NewTableOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *TableOrigin {
	tableOrigin := TableOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}

	for _, capture := range definition.Match {
		switch capture.Id {
		case "table":
			tableOrigin.definition = capture.Node.Text
		case "table.name":
			tableOrigin.name = capture.Node.Text
		}
	}

	return &tableOrigin
}
