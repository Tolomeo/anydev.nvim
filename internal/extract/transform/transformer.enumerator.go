package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getEnumeratorOriginType(enumeratorOrigin *origin.EnumOrigin) (*symbol.Enumerator, error) {
	enumerator := symbol.NewEnumerator(enumeratorOrigin.Name())

	return enumerator, nil
}
