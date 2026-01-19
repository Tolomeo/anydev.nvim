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

func (l *Lexer) Lex() (symbol.Symbol, error) {
	target := l.context.Target()

	cacheId := []string{
		target.Origin().Url(),
		fmt.Sprintf("%d", target.Origin().Line()),
		fmt.Sprintf("%d", target.Origin().Character()),
	}

	if cachedSymbol, hasCachedSymbol := lexedCache.Get(cacheId...); hasCachedSymbol {
		fmt.Printf("\nUsing lexer cached result for symbol '%s': <%v> cache id hit \n", target.Identifier(), cacheId)
		l.context.Logger().Info(fmt.Sprintf("\nUsing lexer cached result for symbol '%s': <%v> cache id hit \n", target.Identifier(), cacheId))
		return cachedSymbol, nil
	}

	switch target.Origin().Type() {
	case treesitter.ASSIGNMENT_STATEMENT,
		treesitter.VARIABLE_DECLARATION,
		treesitter.FUNCTION_DECLARATION:
		lexed, err := l.lexValue()
		if err != nil {
			return nil, err
		}
		lexedCache.Set(lexed, cacheId...)
		return lexed, nil
	case treesitter.ALIAS_ANNOTATION:
		lexed, err := l.lexAliasType()
		if err != nil {
			return nil, err
		}
		lexedCache.Set(lexed, cacheId...)
		return lexed, nil
	case treesitter.CLASS_ANNOTATION:
		lexed, err := l.lexClassType()
		if err != nil {
			return nil, err
		}
		lexedCache.Set(lexed, cacheId...)
		return lexed, nil
	case treesitter.FIELD_ANNOTATION:
		lexed, err := l.lexFieldType()
		if err != nil {
			return nil, err
		}
		lexedCache.Set(lexed, cacheId...)
		return lexed, nil
	}

	return nil, fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", target.Identifier(), target.Origin())
}

func NewLexer(context context) *Lexer {
	return &Lexer{
		context: context,
	}
}
