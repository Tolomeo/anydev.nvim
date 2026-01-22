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

type ctx interface {
	Target() *symbol.Target
	Nvim() *nvim.Nvim
	Logger() *log.Logger
	Extract(symbol.TargetKind, string) error
	ExtractChild(*symbol.Table, string) error
}

type Lexer struct {
	context ctx
}

var lexedCache = cache.NewCache[symbol.Type]()

func (l *Lexer) Lex() (symbol.Type, error) {
	cacheId := []string{
		l.context.Target().Origin().Url(),
		fmt.Sprintf("%d", l.context.Target().Origin().Line()),
		fmt.Sprintf("%d", l.context.Target().Origin().Character()),
	}

	if cachedSymbol, hasCachedSymbol := lexedCache.Get(cacheId...); hasCachedSymbol {
		l.context.Logger().Infof("Using lexer cached result for symbol '%s': <%v> cache id hit", l.context.Target().Identifier(), cacheId)
		return cachedSymbol, nil
	}

	switch l.context.Target().Origin().Type() {
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

	return nil, fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", l.context.Target().Identifier(), l.context.Target().Origin())
}

func NewLexer(context ctx) *Lexer {
	return &Lexer{
		context: context,
	}
}
