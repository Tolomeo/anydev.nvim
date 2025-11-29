package extract

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/extract/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/extract/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

type extractor struct {
	nvim   *nvim.Nvim
	buffer string
}

func (e *extractor) Extract(path string) error {
	err := e.nvim.Start()

	if err != nil {
		return fmt.Errorf("Error opening nvim: %w", err)
	}

	crawler := crawl.NewCrawler(e.nvim)

	source, err := crawler.CrawlRuntime(path)

	if err != nil {
		return fmt.Errorf("Error crawling %s: %w", path, err)
	}

	e.serialize(source)

	outputDir, err := project.GetOutputDir()

	if err != nil {
		return fmt.Errorf("Error getting output location: %w", err)
	}

	out := output.NewOutput(outputDir)

	if err := out.WriteFile("statistics.json", crawler.Statistics()); err != nil {
		return fmt.Errorf("Error collecting crawler statistics: %w", err)
	}

	err = e.nvim.Quit()

	if err != nil {
		return fmt.Errorf("Errot closing nvim process gracefully: %w", err)
	}

	return nil
}

func (e *extractor) scratch(lines []string) error {
	buffer := path.Join(e.nvim.Options().Config().Dir(), "anydev.extractor.lua")

	_, err := e.nvim.Open(buffer)

	if err != nil {
		return err
	}

	err = e.nvim.SetBufferLines(lines)

	if err != nil {
		return err
	}

	return nil
}

func (e *extractor) serialize(source crawl.Source) error {
	switch v := source.(type) {
	case *crawl.TableSource:
		fmt.Println(v, "table")
	case *crawl.FunctionSource:
		return e.lexFunction(v)
	case *crawl.VariableSource:
		fmt.Println("variable")
	}
	// fmt.Printf("%+v", source.Origin())
	return nil
}

func (e *extractor) lexFunction(source *crawl.FunctionSource) error {
	lexedFunction := symbol.LexedFunctionSource{}

	err := e.lexFunctionDefinition(source.Origin().Definition, &lexedFunction)

	if err != nil {
		return fmt.Errorf("Error lexing function %s source: %w", source.Path(), err)
	}

	fmt.Printf("%+v", lexedFunction)

	return nil
}

func (e *extractor) lexFunctionDefinition(definition []string, lexedFunction *symbol.LexedFunctionSource) error {
	queries := map[string]string{
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

	e.scratch(definition)

	classAccess := "class"
	instanceAccess := "instance"

	for _, query := range queries {
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
				lexedFunction.Name = strings.Join(capture.Node.Text, "")
			case "access.class":
				lexedFunction.Access = &classAccess
			case "access.instance":
				lexedFunction.Access = &instanceAccess
			case "arg":
				lexedFunction.Args = append(lexedFunction.Args, symbol.LexedFunctionArgument{
					Name: strings.Join(capture.Node.Text, ""),
				})
			case "vararg":
				lexedFunction.Args = append(lexedFunction.Args, symbol.LexedFunctionArgument{
					Name: strings.Join(capture.Node.Text, ""),
				})
			}
		}

		return nil
	}

	return fmt.Errorf("Function source %v idn't yield any result", definition)
}

func NewExtractor() (*extractor, error) {
	nvimConfigDir, err := project.GetConfigDir()

	if err != nil {
		return nil, fmt.Errorf("Error getting nvim config location: %w", err)
	}

	client, err := nvim.New(nvim.NewConfig(nvimConfigDir))

	if err != nil {
		return nil, fmt.Errorf("Error initialising nvim client: %v", err)
	}

	return &extractor{
		nvim: client,
	}, nil
}
