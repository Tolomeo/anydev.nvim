package main

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

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

func getOutput() (*output.Output, error) {
	outputDir, err := project.GetOutputDir()

	if err != nil {
		return nil, fmt.Errorf("Error getting output location: %w", err)
	}

	out := output.NewOutput(outputDir)

	return out, nil
}

func main() {
	var paths []string = []string{"vim.deepcopy", "vim.validate"}

	client, err := getClient()

	if err != nil {
		panic(err)
	}

	logger := getLogger()
	crawler := getCrawler(client)
	lexer := getLexer()
	lexConfig := lex.NewLexingContext(logger, client, crawler)

	err = lexer.Lex(paths, lexConfig)

	if err != nil {
		panic(err)
	}

	out, err := getOutput()

	if err := out.WriteFile("result.json", lexConfig.Result()); err != nil {
		panic(fmt.Errorf("Error writing result.json: %w", err))
	}

	if err := out.WriteFile("logs.json", logger.Logs()); err != nil {
		panic(fmt.Errorf("Error writing logs.json: %w", err))
	}

	/* err = client.Quit()

	if err != nil {
		return panic(err)
	} */
}
