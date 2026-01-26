package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
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

	origin := l.context.Target().Origin()

	switch o := origin.(type) {
	case *symbol.TableOrigin:
		tableType, err := l.lexTableValue(o)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(tableType, cacheId...)
		return tableType, nil
	case *symbol.FunctionOrigin:
		functionType, err := l.lexFunctionValue(o)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(functionType, cacheId...)
		return functionType, nil
	case *symbol.MetaOrigin:
		metaType, err := l.lexMetaValue(o)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(metaType, cacheId...)
		return metaType, nil
	case *symbol.AliasOrigin:
		aliasType, err := l.lexAliasType(o)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(aliasType, cacheId...)
		return aliasType, nil
	case *symbol.AliasEnumeratorOrigin:
		aliasEnumeratorType, err := l.lexAliasEnumeratorType(o)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(aliasEnumeratorType, cacheId...)
		return aliasEnumeratorType, nil
	case *symbol.ClassOrigin:
		classType, err := l.lexClassType(o)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(classType, cacheId...)
		return classType, nil
	case *symbol.FieldOrigin:
		fieldType, err := l.lexFieldType(o)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(fieldType, cacheId...)
		return fieldType, nil
	}

	return nil, fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", l.context.Target().Identifier(), l.context.Target().Origin())
}

func NewLexer(context ctx) *Lexer {
	return &Lexer{
		context: context,
	}
}
