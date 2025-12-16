package lex

import (
	"fmt"
	"strings"

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

func (l *Lexer) matchFunction(source *crawl.Source) (*lexed.Function, error) {
	function := newFunctionType()

	definitionLines := strings.Split(source.Origin().Definition(), "\n")
	err := l.scratch(definitionLines)

	if err != nil {
		return nil, err
	}

	for _, query := range functionQueries {
		captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "lua", Query: query})

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

	function.Private = annotations.private
	function.Protected = annotations.protected
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
