package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getEnumeratorOriginType(enumeratorOrigin *origin.EnumeratorAnnotationOrigin) (*symbol.Enumerator, error) {
	enumerator := symbol.NewEnumerator(enumeratorOrigin.Name())

	for _, member := range enumeratorOrigin.Members() {
		enumeratorField := symbol.EnumeratorField{Name: member.Name().Text}
		value, err := tr.getType(member.Value().Text)

		if err != nil {
			return nil, err
		}

		enumeratorField.Value = value
		enumerator.Fields = append(enumerator.Fields, enumeratorField)
	}

	return enumerator, nil
}
