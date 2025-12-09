package lex

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/lex/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

type lexingContext struct {
	result  *result
	logs    *logs
}

type lexer struct {
	nvim    *nvim.Nvim
	context *lexingContext
}

func (l *lexer) Lex(paths ...string) error {
	err := l.nvim.Start()

	if err != nil {
		return fmt.Errorf("Error opening nvim: %w", err)
	}

	lexingContext := lexingContext{
		result: newResult(),
		logs:    newLogs(),
	}

	l.context = &lexingContext

	defer func() {
		l.context = nil
	}()

	for _, path := range paths {
		lexingContext.result.Runtime[path] = struct{}{}

		crawler := crawl.NewCrawler(l.nvim)

		source, err := crawler.CrawlRuntime(path)

		if err != nil {
			return fmt.Errorf("Error lexing %s: %w", path, err)
		}

		lexedSource, err := l.lex(source)

		if err != nil {
			return fmt.Errorf("Error lexing %s: %w", path, err)
		}

		lexingContext.result.Runtime[path] = lexedSource
	}

	outputDir, err := project.GetOutputDir()

	if err != nil {
		return fmt.Errorf("Error getting output location: %w", err)
	}

	out := output.NewOutput(outputDir)

	if err := out.WriteFile("result.json", lexingContext.result); err != nil {
		return fmt.Errorf("Error writing result.json: %w", err)
	}

	if err := out.WriteFile("logs.json", lexingContext.logs); err != nil {
		return fmt.Errorf("Error writing logs.json: %w", err)
	}

	err = l.nvim.Quit()

	if err != nil {
		return fmt.Errorf("Errot closing nvim process gracefully: %w", err)
	}

	return nil
}

func (l *lexer) lex(source crawl.Source) (*lexed.Symbol, error) {
	switch v := source.(type) {
	case *crawl.TableSource:
		fmt.Println(v, "table")
		return nil, nil
	case *crawl.FunctionSource:
		lexedFunction, err := l.lexFunction(v)

		if err != nil {
			return nil, err
		}

		return &lexedFunction, nil
	case *crawl.VariableSource:
		fmt.Println("variable")
		return nil, nil
	}

	return nil, fmt.Errorf("Error lexing '%s' source: unknown source type", source.Path())
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

	return &lexer{
		nvim: client,
	}, nil
}
