package crawl

import (
	"fmt"
	"regexp"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
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
	var origin symbol.Origin
	var err error
	originQueryMap := c.getOriginQueryMap()

	for _, location := range locations {
		origin, err = c.getOrigin(location, originQueryMap)

		if err != nil {
			return nil, err
		}

		if origin != nil {
			break
		}
	}

	if origin == nil {
		return nil, nil
	}

	switch o := origin.(type) {
	case *symbol.ModuleOrigin:
		moduleOrigins, err := c.getModuleOrigins(o)
		if err != nil {
			return nil, err
		}
		return symbol.NewOrigins(origin).Merge(moduleOrigins), nil
	case *symbol.VariableOrigin:
		variableOrigins, err := c.getVariableOrigins(o)
		if err != nil {
			return nil, err
		}
		return symbol.NewOrigins(origin).Merge(variableOrigins), nil
	}

	return symbol.NewOrigins(origin), nil
}

func (c *Crawler) getOrigin(location nvim.Location, originQueryMap nvim.TsNodeQueryMap) (symbol.Origin, error) {
	buffer, err := c.context.Nvim().OpenBuffer(location.Url)

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	line, character :=
		uint(location.TargetRange.Start.Line),
		uint(location.TargetRange.Start.Character)
	definition, err := buffer.QueryTsNodeAt(originQueryMap, line, character)

	if err != nil {
		return nil, err
	}

	if definition == nil {
		return nil, nil
	}

	documentation, err := c.getCommentBlock(definition.Node, location)

	if err != nil {
		return nil, err
	}

	for _, capture := range definition.Match {
		switch capture.Id {
		case "origin.function":
			return symbol.NewFunctionOrigin(location, *definition, documentation), nil
		case "origin.table":
			return symbol.NewTableOrigin(location, *definition, documentation), nil
		case "origin.variable":
			return symbol.NewVariableOrigin(location, *definition, documentation), nil
		case "origin.module":
			return symbol.NewModuleOrigin(location, *definition, documentation), nil
		case "origin.class":
			return symbol.NewClassOrigin(location, *definition, documentation), nil
		case "origin.alias":
			return symbol.NewAliasOrigin(location, *definition, documentation), nil
		case "origin.field":
			return symbol.NewFieldOrigin(location, *definition, documentation), nil
		case "origin.meta":
			return symbol.NewMetaOrigin(location, *definition, documentation), nil
		}
	}

	return nil, fmt.Errorf("Unknown origin match received: location <%+v>, definition <%+v>, documentation <%+v>", location, definition, documentation)
}

func (c *Crawler) getModuleOrigins(origin *symbol.ModuleOrigin) (*symbol.Origins, error) {
	moduleName := origin.GetModuleName()
	moduleLocations, err := c.findModuleDefinitionLocations(moduleName)

	if err != nil {
		return nil, err
	}

	return c.getOrigins(*moduleLocations)
}

func (c *Crawler) getVariableOrigins(origin *symbol.VariableOrigin) (*symbol.Origins, error) {
	assignedName := origin.GetAssignedName()
	rightValueLocations, err := c.findDefinitionLocations(assignedName)

	if err != nil {
		return nil, err
	}

	return c.getOrigins(*rightValueLocations)
}
