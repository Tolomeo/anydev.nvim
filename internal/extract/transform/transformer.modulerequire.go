package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getModuleRequireOriginType(requireFnCallOrigin *origin.RequireFunctionCallOrigin) (symbol.Type, error) {
	moduleName := requireFnCallOrigin.RequiredModuleName()
	err := tr.context.Extract("module", moduleName)

	if err != nil {
		return nil, err
	}

	return symbol.NewModuleReference(moduleName), nil
}
