package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var functionQueries = map[string]string{
	/* function fn() end
	function fn(arg1) end
	function fn(arg1, arg2) end
	function fn(arg1, arg2, ...) end
		function fn(...) end */
	"declaration": `
		(function_declaration
			name: (identifier) @name
			parameters: (parameters
				(identifier)? @arg
				("," (identifier) @arg)*
				("," (vararg_expression) @vararg)?
				(vararg_expression)? @vararg
			)
		)
	`,
	/* local T = function() end
	local M = function(arg) end
	local D = function(arg, ...) end
	local E = function(...) end */
	"assignment": `
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
		)
	`,
	/* function api:fn() end
	function api:fn(name) end
	function api:fn(name, value) end
	function api:fn(name, value, ...) end
	function api:fn(...) end */
	"methodDeclaration": `
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
		)
	`,
	/* function api.fn() end
	function api.fn(name) end
	function api.fn(name, value) end
	function api.fn(name, value, ...) end
	function api.fn(...) end */
	"methodIndexDeclaration ": `
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
		)
	`,
	/* api.fn = function() end
	api.fn = function(name) end
	api.fn = function(name, value) end
	api.fn = function(name, value, ...) end
	api.fn = function(...) end */
	"methodAssignment": `
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
		)
	`,
	/* api['fn'] = function() end
	api['fn'] = function(name) end
	api['fn'] = function(name, value) end
	api['fn'] = function(name, value, ...) end
	api['fn'] = function(...) end */
	"methodIndexAssignment": `
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
		)
	`,
}

func (l *Lexer) lexFunctionValue(source *symbol.ValueSource) (*symbol.Function, error) {
	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(source.Origin.DefinitionLines())

	if err != nil {
		return nil, err
	}

	var match *nvim.TsQueryMatch

	for _, query := range functionQueries {
		captures, err := buffer.TsQueryOne(treesitter.Query{Language: "lua", Query: query})

		if err != nil {
			return nil, err
		}

		if captures != nil {
			match = captures
			break
		}
	}

	if match == nil {
		return nil, nil
	}

	function := symbol.NewFunction()
	function.Documentation = source.Origin.DocumentationLines()

	for _, capture := range *match {
		switch capture.Id {
		case "name":
			function.Name = &capture.Node.Text
		case "access.class":
			function.Access = &symbol.FunctionClassAccess
		case "access.instance":
			function.Access = &symbol.FunctionIstanceAccess
		case "arg":
			function.Arguments = append(function.Arguments, *symbol.NewFunctionArgument(capture.Node.Text))
		case "vararg":
			function.Arguments = append(function.Arguments, *symbol.NewFunctionArgument(capture.Node.Text))
		}
	}

	annotations, err := l.lexAnnotations(function.Documentation)

	if err != nil {
		return nil, fmt.Errorf("Error lexing function %s: %w", *function.Name, err)
	}

	function.Overloads = annotations.overloads
	function.Generics = annotations.generics
	function.Returns = annotations.returns

	for argIndex := range function.Arguments {
		name := function.Arguments[argIndex].Name
		annotation, hasAnnotation := annotations.params[name]

		if !hasAnnotation {
			l.context.Logger.Info(fmt.Sprintf("Using '%s' for undocumented argument type '%s'", function.Arguments[argIndex].Type, function.Arguments[argIndex].Name))
			continue
		}

		function.Arguments[argIndex].Type = annotation.Type
		function.Arguments[argIndex].Optional = annotation.Optional
		function.Arguments[argIndex].Documentation = annotation.Documentation
	}

	return function, nil
}

var tableQueries = map[string]string{
	"tableDeclaration": `
	(variable_declaration
		(assignment_statement
			(variable_list
				name: (identifier)
			) @table.name
			(expression_list
				value: (table_constructor)
			) @table.value
		)
	)
`,
	"tableFieldDeclaration": `
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
	)
`,
	"tableIndexFieldDeclaration": `
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
		)
`}

func (l *Lexer) lexTableValue(source *symbol.ValueSource) (*symbol.Table, error) {
	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(source.Origin.DefinitionLines())

	if err != nil {
		return nil, err
	}

	var match *nvim.TsQueryMatch

	for _, query := range tableQueries {
		captures, err := buffer.TsQueryOne(treesitter.Query{Language: "lua", Query: query})

		if err != nil {
			return nil, err
		}

		if captures != nil {
			match = captures
			break
		}
	}

	if match == nil {
		return nil, nil
	}

	table := symbol.NewTable()

	for _, capture := range *match {
		switch capture.Id {
		case "table.name":
			table.Name = &capture.Node.Text
		}
	}

	tableFields, err := l.context.Nvim.GetCompletion(source.Path)

	if err != nil {
		return nil, err
	}

	for _, fieldName := range tableFields {
		tableField := symbol.TableField{Name: fieldName}

		err := l.context.Push(fieldName, func(path string) error {
			source, err := l.sourceValue(l.context.Current())

			if err != nil {
				return err
			}

			annotations, err := l.lexAnnotations(source.Origin.DocumentationLines())

			tableField.Private = annotations.private
			tableField.Protected = annotations.protected
			tableFieldValue, err := l.lexValue(source)

			if err != nil {
				return err
			}

			tableField.Value = tableFieldValue

			return nil
		})

		if err != nil {
			return nil, err
		}

		table.Fields = append(table.Fields, tableField)
	}

	return table, nil
}

var metaQuery string = `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: [
				(vararg_expression) @assignment.right
			] 
		)
	) @assignment
`

func (l *Lexer) lexMetaValue(source *symbol.ValueSource) (symbol.Symbol, error) {
	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(source.Origin.DefinitionLines())

	if err != nil {
		return nil, err
	}

	captures, err := buffer.TsQueryOne(treesitter.Query{Language: "lua", Query: metaQuery})

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	unknown := symbol.NewUnknown()
	unknown.Documentation = source.Origin.DocumentationLines()

	annotations, err := l.lexAnnotations(unknown.Documentation)

	if err != nil {
		return nil, err
	}

	if annotations.tipe != nil {
		return &annotations.tipe, nil
	}

	l.context.Logger.Warn(fmt.Sprintf("Unknown meta type '%s' received", l.context.Current()))
	return unknown, nil
}
