package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
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
	Follow(string) (symbol.Type, error)
	// AddChild(*symbol.Table, string, symbol.Metadata, []string, symbol.Type)
}

type Transformer struct {
	context ctx
}

var lexedCache = cache.NewCache[symbol.Type]()

func (tr *Transformer) GetOriginType() (symbol.Type, error) {
	cacheId := []string{
		tr.context.Target().Origin().Url(),
		fmt.Sprintf("%d", tr.context.Target().Origin().Line()),
		fmt.Sprintf("%d", tr.context.Target().Origin().Character()),
	}

	if cachedSymbol, hasCachedSymbol := lexedCache.Get(cacheId...); hasCachedSymbol {
		tr.context.Logger().Infof("Using transformed cached result for symbol '%s': <%v> cache id hit", tr.context.Target().Identifier(), cacheId)
		return cachedSymbol, nil
	}

	targetOrigin := tr.context.Target().Origin()

	tr.context.Logger().Verbosef("Transforming target origin type <%T>", targetOrigin)

	switch targetOriginType := targetOrigin.(type) {
	case *origin.TableOrigin:
		tableType, err := tr.transformTableOrigin(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(tableType, cacheId...)
		return tableType, nil
	case *origin.FunctionOrigin:
		functionType, err := tr.getFunctionOriginType(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(functionType, cacheId...)
		return functionType, nil
	case *origin.ValueOrigin:
		valueType, err := tr.getValueOriginType(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(valueType, cacheId...)
		return valueType, nil
	case *origin.VirtualOrigin:
		metaType, err := tr.getVirtualOriginType(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(metaType, cacheId...)
		return metaType, nil
	case *origin.AliasOrigin:
		aliasType, err := tr.getAliasOriginType(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(aliasType, cacheId...)
		return aliasType, nil
	case *origin.AliasEnumeratorOrigin:
		aliasEnumeratorType, err := tr.getAliasEnumeratorType(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(aliasEnumeratorType, cacheId...)
		return aliasEnumeratorType, nil
	case *origin.ClassOrigin:
		classType, err := tr.getClassOriginType(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(classType, cacheId...)
		return classType, nil
	case *origin.FieldAnnotationOrigin:
		fieldType, err := tr.getFieldOriginType(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(fieldType, cacheId...)
		return fieldType, nil
	case *origin.EnumeratorAnnotationOrigin:
		enumSymbol, err := tr.getEnumeratorOriginType(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(enumSymbol, cacheId...)
		return enumSymbol, nil
	case *origin.VariableOrigin:
		variableSymbol, err := tr.getVariableOriginType(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(variableSymbol, cacheId...)
		return variableSymbol, nil
	case *origin.RequireFunctionCallOrigin:
		requiredModuleSymbol, err := tr.getModuleRequireOriginType(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(requiredModuleSymbol, cacheId...)
		return requiredModuleSymbol, nil
	case *origin.FunctionCallOrigin:
		originType, err := tr.transformFunctionCallOrigin(targetOriginType)
		if err != nil {
			return nil, err
		}
		lexedCache.Set(originType, cacheId...)
		return originType, nil
	case *origin.UnknownOrigin:
		return tr.getUnknownOriginType(targetOriginType), nil
	}

	/* switch ot := o.(type) {
	default:
		tr.context.Logger().Debugf("%T", ot)
	} */

	return nil, fmt.Errorf("Unknown origin type received for source '%s' <%T>", tr.context.Target().Identifier(), tr.context.Target().Origin())
}

func (tr *Transformer) GetType(typ string) (symbol.Type, error) {
	return tr.getType(typ)
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
