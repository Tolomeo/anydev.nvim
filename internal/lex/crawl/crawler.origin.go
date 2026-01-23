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

func (c *Crawler) findOrigin(locations []nvim.Location) (*symbol.Origins, error) {
	c.context.Logger().Debugf("QueryMap for %s: \n %+v", c.context.Target().Identifier(), c.getOriginQueryMap())

	for _, location := range locations {
		c.context.Logger().Debugf("Location: %+v", location.Url)

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

		c.context.Logger().Debugf("QueryMatch for %s: \n %+v", c.context.Target().Identifier(), queryMatch)

		if queryMatch == nil {
			continue
		}

		var origin *symbol.Origin

		documentation, err := c.getCommentBlock(queryMatch.Node, location)

		if err != nil {
			return nil, err
		}

		if documentation != nil {
			origin = symbol.NewOrigin(location, queryMatch.Node, *documentation)
		} else {
			c.context.Logger().Warnf("No documentation found for location <%v>", location)
			origin = symbol.NewOrigin(location, queryMatch.Node, treesitter.TsNode{})
		}

		targetOrigin := symbol.NewOrigins(origin)

		switch targetOrigin.Last().Type() {
		case treesitter.ASSIGNMENT_STATEMENT:
			followed, err := c.followModuleRequireAssignment(targetOrigin)

			if err != nil {
				return nil, err
			}

			if followed {
				return targetOrigin, nil
			}

			followed, err = c.followVariableAssignment(targetOrigin)

			if err != nil {
				return nil, err
			}

			if followed {
				return targetOrigin, nil
			}
		}

		return targetOrigin, nil
	}

	return nil, nil
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

func (c *Crawler) followModuleRequireAssignment(origin *symbol.Origins) (bool, error) {
	buffer, err := c.context.Nvim().NewBuffer()

	if err != nil {
		return false, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.Last().Definition())

	if err != nil {
		return false, err
	}

	captures, err := buffer.TsQueryOne(requireAssignmentQuery)

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	moduleNameCapture, found := slicesx.FindFunc(*captures, func(capture treesitter.Capture) bool {
		return capture.Id == "require.module"
	})

	if !found {
		return false, fmt.Errorf("Error retrieving required module name from require statement in '%s'", origin.Last().Definition())
	}

	moduleLocations, err := c.findModuleDefinitionLocations(moduleNameCapture.Node.Text)

	if err != nil {
		return false, err
	}

	moduleOrigin, err := c.findOrigin(*moduleLocations)

	switch {
	case err != nil:
		return false, err
	case moduleOrigin == nil:
		return false, fmt.Errorf("Error following require statement '%s'", origin.Last().Definition())
	}

	origin.Merge(moduleOrigin)

	return true, nil
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

func (c *Crawler) followVariableAssignment(origin *symbol.Origins) (bool, error) {
	buffer, err := c.context.Nvim().NewBuffer()

	if err != nil {
		return false, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.Last().Definition())

	if err != nil {
		return false, err
	}

	var captures *nvim.TsQueryMatch = nil

	for _, variableAssignmentQuery := range variableAssignmentQueries {
		assignmentCaptures, err := buffer.TsQueryOne(variableAssignmentQuery)

		switch {
		case err != nil:
			return false, err
		case assignmentCaptures == nil:
			continue
		}

		captures = assignmentCaptures
	}

	if captures == nil {
		return false, nil
	}

	rightValue, found := slicesx.FindFunc(*captures, func(capture treesitter.Capture) bool {
		return capture.Id == "assignment.right"
	})

	if !found {
		return false, fmt.Errorf("Error retrieving read variable name from variable to variable assignment in '%s'", origin.Last().Definition())
	}

	rightValueLocations, err := c.findDefinitionLocations(rightValue.Node.Text)

	if err != nil {
		return false, err
	}

	rightValueOrigin, err := c.findOrigin(*rightValueLocations)

	if err != nil {
		return false, err
	}

	origin.Merge(rightValueOrigin)

	return true, nil
}
