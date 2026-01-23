package crawl

import (
	"fmt"
	"regexp"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (c *Crawler) getOriginQueryMap() nvim.TsNodeQueryMap {
	return nvim.TsNodeQueryMap{
		treesitter.ASSIGNMENT_STATEMENT: []treesitter.Query{
			// table field assignment
			// F.T = {}
			{
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
				) @origin.table`,
			},
			// table field index assignment
			// F['T'] = {}
			{
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
				) @origin.table`,
			},
			// meta
			// local F = ...
			{
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
				) @origin.meta`,
			},
			// method assignment
			/* api.fn = function() end
			api.fn = function(name) end
			api.fn = function(name, value) end
			api.fn = function(name, value, ...) end
			api.fn = function(...) end */
			{
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
				) @origin.function`,
			},
			/* api['fn'] = function() end
			api['fn'] = function(name) end
			api['fn'] = function(name, value) end
			api['fn'] = function(name, value, ...) end
			api['fn'] = function(...) end */
			{
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
				) @origin.function`,
			},
			// module assignment
			// F = require("T")
			{
				Language: "lua",
				Query: `
				(assignment_statement
					(variable_list)
					(expression_list
						value: (function_call
							name: (identifier) @require.call
							arguments: (arguments
								(string
									content: (string_content) @require.module
								)
							)
						)
					) @require
					(#eq? @require.call "require")
				) @origin.module`,
			},
			// dotindex assignment
			// F.T = X.V
			{
				Language: "lua",
				Query: `
				(assignment_statement
					(variable_list
						name: (_)
					) @assignment.left
					(expression_list
						value: (dot_index_expression) @assignment.right 
					)
				) @origin.variable`,
			},
			// variable assignment
			// F.T = X
			{
				Language: "lua",
				Query: `
				(assignment_statement
					(variable_list
						name: (_)
					) @assignment.left
					(expression_list
						value: (identifier) @assignment.right 
					)
				) @origin.variable`,
			},
		},
		treesitter.VARIABLE_DECLARATION: []treesitter.Query{
			// static method assignment
			/* local T = function() end
			local M = function(arg) end
			local D = function(arg, ...) end
			local E = function(...) end */
			{
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
				) @origin.function`,
			},
			// table assignment
			// local T = {}
			{
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
				) @origin.table`,
			},
		},
		treesitter.FUNCTION_DECLARATION: []treesitter.Query{
			// function
			/* function fn() end
			function fn(arg1) end
			function fn(arg1, arg2) end
			function fn(arg1, arg2, ...) end
				function fn(...) end */
			{
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
				) @origin.function`,
			},
			// instance method
			/* function api:fn() end
			function api:fn(name) end
			function api:fn(name, value) end
			function api:fn(name, value, ...) end
			function api:fn(...) end */
			{
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
				) @origin.function`,
			},
			// static method
			/* function api.fn() end
			function api.fn(name) end
			function api.fn(name, value) end
			function api.fn(name, value, ...) end
			function api.fn(...) end */
			{
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
				) @origin.function`,
			},
		},
		treesitter.ALIAS_ANNOTATION: []treesitter.Query{
			{
				Language: "luadoc",
				Query: fmt.Sprintf(`(
					(alias_annotation) @alias
					(#match? @alias "\\@alias *%s($|[^a-zA-Z0-9_])")
				) @origin.alias`, regexp.QuoteMeta(c.context.Target().Identifier())),
			},
		},
		treesitter.CLASS_ANNOTATION: []treesitter.Query{
			{
				Language: "luadoc",
				Query: fmt.Sprintf(`(
					(class_annotation) @class_annotation
					(#match? @class_annotation "\\@class *%s($|[^a-zA-Z0-9_])")
				) @origin.class`, regexp.QuoteMeta(c.context.Target().Identifier())),
			},
		},
		treesitter.FIELD_ANNOTATION: []treesitter.Query{
			{
				Language: "luadoc",
				Query: fmt.Sprintf(`(
					(field_annotation) @field_annotation
					(#match? @field_annotation "\\@field *%s($|[^a-zA-Z0-9_])")
				) @origin.field`, regexp.QuoteMeta(c.context.Target().Name())),
			},
		},
	}
}

