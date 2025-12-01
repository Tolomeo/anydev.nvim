package extract

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/extract/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/extract/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

var (
	functionClassAccess string = "class"
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

func newFunctionArg(name string) symbol.LexedFunctionArg {
	return symbol.LexedFunctionArg{
		Name:     name,
		Type:     "uknown",
		Optional: false,
	}
}

func newFunctionVararg() symbol.LexedFunctionVararg {
	return symbol.LexedFunctionVararg{
		Type: "uknown",
	}
}

func (e *extractor) lexFunction(source *crawl.FunctionSource) error {
	lexedFunction := symbol.LexedFunction{}

	err := e.lexFunctionDefinition(source.Origin().Definition, &lexedFunction)

	if err != nil {
		return fmt.Errorf("Error lexing function %s source: %w", source.Path(), err)
	}

	err = e.lexFunctionDocumentation(source.Origin().Documentation, &lexedFunction)

	if err != nil {
		return fmt.Errorf("Error lexing function %s source: %w", source.Path(), err)
	}

	fmt.Printf("%+v", lexedFunction)

	return nil
}

func (e *extractor) lexFunctionDocumentation(documentation []string, function *symbol.LexedFunction) error {
	e.scratch(documentation)

	for _, argument := range function.Args {
		switch arg := argument.(type) {
		case symbol.LexedFunctionArg:
			fmt.Println("arg", arg.Name)
		case symbol.LexedFunctionVararg:
			fmt.Println("vararg")
		}
	}

	return nil
}

func (e *extractor) lexFunctionDefinition(definition []string, function *symbol.LexedFunction) error {
	e.scratch(definition)

	for _, query := range functionQueries {
		captures, err := e.nvim.TsQuery(query)

		switch {
		case errors.Is(nvim.ErrNotFound, err):
			continue
		case err != nil:
			fmt.Printf("%v", err)
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
				function.Args = append(function.Args, newFunctionArg(strings.Join(capture.Node.Text, "")))
			case "vararg":
				function.Args = append(function.Args, newFunctionVararg())
			}
		}

		return nil
	}

	return fmt.Errorf("Function source %v didn't yield any result", definition)
}
