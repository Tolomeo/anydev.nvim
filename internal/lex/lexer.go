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
	Target() symbol.Target
	Extract(symbol.TargetKind, string) error
	ExtractChild(string) (symbol.Symbol, error)
}

type Lexer struct {
	context context
}

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
