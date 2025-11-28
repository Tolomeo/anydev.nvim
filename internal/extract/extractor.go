package extract

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/extract/crawl"
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

func (e *extractor) serialize(source crawl.Source) {
	switch v := source.(type) {
	case *crawl.TableSource:
		fmt.Println(v, "table")
	case *crawl.FunctionSource:
		e.serializeFunction(v)
	case *crawl.VariableSource:
		fmt.Println("variable")
	}
	// fmt.Printf("%+v", source.Origin())
}

func (e *extractor) serializeFunction(source *crawl.FunctionSource) error {
	e.scratch(source.Origin().Definition)

	_, err := e.nvim.GetBufferLines(0, -1)

	if err != nil {
		return err
	}

	fmt.Println(source.Origin().Url)
	// fmt.Printf("%v", lines)

	/* function fn() end
	function fn(arg1) end
	function fn(arg1, arg2) end
	function fn(arg1, arg2, ...) end
		function fn(...) end */
	/* declarationQuery := `
		(function_declaration
			name: (identifier) @name
			parameters: (parameters
				(identifier)? @parameter
				("," (identifier) @parameter)*
				("," (vararg_expression) @parameter)?
				(vararg_expression)? @parameter
			)
		)
	` */

	/* function api:fn() end
	function api:fn(name) end
	function api:fn(name, value) end
	function api:fn(name, value, ...) end
	function api:fn(...) end */
	/* methodDeclarationQuery := `
		(function_declaration
			name: (method_index_expression
				method: (identifier) @function.name
			) @function.access
			parameters: (parameters
				(identifier)? @parameter
				("," (identifier) @parameter)*
				("," (vararg_expression) @parameter)?
				(vararg_expression)? @parameter
			)
		)
	` */

	/* function api.fn() end
	function api.fn(name) end
	function api.fn(name, value) end
	function api.fn(name, value, ...) end
	function api.fn(...) end */
	dotIndexDeclarationQuery := `
		(function_declaration
			name: (dot_index_expression
				field: (identifier) @name
			)
			parameters: (parameters
				(identifier)? @parameter
				("," (identifier) @parameter)*
				("," (vararg_expression) @parameter)?
				(vararg_expression)? @parameter
			)
		)
	`

	captures, err := e.nvim.ExecTsQuery(dotIndexDeclarationQuery)

	if err != nil {
		return err
	}

	fmt.Printf("%+v", captures)

	return nil
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
