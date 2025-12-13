package lex

import (
	"errors"
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
	"dotIndexDeclaration ": `
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
}

func (l *Lexer) lexFunctionDefinition(functionDefinition string, function *lexed.Function) error {
	l.scratch([]string{functionDefinition})

	for _, query := range functionQueries {
		captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "lua", Query: query})

		switch {
		case errors.Is(nvim.ErrTSQueryNoMatch, err):
			continue
		case err != nil:
			return err
		}

		for _, capture := range captures {
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

		return nil
	}

	return fmt.Errorf("Function source %v didn't yield any result", functionDefinition)
}

func (l *Lexer) lexFunction(source *crawl.Source) (lexed.Symbol, error) {
	function := newFunctionType()

	function.Documentation = source.Origin().Documentation()

	annotations, err := l.lexAnnotations(function.Documentation)

	if err != nil {
		return struct{}{}, fmt.Errorf("Error lexing function %s: %w", source.Path(), err)
	}

	function.Private = annotations.private
	function.Protected = annotations.protected
	function.Overloads = annotations.overloads
	function.Generics = annotations.generics
	function.Return = annotations.returns

	functionDefinition := source.Origin().Definition()
	err = l.lexFunctionDefinition(functionDefinition, &function)

	if err != nil {
		return struct{}{}, fmt.Errorf("Error lexing function %s: %w", source.Path(), err)
	}

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

	return function, nil
}
