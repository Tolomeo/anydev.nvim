package lex

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
)

type Lexer struct {
	context *lexingContext
}

func (l *Lexer) Lex(paths []string, context *lexingContext) error {
	l.context = context
	defer func() { l.context = nil }()

	for _, path := range paths {
		if _, exists := l.context.result.Runtime[path]; exists {
			return fmt.Errorf("Error lexing '%s': path already found", path)
		}

		l.context.Result().Runtime[path] = struct{}{}

		err := l.context.provide(path, func(path string) error {
			source, err := l.source(path)

			if err != nil {
				return err
			}

			symbol, err := l.lex(source)

			if err != nil {
				return err
			}

			l.context.result.Runtime[path] = symbol

			return nil
		})

		if err != nil {
			return fmt.Errorf("Error lexing %s: %w", path, err)
		}

	}

	return nil
}

func (l *Lexer) source(path string) (*crawl.Source, error) {
	source, err := l.context.crawler.SourceRuntime(path)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing %s: %w", path, err)
	case source == nil:
		return nil, fmt.Errorf("Error lexing %s: No source found", path)
	}

	return source, nil
}

func (l *Lexer) lex(source *crawl.Source) (lexed.Symbol, error) {
	sourcePath := source.Path()
	sourceOrigin := source.Origin()

	if sourceOrigin == nil {
		symbol := newUnknownType()
		l.context.logger.Warn(fmt.Sprintf("Using '%v' for symbol '%s' without origin", symbol, source.Path()))
		return symbol, nil
	}

	sourceType := sourceOrigin.Type()
	sourceDefinition := sourceOrigin.Definition()

	err := l.scratch(sourceDefinition)

	if err != nil {
		return nil, err
	}

	switch sourceType {
	case "variable_declaration":
		symbol, err := l.lexVariableDeclaration(source)

		switch {
		case err != nil:
			return nil, err
		case symbol == nil:
			return nil, fmt.Errorf("Error lexing variable declaration: no symbol found in '%s'", sourceDefinition)
		}

		return symbol, nil

	case "function_declaration":
		symbol, err := l.lexFunctionDeclaration(source)

		switch {
		case err != nil:
			return nil, err
		case symbol == nil:
			return nil, fmt.Errorf("Error lexing variable declaration: no symbol found in '%s'", sourceDefinition)
		}

		return symbol, nil
	default:
		return nil, fmt.Errorf("Error lexing %s: Unkonw origin type '%s'", sourcePath, sourceType)
	}
}

func (l *Lexer) lexVariableDeclaration(source *crawl.Source) (lexed.Symbol, error) {
	table, err := l.matchTable(source)

	switch {
	case err != nil:
		return nil, err
	case table == nil:
		return nil, nil
	}

	err = l.lexTable(table)

	if err != nil {
		return nil, err
	}

	return table, nil
}

func (l *Lexer) lexFunctionDeclaration(source *crawl.Source) (lexed.Symbol, error) {
	function, err := l.matchFunction(source)

	if err != nil {
		return nil, err
	}

	err = l.lexFunction(function)

	if err != nil {
		return nil, err
	}

	return function, nil
}

func (l *Lexer) scratch(lines []string) error {
	buffer := path.Join(l.context.nvim.Options().Config().Dir(), "anydev.lexer.lua")

	_, err := l.context.nvim.Open(buffer)

	if err != nil {
		return fmt.Errorf("Error writing to scratch buffer: %w", err)
	}

	err = l.context.nvim.SetBufferLines(lines)

	if err != nil {
		return fmt.Errorf("Error writing to scratch buffer: %w", err)
	}

	return nil
}

func NewLexer() *Lexer {
	return &Lexer{}
}
