package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/context"
	"github.com/Tolomeo/anydev.nvim/internal/lex/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

type Lexer struct {
	crawler *crawl.Crawler
	context *context.Context
}

func (l *Lexer) LexValue(path string, context *context.Context) error {
	currentContext := l.context
	currentCrawler := l.crawler
	l.context = context
	l.crawler = crawl.NewCrawler(l.context)
	defer func() {
		l.context = currentContext
		l.crawler = currentCrawler
	}()

	if _, exists := l.context.Result().Runtime[path]; exists {
		l.context.Logger.Info("Skipping '%s': lexed symbol already found")
		return nil
	}

	l.context.Result().Runtime[path] = struct{}{}

	err := l.context.Push(path, func(path string) error {
		fmt.Println("Lexing:", l.context.Current())

		source, err := l.sourceValue(path)

		if err != nil {
			return err
		}

		// fmt.Printf("\n%+v\n", source.Origin.DocumentationLines())

		symbol, err := l.lexValue(source)

		if err != nil {
			return err
		}

		l.context.Result().Runtime[path] = symbol

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error lexing %s: %w", path, err)
	}

	return nil
}

func (l *Lexer) sourceValue(path string) (*symbol.ValueSource, error) {
	source, err := l.crawler.SourceValue(path)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing %s: %w", path, err)
	case source == nil:
		return nil, fmt.Errorf("Error lexing %s: No source found", path)
	}

	return source, nil
}

func NewLexer() *Lexer {
	return &Lexer{}
}
