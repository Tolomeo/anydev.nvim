package main

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/crawler"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

func main() {
	client, err := nvim.New()

	if err != nil {
		panic(fmt.Errorf("Error initialising nvim client: %v", err))
	}

	err = client.Start()

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %w", err))
	}

	crawler := crawler.New(client)

	path := "vim.uv.fs_stat"

	_, err = crawler.Crawl(path)

	if err != nil {
		panic(fmt.Errorf("Error crawling %s: %w", path, err))
	}

	outputDir, err := project.GetOutputDir()

	if err != nil {
		panic(fmt.Errorf("Error getting output location: %w", err))
	}

	out := output.NewOutput(outputDir)

	err = out.WriteFile("statistics.json", crawler.Statistics())

	if err != nil {
		panic(fmt.Errorf("Error collecting crawler statistics: %w", err))
	}

	fmt.Printf("%+v", crawler.Statistics())

	if err := client.Quit(); err != nil {
		fmt.Println(fmt.Errorf("Error closing nvim gracefully: %w", err))
	}
}
