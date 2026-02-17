package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getClassOriginType(classOrigin *origin.ClassOrigin) (*symbol.Table, error) {
	class := symbol.NewTable()
	class.Name = classOrigin.Name()
	classChildren, err := tr.context.Nvim().GetTypeCompletion(class.Name)

	if err != nil {
		return nil, err
	}

	for _, childName := range classChildren {
		err := tr.context.ExtractChild(class, childName)

		if err != nil {
			return nil, err
		}
	}

	return class, nil
}
