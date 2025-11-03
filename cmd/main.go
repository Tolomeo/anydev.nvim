package main

import (
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/crawler"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
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

	result, err := crawler.Crawl("vim")


	if err != nil {
		panic(fmt.Errorf("Error crawling %s: %w", "vim", err))
	}

	fmt.Printf("%v", result)

	if err := client.Quit(); err != nil {
		fmt.Println(fmt.Errorf("Error closing nvim gracefully: %w", err))
	}
}
