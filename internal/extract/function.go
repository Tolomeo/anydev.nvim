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

func (e *extractor) lexFunctionDocumentation(functionDocumentation []string, function *symbol.LexedFunction) error {
	paramAnnotationQuery := `
		(documentation) @documentation
	`

	for _, documentationLine := range functionDocumentation {
		err := e.scratch([]string{documentationLine})

		if err != nil {
			return err
		}

		fmt.Println(documentationLine)

		captures, err := e.nvim.TsQuery("luadoc", paramAnnotationQuery)

		switch {
		case errors.Is(nvim.ErrNotFound, err):
			continue
		case err != nil:
			return err
		}

		fmt.Printf("%+v\n", captures)
	}

	return nil
}

func (e *extractor) lexFunctionDefinition(functionDefinition []string, function *symbol.LexedFunction) error {
	e.scratch(functionDefinition)

	for _, query := range functionQueries {
		captures, err := e.nvim.TsQuery("lua", query)

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
				function.Args = append(function.Args, newFunctionVararg())
			}
		}

		return nil
	}

	return fmt.Errorf("Function source %v didn't yield any result", functionDefinition)
}

func (e *extractor) lexFunction(source *crawl.FunctionSource) error {
	// fmt.Printf("%+v\n", source.Origin())

	lexedFunction := symbol.LexedFunction{}

	functionDefinition := source.Origin().Definition
	err := e.lexFunctionDefinition(functionDefinition, &lexedFunction)

	if err != nil {
		return fmt.Errorf("Error lexing function %s definition: %w", source.Path(), err)
	}

	functionDocumentation := source.Origin().Documentation
	err = e.lexFunctionDocumentation(functionDocumentation, &lexedFunction)

	if err != nil {
		fmt.Printf("%v", err)
		return fmt.Errorf("Error lexing function %s documentation: %w", source.Path(), err)
	}

	// fmt.Printf("%+v", lexedFunction)

	return nil
}
