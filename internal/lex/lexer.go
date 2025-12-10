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
	nvim    *nvim.Nvim
	context *lexingContext
	crawler *crawl.Crawler
}

func (l *lexer) Lex(paths ...string) error {
	l.context = newLexingContext()
	defer func() { l.context = nil }()

	l.crawler = crawl.NewCrawler(crawl.CrawlerOptions{
		Nvim: l.nvim,
		Log:  l.context.Info,
	})
	defer func() { l.crawler = nil }()

	for _, path := range paths {
		err := l.context.Provide(path, l.lex)

		if err != nil {
			return fmt.Errorf("Error lexing %s: %w", path, err)
		}
	}

	outputDir, err := project.GetOutputDir()

	if err != nil {
		return fmt.Errorf("Error getting output location: %w", err)
	}

	out := output.NewOutput(outputDir)

	if err := out.WriteFile("result.json", l.context.Result()); err != nil {
		return fmt.Errorf("Error writing result.json: %w", err)
	}

	if err := out.WriteFile("logs.json", l.context.Logs()); err != nil {
		return fmt.Errorf("Error writing logs.json: %w", err)
	}

	err = l.nvim.Quit()

	if err != nil {
		return fmt.Errorf("Errot closing nvim process gracefully: %w", err)
	}

	return nil
}

func (l *lexer) lex(path string) error {
	l.context.Result().Runtime[path] = struct{}{}

	source, err := l.crawler.CrawlRuntime(path)

	if err != nil {
		return fmt.Errorf("Error lexing %s: %w", path, err)
	}

	switch v := source.(type) {
	case *crawl.TableSource:
		fmt.Println(v, "table")
		l.context.Result().Runtime[path] = struct{}{}
	case *crawl.FunctionSource:
		lexedFunction, err := l.lexFunction(v)

		if err != nil {
			return err
		}

		l.context.Result().Runtime[path] = lexedFunction
		return nil
	case *crawl.VariableSource:
		fmt.Println("variable")
		l.context.Result().Runtime[path] = struct{}{}
	default:
		return fmt.Errorf("Error lexing '%s' source: unknown source type", path)
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

func NewLexer() (*lexer, error) {
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

	return &lexer{
		nvim: client,
	}, nil
}
