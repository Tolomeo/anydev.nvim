package lex

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/lex/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

type lexer struct {
	nvim   *nvim.Nvim
	buffer string
}

func (l *lexer) Lex(path string) error {
	err := l.nvim.Start()

	if err != nil {
		return fmt.Errorf("Error opening nvim: %w", err)
	}

	crawler := crawl.NewCrawler(l.nvim)

	source, err := crawler.CrawlRuntime(path)

	if err != nil {
		return fmt.Errorf("Error lexing %s: %w", path, err)
	}

	err = l.lex(source)

	if err != nil {
		return fmt.Errorf("Error lexing %s: %w", path, err)
	}

	outputDir, err := project.GetOutputDir()

	if err != nil {
		return fmt.Errorf("Error getting output location: %w", err)
	}

	out := output.NewOutput(outputDir)

	if err := out.WriteFile("statistics.json", crawler.Statistics()); err != nil {
		return fmt.Errorf("Error collecting crawler statistics: %w", err)
	}

	err = l.nvim.Quit()

	if err != nil {
		return fmt.Errorf("Errot closing nvim process gracefully: %w", err)
	}

	return nil
}

func (l *lexer) scratch(lines []string) error {
	buffer := path.Join(l.nvim.Options().Config().Dir(), "anydev.extractor.lua")

	_, err := l.nvim.Open(buffer)

	if err != nil {
		return err
	}

	err = l.nvim.SetBufferLines(lines)

	if err != nil {
		return err
	}

	return nil
}

func (l *lexer) lex(source crawl.Source) error {
	switch v := source.(type) {
	case *crawl.TableSource:
		fmt.Println(v, "table")
	case *crawl.FunctionSource:
		return l.lexFunction(v)
	case *crawl.VariableSource:
		fmt.Println("variable")
	}
	// fmt.Printf("%+v", source.Origin())
	return nil
}

func NewLexer() (*lexer, error) {
	nvimConfigDir, err := project.GetConfigDir()

	if err != nil {
		return nil, fmt.Errorf("Error getting nvim config location: %w", err)
	}

	client, err := nvim.New(nvim.NewConfig(nvimConfigDir))

	if err != nil {
		return nil, fmt.Errorf("Error initialising nvim client: %v", err)
	}

	return &lexer{
		nvim: client,
	}, nil
}
