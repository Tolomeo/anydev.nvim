package main

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/crawler"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

const path string = "vim.fs"

func main() {
	nvimConfigDir, err := project.GetConfigDir()

	if err != nil {
		panic(fmt.Errorf("Error getting nvim config location: %w", err))
	}

	client, err := nvim.New(nvim.NewConfig(nvimConfigDir))

	if err != nil {
		panic(fmt.Errorf("Error initialising nvim client: %v", err))
	}

	err = client.Start()

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %w", err))
	}

	crawler := crawler.New(client)

	_, err = crawler.Crawl(path)

	if err != nil {
		panic(fmt.Errorf("Error crawling %s: %w", path, err))
	}

	outputDir, err := project.GetOutputDir()

	if err != nil {
		panic(fmt.Errorf("Error getting output location: %w", err))
	}

	out := output.NewOutput(outputDir)

	if err := out.WriteFile("statistics.json", crawler.Statistics()); err != nil {
		panic(fmt.Errorf("Error collecting crawler statistics: %w", err))
	}

	if err := client.Quit(); err != nil {
		fmt.Println(fmt.Errorf("Error closing nvim gracefully: %w", err))
	}
}
