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
		source, err := l.sourceValue(path)

		if err != nil {
			return err
		}

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

func (l *Lexer) lexValue(source *symbol.ValueSource) (symbol.Symbol, error) {
	sourcePath := source.Path
	sourceOrigin := source.Origin

	if sourceOrigin == nil {
		symbol := newUnknownType()
		l.context.Logger.Warn(fmt.Sprintf("Using '%v' for symbol '%s' without origin", symbol, source.Path))
		return symbol, nil
	}

	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	sourceDefinition := sourceOrigin.DefinitionLines()
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

		/* fmt.Println("Meta origin:")
		fmt.Printf("%+v", source.Origin) */

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
