package lex

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	// "github.com/Tolomeo/anydev.nvim/internal/nvim/ts"
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

func newFunctionArg(name string) lexed.FunctionArg {
	return lexed.FunctionArg{
		Name:     name,
		Type:     "any",
		Optional: false,
	}
}

func (l *lexer) lexFunctionDefinition(functionDefinition []string, function *lexed.Function) error {
	l.scratch(functionDefinition)

	for _, query := range functionQueries {
		captures, err := l.nvim.TsQuery("lua", query)

		switch {
		case errors.Is(nvim.ErrNotFound, err):
			continue
		case err != nil:
			return err
		}

		for _, capture := range captures {
			switch capture.Id {
			case "name":
				function.Name = strings.Join(capture.Node.Text, "")
			case "access.class":
				function.Access = &functionClassAccess
			case "access.instance":
				function.Access = &functionIstanceAccess
			case "arg":
				argName := strings.Join(capture.Node.Text, "")
				function.Args = append(function.Args, newFunctionArg(argName))
			case "vararg":
				argName := strings.Join(capture.Node.Text, "")
				function.Args = append(function.Args, newFunctionArg(argName))
			}
		}

		return nil
	}

	return fmt.Errorf("Function source %v didn't yield any result", functionDefinition)
}

func (l *lexer) lexFunction(source *crawl.FunctionSource) error {
	function := lexed.Function{}

	annotations, err := l.lexAnnotations(source.Origin().Documentation)

	if err != nil {
		return fmt.Errorf("Error lexing function %s: %w", source.Path(), err)
	}

	functionDefinition := source.Origin().Definition
	err = l.lexFunctionDefinition(functionDefinition, &function)

	if err != nil {
		return fmt.Errorf("Error lexing function %s: %w", source.Path(), err)
	}

	for argIndex := range function.Args {
		name := function.Args[argIndex].Name
		annotation, hasAnnotation := annotations.params[name]

		if !hasAnnotation {
			//TODO: trace not found param annotation
			continue
		}

		function.Args[argIndex].Type = annotation.Type
		function.Args[argIndex].Optional = annotation.Optional
		function.Args[argIndex].Documentation = annotation.Documentation
	}

	fmt.Printf("%+v", function)

	return nil
}
