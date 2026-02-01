package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
)

func (tr *Transformer) getFieldOriginType(fieldOrigin *origin.FieldOrigin) (annotation.Type, error) {
	lexedField, err := tr.getType(fieldOrigin.Type())

	if err != nil {
		return nil, err
	}

	return lexedField, nil
}
