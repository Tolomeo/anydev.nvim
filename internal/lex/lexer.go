package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/context"
	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/symbol"
)

type Lexer struct {
	crawler *crawl.Crawler
	context *context.Context
}

func (l *Lexer) Lex(paths []string, context *context.Context) error {
	l.context = context
	l.crawler = crawl.NewCrawler(l.context)
	defer func() {
		l.context = nil
		l.crawler = nil
	}()

	for _, path := range paths {
		if _, exists := l.context.Result().Runtime[path]; exists {
			return fmt.Errorf("Error lexing '%s': path already found", path)
		}

		l.context.Result().Runtime[path] = struct{}{}

		err := l.context.Provide(path, func(path string) error {
			source, err := l.source(path)

			if err != nil {
				return err
			}

			symbol, err := l.lex(source)

			if err != nil {
				return err
			}

			l.context.Result().Runtime[path] = symbol

			return nil
		})

		if err != nil {
			return fmt.Errorf("Error lexing %s: %w", path, err)
		}

	}

	return nil
}

func (l *Lexer) source(path string) (*crawl.Source, error) {
	source, err := l.crawler.Source(path)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing %s: %w", path, err)
	case source == nil:
		return nil, fmt.Errorf("Error lexing %s: No source found", path)
	}

	return source, nil
}

func (l *Lexer) lex(source *crawl.Source) (symbol.Symbol, error) {
	sourcePath := source.Path()
	sourceOrigin := source.Origin()

	if sourceOrigin == nil {
		symbol := newUnknownType()
		l.context.Logger.Warn(fmt.Sprintf("Using '%v' for symbol '%s' without origin", symbol, source.Path()))
		return symbol, nil
	}

	buffer, err := l.context.Nvim.Buffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	sourceDefinition := sourceOrigin.Definition()
	err = buffer.SetLines(sourceDefinition)

	if err != nil {
		return nil, err
	}

	table, err := l.matchTable(source)

	switch {
	case err != nil:
		return nil, err
	case table != nil:
		err = l.lexTable(table)

		if err != nil {
			return nil, err
		}

		return table, nil
	}

	function, err := l.matchFunction(source)

	switch {
	case err != nil:
		return nil, err
	case function != nil:
		err := l.lexFunction(function)

		if err != nil {
			return nil, err
		}

		return function, nil
	}

	meta, err := l.matchMeta(source)

	switch {
	case err != nil:
		return nil, err
	case meta != nil:
		symbol, err := l.lexMeta(meta)

		fmt.Println("Meta origin:")
		fmt.Printf("%+v", source.Origin())

		if err != nil {
			return nil, err
		}

		return symbol, nil
	}
	return nil, fmt.Errorf("Error lexing %s: unknown origin [%+v]", sourcePath, sourceOrigin)
}

func NewLexer() *Lexer {
	return &Lexer{}
}
