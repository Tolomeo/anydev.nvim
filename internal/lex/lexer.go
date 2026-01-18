package lex

import (
	"fmt"

	// "github.com/Tolomeo/anydev.nvim/internal/lex/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/cache"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type context interface {
	Nvim() *nvim.Nvim
	Logger() *log.Logger
	CurrentName() string
}

type Lexer struct {
	context context
}

/* func (l *Lexer) LexValue(path string, context *context.Context) error {
	currentContext := l.context
	currentCrawler := l.crawler
	l.context = context
	l.crawler = crawl.NewCrawler(l.context)
	defer func() {
		l.context = currentContext
		l.crawler = currentCrawler
	}()

	if _, exists := l.context.Result().Runtime[path]; exists {
		l.context.Logger().Info("Skipping '%s': lexed symbol already found")
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

		symbol, err := l.Lex(source)

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
*/

/* func (l *Lexer) LexType(name string, context *context.Context) error {
	currentContext := l.context
	currentCrawler := l.crawler
	l.context = context
	l.crawler = crawl.NewCrawler(l.context)
	defer func() {
		l.context = currentContext
		l.crawler = currentCrawler
	}()

	if _, alreadyLexed := l.context.Result().Types[name]; alreadyLexed {
		l.context.Logger().Info(fmt.Sprintf("Skipping '%s': lexed type already found", name))
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

		// fmt.Printf("\nReference '%s' source:\n%+v\n\n", name, source.Origin)
		// fmt.Printf("\nType: %+v\n\n", source.Origin.DefinitionText())

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

		l.context.Logger().Warn(fmt.Sprintf("No type definitions found for '%s' name", name))
		return nil
	})

	if err != nil {
		return err
	}

	return nil
} */

var lexedCache = cache.NewCache[symbol.Symbol]()

func (l *Lexer) Lex(source symbol.Source) (symbol.Symbol, error) {
	cacheId := []string{
		source.GetOrigin().Url(),
		fmt.Sprintf("%d", source.GetOrigin().Line()),
		fmt.Sprintf("%d", source.GetOrigin().Character()),
	}

	if cachedSymbol, hasCachedSymbol := lexedCache.Get(cacheId...); hasCachedSymbol {
		fmt.Printf("\nUsing lexer cached result for symbol '%s': <%v> cache id hit \n", source.Identifier(), cacheId)
		l.context.Logger().Info(fmt.Sprintf("\nUsing lexer cached result for symbol '%s': <%v> cache id hit \n", source.Identifier(), cacheId))
		return cachedSymbol, nil
	}

	switch source.GetOrigin().Type() {
	case treesitter.ASSIGNMENT_STATEMENT,
		treesitter.VARIABLE_DECLARATION,
		treesitter.FUNCTION_DECLARATION:
		lexed, err := l.lexValue(source)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(lexed, cacheId...)
		return lexed, nil
	case treesitter.ALIAS_ANNOTATION:
		lexed, err := l.lexAliasType(source)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(lexed, cacheId...)
		return lexed, nil
	case treesitter.CLASS_ANNOTATION:
		lexed, err := l.lexClassType(source)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(lexed, cacheId...)
		return lexed, nil
	case treesitter.FIELD_ANNOTATION:
		lexed, err := l.lexFieldType(source)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(lexed, cacheId...)
		return lexed, nil
	}

	return nil, fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", source.Identifier(), source)
}

func NewLexer(context context) *Lexer {
	return &Lexer{
		context: context,
	}
}
