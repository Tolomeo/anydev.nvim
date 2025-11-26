package extract

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/extract/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

type extractor struct{}

func (e *extractor) Extract(path string) error {
	nvimConfigDir, err := project.GetConfigDir()

	if err != nil {
		return fmt.Errorf("Error getting nvim config location: %w", err)
	}

	nvimClient, err := nvim.New(nvim.NewConfig(nvimConfigDir))

	if err != nil {
		return fmt.Errorf("Error initialising nvim client: %v", err)
	}

	err = nvimClient.Start()

	if err != nil {
		return fmt.Errorf("Error opening nvim: %w", err)
	}

	crawler := crawl.NewCrawler(nvimClient)

	source, err := crawler.CrawlRuntime(path)

	if err != nil {
		return fmt.Errorf("Error crawling %s: %w", path, err)
	}

	e.transform(source)

	outputDir, err := project.GetOutputDir()

	if err != nil {
		return fmt.Errorf("Error getting output location: %w", err)
	}

	out := output.NewOutput(outputDir)

	if err := out.WriteFile("statistics.json", crawler.Statistics()); err != nil {
		return fmt.Errorf("Error collecting crawler statistics: %w", err)
	}

	return nil
}

func (e *extractor) transform(source crawl.Source) {
	switch v := source.(type) {
	case *crawl.TableSource:
		fmt.Println(v, "table")
	case *crawl.FunctionSource:
		fmt.Println("function")
	case *crawl.VariableSource:
		fmt.Println("variable")
	}
	fmt.Printf("%+v", source.Origin())
}

func NewExtractor() *extractor {
	return &extractor{}
}
