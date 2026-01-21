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

type target interface {
	symbol.Target
	Nvim() *nvim.Nvim
	Logger() *log.Logger
	Extract(symbol.TargetKind, string) error
	ExtractChild(string) (symbol.Type, error)
}

type Lexer struct {
	target target
}

var lexedCache = cache.NewCache[symbol.Type]()

func (l *Lexer) Lex() (symbol.Type, error) {
	cacheId := []string{
		l.target.Origin().Url(),
		fmt.Sprintf("%d", l.target.Origin().Line()),
		fmt.Sprintf("%d", l.target.Origin().Character()),
	}

	if cachedSymbol, hasCachedSymbol := lexedCache.Get(cacheId...); hasCachedSymbol {
		fmt.Printf("\nUsing lexer cached result for symbol '%s': <%v> cache id hit \n", l.target.Identifier(), cacheId)
		l.target.Logger().Info(fmt.Sprintf("\nUsing lexer cached result for symbol '%s': <%v> cache id hit \n", l.target.Identifier(), cacheId))
		return cachedSymbol, nil
	}

	switch l.target.Origin().Type() {
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

	return nil, fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", l.target.Identifier(), l.target.Origin())
}

func NewLexer(context target) *Lexer {
	return &Lexer{
		target: context,
	}
}
