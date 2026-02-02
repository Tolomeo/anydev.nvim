package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
)

func (tr *Transformer) getValueOriginType(valueOrigin *origin.ValueOrigin) (annotation.Type, error) {
	valueType, err := tr.getType(valueOrigin.Type())

	if err != nil {
		return nil, err
	}

	return valueType, nil
}
