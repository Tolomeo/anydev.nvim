package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
)

func (tr *Transformer) getModuleRequireOriginType(_ *origin.RequireFunctionCallOrigin) (annotation.Type, error) {
	return tr.context.Follow()
}
