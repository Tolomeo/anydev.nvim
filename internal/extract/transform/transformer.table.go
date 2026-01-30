package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getTableOriginType(tableOrigin *origin.TableOrigin) (*symbol.Table, error) {
	table := symbol.NewTable()
	table.Name = tableOrigin.Name()
	tableFields, err := tr.context.Nvim().GetValueCompletion(tr.context.Target().Identifier())

	if err != nil {
		return nil, err
	}

	for _, fieldName := range tableFields {
		err := tr.context.ExtractChild(table, fieldName)

		if err != nil {
			return nil, err
		}
	}

	return table, nil
}
