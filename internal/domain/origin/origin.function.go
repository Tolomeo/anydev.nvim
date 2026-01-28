package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

/*
function fn() end
function fn(arg1) end
function fn(arg1, arg2) end
function fn(arg1, arg2, ...) end
function fn(...) end
*/
var FunctionDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(function_declaration
		name: (identifier) @name
		parameters: (parameters
			(identifier)? @arg
			("," (identifier) @arg)*
			("," (vararg_expression) @vararg)?
			(vararg_expression)? @vararg
		)
	) @function`,
}

/*
local T = function() end
local M = function(arg) end
local D = function(arg, ...) end
local E = function(...) end
*/
var FunctionVariableDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(variable_declaration
		(assignment_statement
			(variable_list
				name: (identifier) @name
			)
			(expression_list
				value: (function_definition
					parameters: (parameters
						(identifier)? @arg
						("," (identifier) @arg)*
						("," (vararg_expression) @vararg)?
						(vararg_expression)? @vararg
					)
				)
			)
		)
	) @function`,
}

/*
api['fn'] = function() end
api['fn'] = function(name) end
api['fn'] = function(name, value) end
api['fn'] = function(name, value, ...) end
api['fn'] = function(...) end
*/
var FunctionFieldIndexAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (bracket_index_expression
				table: (_)
				field: (string
					content: (string_content) @name
				)
			) @access.class
		)
		(expression_list
			value: (function_definition
				parameters: (parameters
					(identifier)? @arg
					("," (identifier) @arg)*
					("," (vararg_expression) @vararg)?
					(vararg_expression)? @vararg
				)
			)
		)
	) @function`,
}

/*
api.fn = function() end
api.fn = function(name) end
api.fn = function(name, value) end
api.fn = function(name, value, ...) end
api.fn = function(...) end
*/
var FunctionFieldDotAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (dot_index_expression
				field: (identifier) @name
			) @access.class
		)
		(expression_list
			value: (function_definition
				parameters: (parameters
					(identifier)? @arg
					("," (identifier) @arg)*
					("," (vararg_expression) @vararg)?
					(vararg_expression)? @vararg
				)
			)
		)
	) @function`,
}

/*
function api.fn() end
function api.fn(name) end
function api.fn(name, value) end
function api.fn(name, value, ...) end
function api.fn(...) end
*/
var FunctionFieldDotDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(function_declaration
		name: (dot_index_expression
			field: (identifier) @name
		) @access.class
		parameters: (parameters
			(identifier)? @arg
			("," (identifier) @arg)*
			("," (vararg_expression) @vararg)?
			(vararg_expression)? @vararg
		)
	) @function`,
}

/*
function api:fn() end
function api:fn(name) end
function api:fn(name, value) end
function api:fn(name, value, ...) end
function api:fn(...) end
*/
var FunctionFieldMethodDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(function_declaration
		name: (method_index_expression
			method: (identifier) @name
		) @access.instance
		parameters: (parameters
			(identifier)? @arg
			("," (identifier) @arg)*
			("," (vararg_expression) @vararg)?
			(vararg_expression)? @vararg
		)
	) @function`,
}

type FunctionOrigin struct {
	origin
}

func (fo *FunctionOrigin) Definition() []string {
	root, _ := fo.definition.Match.Find("function")
	return strings.Split(root.Node.Text, "\n")
}

func NewFunctionOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *FunctionOrigin {
	return &FunctionOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}
