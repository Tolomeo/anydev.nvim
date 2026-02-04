package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getClassOriginType(classOrigin *origin.ClassOrigin) (*symbol.Table, error) {
	class := symbol.NewTable()
	class.Name = classOrigin.Name()
	classFields, err := tr.context.Nvim().GetTypeCompletion(class.Name)

	if err != nil {
		return nil, err
	}

	for _, fieldName := range classFields {
		fieldSymbol, err := tr.context.ExtractChild(class, fieldName)

		if err != nil {
			return nil, err
		}

		class.Fields = append(class.Fields, *fieldSymbol)
	}

	return class, nil
}
