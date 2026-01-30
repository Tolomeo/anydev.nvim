package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/cache"
)

type ctx interface {
	Target() *target.Target
	Nvim() *nvim.Nvim
	Logger() *log.Logger
	Extract(target.TargetKind, string) error
	ExtractChild(*symbol.Table, string) error
}

type Transformer struct {
	context ctx
}

var lexedCache = cache.NewCache[annotation.Type]()

func (tr *Transformer) GetType() (annotation.Type, error) {
	cacheId := []string{
		tr.context.Target().Origin().Url(),
		fmt.Sprintf("%d", tr.context.Target().Origin().Line()),
		fmt.Sprintf("%d", tr.context.Target().Origin().Character()),
	}

	if cachedSymbol, hasCachedSymbol := lexedCache.Get(cacheId...); hasCachedSymbol {
		tr.context.Logger().Infof("Using transformed cached result for symbol '%s': <%v> cache id hit", tr.context.Target().Identifier(), cacheId)
		return cachedSymbol, nil
	}

	o := tr.context.Target().Origin()

	switch ot := o.(type) {
	case *origin.TableOrigin:
		tableType, err := tr.getTableOriginType(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(tableType, cacheId...)
		return tableType, nil
	case *origin.FunctionOrigin:
		functionType, err := tr.getFunctionOriginType(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(functionType, cacheId...)
		return functionType, nil
	case *origin.VirtualOrigin:
		metaType, err := tr.getVirtualOriginType(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(metaType, cacheId...)
		return metaType, nil
	case *origin.AliasOrigin:
		aliasType, err := tr.getAliasOriginType(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(aliasType, cacheId...)
		return aliasType, nil
	case *origin.AliasEnumeratorOrigin:
		aliasEnumeratorType, err := tr.lexAliasEnumeratorType(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(aliasEnumeratorType, cacheId...)
		return aliasEnumeratorType, nil
	case *origin.ClassOrigin:
		classType, err := tr.getClassOriginType(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(classType, cacheId...)
		return classType, nil
	case *origin.FieldOrigin:
		fieldType, err := tr.getFieldOriginType(ot)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(fieldType, cacheId...)
		return fieldType, nil
	}

	return nil, fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", tr.context.Target().Identifier(), tr.context.Target().Origin())
}

func (tr *Transformer) GetMetadata() (*symbol.Metadata, error) {
	origin := tr.context.Target().OriginChain().First()

	return tr.getMetadata(origin)
}

func NewTransformer(context ctx) *Transformer {
	return &Transformer{
		context: context,
	}
}
