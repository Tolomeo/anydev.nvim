package main

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

var paths []string = []string{"vim.deepcopy", "vim.validate"}

func getClient() (*nvim.Nvim, error) {
	nvimConfigDir, err := project.GetConfigDir()

	if err != nil {
		return nil, fmt.Errorf("Error getting nvim config location: %w", err)
	}

	client, err := nvim.New(nvim.NewConfig(nvimConfigDir))

	if err != nil {
		return nil, fmt.Errorf("Error initialising nvim client: %v", err)
	}

	err = client.Start()

	if err != nil {
		return nil, fmt.Errorf("Error starting nvim client: %w", err)
	}

	return client, nil
}

func getLogger() *log.Logger {
	return log.NewLogger("")
}

func getCrawler(client *nvim.Nvim) *crawl.Crawler {
	return crawl.NewCrawler(crawl.CrawlerOptions{
		Nvim: client,
		Log:  func(message string) {},
	})
}

func getLexer() *lex.Lexer {
	return lex.NewLexer()
}

func main() {
	logger := getLogger()
	client, err := getClient()
	crawler := getCrawler(client)
	lexer := getLexer()

	if err != nil {
		panic(err)
	}

	lexConfig := lex.NewLexingContext(logger, client, crawler)
	err = lexer.Lex(paths, lexConfig)

	if err != nil {
		panic(err)
	}
}
