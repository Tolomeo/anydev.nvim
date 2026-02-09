package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

/*
	F = {
		...
	    T = {},
		...
	}
*/
var TableConstructorFieldAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
		(field
			name: (identifier) @table.name
			value: (table_constructor) @table.value
		) @table
		(#not-has-ancestor? @table "field") ; Avoiding nested matches
	`,
}

/*
	F = {
		...
	    ['T'] = {},
		...
	}
*/
var TableConstructorFieldIndexAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(
		(field
			name: (string
				content: (string_content) @table.name
			)
			value: (table_constructor) @table.value
		) @field 
		(#not-has-ancestor? @field "field") ; Avoiding nested matches
	) @table`,
}

// T = {}
var TableDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (identifier)
		) @table.name
		(expression_list
			value: (table_constructor)
		) @table.value
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

var TableReturnAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement 
		(variable_list 
			name: (_) @table.name
		) 
		(expression_list 
			value: (function_call 
				name: (_) @function.name 
			)
		)
		(#any-eq? @function.name "vim._defer_require")
	) @table`,
}

var MetatableAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list name: (_) @table.name)
		(expression_list
			value: 
				(function_call 
					name: (identifier) @setmetatable
				)
		) @table.value
		(#eq? @setmetatable "setmetatable")
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

func NewTableOrigin(location nvim.Location, captures nvim.TsQueryMatch, documentation []string) *TableOrigin {
	tableOrigin := TableOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: documentation,
		},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "table":
			tableOrigin.definition = capture.Node.Text
		case "table.name":
			tableOrigin.name = capture.Node.Text
		}
	}

	return &tableOrigin
}
