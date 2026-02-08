package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getUnknownOriginType(variableOrigin *origin.UnknownOrigin) (*symbol.Unknown) {
	return symbol.NewUnknown()
}
