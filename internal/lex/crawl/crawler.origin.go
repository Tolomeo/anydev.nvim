package crawl

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

var typeAnnotationQueries = map[string]string{
	"builtin_type":         "(builtin_type)",
	"identifier":           "(identifier)",
	"array_type":           "(array_type)",
	"table_type":           "(table_type)",
	"table_literal_type":   "(table_literal_type)",
	"union_type":           "(union_type)",
	"parenthesized_type":   "(parenthesized_type)",
	"tuple_type":           "(tuple_type)",
	"function_type":        "(function_type)",
	"member_type":          "(member_type)",
	"optional_type":        "(optional_type)",
	"literal_type":         "(literal_type)",
	"numeric_literal_type": "(numeric_literal_type)",
	"custom_type":          "(custom_type)",
}

var anyTypeAnnotationQuery = fmt.Sprintf(`[%s]`, strings.Join(mapx.Values(typeAnnotationQueries), " "))

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
					(alias_annotation
						"@alias"
						.
						(identifier) @alias.name
						.
						(%s) @alias.type
						.
						(comment)? @alias.documentation
						.
					) @origin.alias
					(#eq? @alias.name "%s")
					(#not-eq? @alias.type "")
				)`, anyTypeAnnotationQuery, c.context.Target().Identifier()),
			},
			// Luadoc matches an empty type node even when the type is not present
			// That means that enum aliases have an empty type node defined
			{
				Language: "luadoc",
				Query: fmt.Sprintf(`(
					(alias_annotation
						"@alias"
						.
						(identifier) @alias.name
						.
						(%s) @alias.emptytype
						.
						(comment)? @alias.documentation
						.
					) @origin.alias.enumerator
					(#eq? @alias.name "%s")
					(#eq? @alias.emptytype "")
				)`, anyTypeAnnotationQuery, c.context.Target().Identifier()),
			},
		},
		treesitter.CLASS_ANNOTATION: []treesitter.Query{
			{
				Language: "luadoc",
				Query: fmt.Sprintf(`(
					(class_annotation
						(identifier) @class.name
					) @class_annotation
					(#eq? @class.name "%s")
				) @origin.class`, c.context.Target().Identifier()),
			},
		},
		treesitter.FIELD_ANNOTATION: []treesitter.Query{
			{
				Language: "luadoc",
				Query: fmt.Sprintf(`(
					(field_annotation
						(identifier) @field.name
					) @field_annotation
					(#eq? @field.name "%s")
				) @origin.field`, c.context.Target().Name()),
			},
		},
	}
}

var enumAliasEnumeratorMemberQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(continuation
		(%s) @alias.type
	)
`, anyTypeAnnotationQuery)}

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
	queryResult, err := buffer.QueryTsNodeAt(originQueryMap, line, character)

	if err != nil {
		return nil, err
	}

	if queryResult == nil {
		return nil, nil
	}

	documentation, err := c.getCommentBlock(queryResult.Node, location)

	if err != nil {
		return nil, err
	}

	if _, isFunction := queryResult.Match.Find("origin.function"); isFunction {
		return symbol.NewFunctionOrigin(location, *queryResult, documentation), nil
	}

	if _, isTable := queryResult.Match.Find("origin.table"); isTable {
		return symbol.NewTableOrigin(location, *queryResult, documentation), nil
	}

	if _, isVariable := queryResult.Match.Find("origin.variable"); isVariable {
		return symbol.NewVariableOrigin(location, *queryResult, documentation), nil
	}

	if _, isModule := queryResult.Match.Find("origin.module"); isModule {
		return symbol.NewModuleOrigin(location, *queryResult, documentation), nil
	}

	if _, isClass := queryResult.Match.Find("origin.class"); isClass {
		return symbol.NewClassOrigin(location, *queryResult, documentation), nil
	}

	if _, isAlias := queryResult.Match.Find("origin.alias"); isAlias {
		return symbol.NewAliasOrigin(location, *queryResult, documentation), nil
	}

	if _, isAliasEnumerator := queryResult.Match.Find("origin.alias.enumerator"); isAliasEnumerator {
		c.context.Logger().Debugf("Alias enum match: <%+v>", queryResult.Match)

		nextLines, err := buffer.NextLineIterator(uint(queryResult.Match.LineRange().Start + 1))

		if err != nil {
			return nil, err
		}

		for line, err := range nextLines {
			if err != nil {
				return nil, err
			}

			match, err := line.TsQueryOne(enumAliasEnumeratorMemberQuery)

			if err != nil {
				return nil, err
			}

			if match == nil {
				break
			}

			queryResult.Match = queryResult.Match.Append(*match...)
		}

		return symbol.NewAliasEnumeratorOrigin(location, *queryResult, documentation), nil
	}

	if _, isField := queryResult.Match.Find("origin.field"); isField {
		return symbol.NewFieldOrigin(location, *queryResult, documentation), nil
	}

	if _, isMeta := queryResult.Match.Find("origin.meta"); isMeta {
		return symbol.NewMetaOrigin(location, *queryResult, documentation), nil
	}

	return nil, fmt.Errorf("Unknown origin match received: location <%+v>, definition <%+v>, documentation <%+v>", location, queryResult, documentation)
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
