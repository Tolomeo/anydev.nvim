package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/context"
	"github.com/Tolomeo/anydev.nvim/internal/lex/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
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

	l.context.Result().Runtime[path] = symbol.NewUnknown()

	err := l.context.Push(path, func(path string) error {
		fmt.Printf("\nLexing: %s value\n", l.context.Current())

		source, err := l.crawler.SourceValue(path)

		if err != nil {
			return err
		}

		// fmt.Printf("\n%+v\n", source.Origin.DocumentationLines())

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

	return nil
}

func (l *Lexer) LexType(name string, context *context.Context) error {
	currentContext := l.context
	currentCrawler := l.crawler
	l.context = context
	l.crawler = crawl.NewCrawler(l.context)
	defer func() {
		l.context = currentContext
		l.crawler = currentCrawler
	}()

	if _, alreadyLexed := l.context.Result().Types[name]; alreadyLexed {
		l.context.Logger.Info(fmt.Sprintf("Skipping '%s': lexed type already found", name))
		return nil
	}

	l.context.Result().Types[name] = symbol.NewUnknown()

	err := l.context.Fork(name, func(name string) error {
		fmt.Printf("\nLexing: %s type\n", l.context.Current())

		source, err := l.crawler.SourceType(name, "")

		switch {
		case err != nil:
			return err
		case source == nil:
			return nil
		}

		/* fmt.Printf("\nReference '%s' source:\n%+v\n\n", name, source.Origin)
		fmt.Printf("\nType: %+v\n\n", source.Origin.DefinitionText()) */

		lexedAliasType, err := l.lexAliasType(source)

		// fmt.Printf("\nLexed '%s' alias: %+v\n\n", name, lexedAliasType)

		switch {
		case err != nil:
			return err
		case lexedAliasType != nil:
			// fmt.Printf("\nLexed '%s' alias: %+v\n\n", name, lexedAliasType)
			l.context.Result().Types[name] = lexedAliasType
			return nil
		}

		lexedClassType, err := l.lexClassType(source)

		switch {
		case err != nil:
			return err
		case lexedClassType != nil:
			// fmt.Printf("\nLexed '%s' class: %+v\n\n", name, lexedClassType)
			l.context.Result().Types[name] = lexedClassType
			return nil
		}

		fmt.Printf(fmt.Sprintf("No types found for '%s' name", name))
		l.context.Logger.Warn(fmt.Sprintf("No type definitions found for '%s' name", name))
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (l *Lexer) lex(source symbol.Source) (symbol.Symbol, error) {
	switch source.GetOrigin().Type() {
	case treesitter.ASSIGNMENT_STATEMENT,
		treesitter.VARIABLE_DECLARATION,
		treesitter.FUNCTION_DECLARATION:
		return l.lexValue(source)
	case treesitter.ALIAS_ANNOTATION:
		return l.lexAliasType(source)
	case treesitter.CLASS_ANNOTATION:
		return l.lexClassType(source)
	case treesitter.FIELD_ANNOTATION:
		return l.lexFieldType(source)
	}

	return nil, fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", source.Identifier(), source)
}

// TODO: remove
func (l *Lexer) sourceValue(path string) (symbol.Source, error) {
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
