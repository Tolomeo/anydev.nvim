package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
)

func (tr *Transformer) getFieldOriginType(fieldOrigin *origin.FieldAnnotationOrigin) (symbol.Type, error) {
	lexedField, err := tr.getType(fieldOrigin.Type())

	if err != nil {
		return nil, err
	}

	return lexedField, nil
}
