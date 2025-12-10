package lex

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

type Lexer struct {
	context *lexingContext
}

func (l *Lexer) Lex(paths []string, context *lexingContext) error {
	l.context = context
	defer func() { l.context = nil }()

	for _, path := range paths {
		err := l.context.provide(path, l.lex)

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

	if err := out.WriteFile("logs.json", l.context.logger.Logs()); err != nil {
		return fmt.Errorf("Error writing logs.json: %w", err)
	}

	/* err = l.nvim.Quit()

	if err != nil {
		return fmt.Errorf("Errot closing nvim process gracefully: %w", err)
	} */

	return nil
}

func (l *Lexer) lex(path string) error {
	l.context.Result().Runtime[path] = struct{}{}

	source, err := l.context.crawler.CrawlRuntime(path)

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

func (l *Lexer) scratch(lines []string) error {
	buffer := path.Join(l.context.nvim.Options().Config().Dir(), "anydev.extractor.lua")

	_, err := l.context.nvim.Open(buffer)

	if err != nil {
		return err
	}

	err = l.context.nvim.SetBufferLines(lines)

	if err != nil {
		return err
	}

	return nil
}

func NewLexer() *Lexer {
	return &Lexer{}
}