func (c *Crawler) getOrigins(locations []nvim.Location) (*symbol.Origins, error) {
	var origin *symbol.Origin

	for _, location := range locations {
		buffer, err := c.context.Nvim().OpenBuffer(location.Url)

		if err != nil {
			return nil, err
		}

		defer buffer.Close()

		queryMap := c.getOriginQueryMap()
		line, character :=
			uint(location.TargetRange.Start.Line),
			uint(location.TargetRange.Start.Character)
		queryMatch, err := buffer.QueryTsNodeAt(queryMap, line, character)

		if err != nil {
			return nil, err
		}

		if queryMatch == nil {
			continue
		}

		node := queryMatch.Node
		documentation, err := c.getCommentBlock(node, location)

		if err != nil {
			return nil, err
		}

		origin = symbol.NewOrigin(location, node, documentation)
		break
	}

	if origin == nil {
		return nil, nil
	}

	origins := symbol.NewOrigins(origin)

	switch origins.Last().Type() {
	case treesitter.ASSIGNMENT_STATEMENT:
		moduleOrigins, err := c.getModuleOrigins(origins.Last())

		if err != nil {
			return nil, err
		}

		if moduleOrigins != nil {
			origins.Merge(moduleOrigins)
			return origins, nil
		}

		variableOrigins, err := c.getVariableOrigins(origins.Last())

		if err != nil {
			return nil, err
		}

		if variableOrigins != nil {
			origins.Merge(variableOrigins)
			return origins, nil
		}
	}

	return origins, nil
}

var requireAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list)
		(expression_list
			value: (function_call
				name: (identifier) @require.call
				arguments: (arguments
					(string
						content: (string_content) @require.module
					)
				)
			)
		) @require
		(#eq? @require.call "require")
	)`,
}

func (c *Crawler) getModuleOrigins(origin *symbol.Origin) (*symbol.Origins, error) {
	buffer, err := c.context.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.Definition())

	if err != nil {
		return nil, err
	}

	captures, err := buffer.TsQueryOne(requireAssignmentQuery)

	if err != nil {
		return nil, err
	}

	if captures == nil {
		return nil, nil
	}

	moduleNameCapture, found := slicesx.FindFunc(*captures, func(capture treesitter.Capture) bool {
		return capture.Id == "require.module"
	})

	if !found {
		return nil, fmt.Errorf("Error retrieving required module name from require statement in '%s'", origin.Definition())
	}

	moduleLocations, err := c.findModuleDefinitionLocations(moduleNameCapture.Node.Text)

	if err != nil {
		return nil, err
	}

	return c.getOrigins(*moduleLocations)
	/* moduleOrigins, err := c.getOrigins(*moduleLocations)

	switch {
	case err != nil:
		return false, err
	case moduleOrigins == nil:
		return false, fmt.Errorf("Error following require statement '%s'", origins.Last().Definition())
	}

	origins.Merge(moduleOrigins)

	return true, nil */
}

var variableAssignmentQueries = map[string]treesitter.Query{
	"dotIndexAssignment": {
		Language: "lua",
		Query: `
		(assignment_statement
			(variable_list
				name: (_)
			) @assignment.left
			(expression_list
				value: (dot_index_expression) @assignment.right 
			)
		)`},
	"identifierAssignment": {
		Language: "lua",
		Query: `
		(assignment_statement
			(variable_list
				name: (_)
			) @assignment.left
			(expression_list
				value: (identifier) @assignment.right 
			)
		)`,
	},
}

func (c *Crawler) getVariableOrigins(origin *symbol.Origin) (*symbol.Origins, error) {
	buffer, err := c.context.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.Definition())

	if err != nil {
		return nil, err
	}

	var captures *nvim.TsQueryMatch = nil

	for _, variableAssignmentQuery := range variableAssignmentQueries {
		assignmentCaptures, err := buffer.TsQueryOne(variableAssignmentQuery)

		if err != nil {
			return nil, err
		}

		if assignmentCaptures == nil {
			continue
		}

		captures = assignmentCaptures
	}

	if captures == nil {
		return nil, nil
	}

	rightValue, found := slicesx.FindFunc(*captures, func(capture treesitter.Capture) bool {
		return capture.Id == "assignment.right"
	})

	if !found {
		return nil, fmt.Errorf("Error retrieving read variable name from variable to variable assignment in '%s'", origin.Definition())
	}

	rightValueLocations, err := c.findDefinitionLocations(rightValue.Node.Text)

	if err != nil {
		return nil, err
	}

	rightValueOrigins, err := c.getOrigins(*rightValueLocations)

	if err != nil {
		return nil, err
	}

	return rightValueOrigins, nil
}
