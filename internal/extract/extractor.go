package extract

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/extract/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

type extractor struct {
	nvim *nvim.Nvim
}

func (e *extractor) Extract(path string) error {
	crawler := crawl.NewCrawler(e.nvim)

	_, err := crawler.CrawlRuntime(path)

	if err != nil {
		return fmt.Errorf("Error crawling %s: %w", path, err)
	}

	outputDir, err := project.GetOutputDir()

	if err != nil {
		return fmt.Errorf("Error getting output location: %w", err)
	}

	out := output.NewOutput(outputDir)

	if err := out.WriteFile("statistics.json", crawler.Statistics()); err != nil {
		return fmt.Errorf("Error collecting crawler statistics: %w", err)
	}

	if err := e.nvim.Quit(); err != nil {
		fmt.Println(fmt.Errorf("Error closing nvim gracefully: %w", err))
	}

	return nil
}

func NewExtractor() (*extractor, error) {
	nvimConfigDir, err := project.GetConfigDir()

	if err != nil {
		return nil, fmt.Errorf("Error getting nvim config location: %w", err)
	}

	nvimClient, err := nvim.New(nvim.NewConfig(nvimConfigDir))

	if err != nil {
		return nil, fmt.Errorf("Error initialising nvim client: %v", err)
	}

	err = nvimClient.Start()

	if err != nil {
		return nil, fmt.Errorf("Error opening nvim: %w", err)
	}

	extractor := extractor{
		nvim: nvimClient,
	}

	return &extractor, nil
}
