package definition

import (
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
	) @function
	(#not-has-ancestor? @function "function_declaration") ; Avoiding nested matches
	(#not-has-ancestor? @function "function_definition") ; Avoiding nested matches
	`,
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
				) @signature
				(#not-has-ancestor? @signature "function_declaration") ; Avoiding nested matches
				(#not-has-ancestor? @signature "function_definition") ; Avoiding nested matches
			)
		)
	) @function`,
}

/*
	F = {
		...
			fn = function() end
			fn = function(name) end
			fn = function(name, value) end
			fn = function(name, value, ...) end
			fn = function(...) end
		...
	}
*/
var FunctionFieldAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(field
		name: (identifier) @name
		value: (function_definition
			parameters: (parameters
				(identifier)? @arg
				("," (identifier) @arg)*
				("," (vararg_expression) @vararg)?
				(vararg_expression)? @vararg
			)
		) @signature
		; (#not-has-ancestor? @signature "function_declaration") ; Avoiding nested matches
		; (#not-has-ancestor? @signature "function_definition") ; Avoiding nested matches
	) @function 
	(#not-has-ancestor? @function "field") ; Avoiding nested matches
	`,
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
			) @signature
			(#not-has-ancestor? @signature "function_declaration") ; Avoiding nested matches
			(#not-has-ancestor? @signature "function_definition") ; Avoiding nested matches
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
			) @signature
			(#not-has-ancestor? @signature "function_declaration") ; Avoiding nested matches
			(#not-has-ancestor? @signature "function_definition") ; Avoiding nested matches
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
	) @function
	(#not-has-ancestor? @function "function_declaration") ; Avoiding nested matches
	(#not-has-ancestor? @function "function_definition") ; Avoiding nested matches
	`,
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
	) @function
	(#not-has-ancestor? @function "function_declaration") ; Avoiding nested matches
	(#not-has-ancestor? @function "function_definition") ; Avoiding nested matches
	`,
}

type Function struct {
	root   treesitter.TsNode
	name   treesitter.TsNode
	static *treesitter.TsNode
	args   []treesitter.TsNode
}

func (f *Function) Root() treesitter.TsNode {
	return f.root
}

func (f *Function) Name() treesitter.TsNode {
	return f.name
}

func (f *Function) Static() *treesitter.TsNode {
	return f.static
}

func (f *Function) Args() []treesitter.TsNode {
	return f.args
}

func NewFunction(match nvim.TsQueryMatch) *Function {
	function := Function{
		args: []treesitter.TsNode{},
	}

	for _, capture := range match {
		switch capture.Id {
		case "function":
			function.root = capture.Node
		case "name":
			function.name = capture.Node
		case "access.class":
			function.static = &capture.Node
		case "arg":
			switch capture.Node.Text {
			case "self":
				function.static = &capture.Node
			default:
				function.args = append(function.args, capture.Node)
			}
		case "vararg":
			function.args = append(function.args, capture.Node)
		}
	}

	return &function
}
