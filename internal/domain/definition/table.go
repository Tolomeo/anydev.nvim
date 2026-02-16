package definition

import (
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
		(#any-of? @function.name "vim._defer_require" "create_option_accessor")
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

type Table struct {
	root treesitter.TsNode
	name treesitter.TsNode
}

func (t *Table) Root() treesitter.TsNode {
	return t.root
}

func (t *Table) Name() treesitter.TsNode {
	return t.name
}

func NewTable(match nvim.TsQueryMatch) *Table {
	table := Table{}

	for _, capture := range match {
		switch capture.Id {
		case "table":
			table.root = capture.Node
		case "table.name":
			table.name = capture.Node
		}
	}

	return &table
}
