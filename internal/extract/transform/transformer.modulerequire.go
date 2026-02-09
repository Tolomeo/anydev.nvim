package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getModuleRequireOriginType(_ *origin.RequireFunctionCallOrigin) (symbol.Type, error) {
	return tr.context.Follow()
}
