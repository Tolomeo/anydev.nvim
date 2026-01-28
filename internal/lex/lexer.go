package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/cache"
	"github.com/Tolomeo/anydev.nvim/internal/log"
)

type ctx interface {
	Target() *target.Target
	Nvim() *nvim.Nvim
	Logger() *log.Logger
	Extract(target.TargetKind, string) error
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

	o := l.context.Target().Origin()

	switch ot := o.(type) {
	case *origin.TableOrigin:
		tableType, err := l.lexTableValue(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(tableType, cacheId...)
		return tableType, nil
	case *origin.FunctionOrigin:
		functionType, err := l.lexFunctionValue(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(functionType, cacheId...)
		return functionType, nil
	case *origin.MetaOrigin:
		metaType, err := l.lexMetaValue(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(metaType, cacheId...)
		return metaType, nil
	case *origin.AliasOrigin:
		aliasType, err := l.lexAliasType(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(aliasType, cacheId...)
		return aliasType, nil
	case *origin.AliasEnumeratorOrigin:
		aliasEnumeratorType, err := l.lexAliasEnumeratorType(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(aliasEnumeratorType, cacheId...)
		return aliasEnumeratorType, nil
	case *origin.ClassOrigin:
		classType, err := l.lexClassType(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(classType, cacheId...)
		return classType, nil
	case *origin.FieldOrigin:
		fieldType, err := l.lexFieldType(ot)
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
