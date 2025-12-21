package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

var (
	functionClassAccess   string = "class"
	functionIstanceAccess string = "instance"
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

func (l *Lexer) matchFunction(source *crawl.Source) (*lexed.Function, error) {
	function := newFunctionType()

	err := l.scratch(source.Origin().Definition())

	if err != nil {
		return nil, err
	}

	for _, query := range functionQueries {
		captures, err := l.context.nvim.TsQueryOne(nvim.TsQueryConfig{Language: "lua", Query: query})

		switch {
		case err != nil:
			return nil, err
		case captures == nil:
			continue
		}

		for _, capture := range *captures {
			switch capture.Id {
			case "name":
				function.Name = &capture.Node.Text
			case "access.class":
				function.Access = &functionClassAccess
			case "access.instance":
				function.Access = &functionIstanceAccess
			case "arg":
				function.Args = append(function.Args, newFunctionTypeArg(capture.Node.Text))
			case "vararg":
				function.Args = append(function.Args, newFunctionTypeArg(capture.Node.Text))
			}
		}

		function.Documentation = source.Origin().Documentation()
		return function, nil
	}

	return nil, nil
}

func (l *Lexer) lexFunction(function *lexed.Function) error {
	annotations, err := l.lexAnnotations(function.Documentation)

	if err != nil {
		return fmt.Errorf("Error lexing function %s: %w", *function.Name, err)
	}

	function.Overloads = annotations.overloads
	function.Generics = annotations.generics
	function.Return = annotations.returns

	for argIndex := range function.Args {
		name := function.Args[argIndex].Name
		annotation, hasAnnotation := annotations.params[name]

		if !hasAnnotation {
			l.context.logger.Info(fmt.Sprintf("Using '%s' for undocumented argument type '%s'", function.Args[argIndex].Type, function.Args[argIndex].Name))
			continue
		}

		function.Args[argIndex].Type = annotation.Type
		function.Args[argIndex].Optional = annotation.Optional
		function.Args[argIndex].Documentation = annotation.Documentation
	}

	return nil
}
