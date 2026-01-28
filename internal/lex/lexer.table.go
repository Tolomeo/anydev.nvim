package lex

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (l *Lexer) lexTableValue(tableOrigin *origin.TableOrigin) (*symbol.Table, error) {
	table := symbol.NewTable()
	table.Name = tableOrigin.Name()
	tableFields, err := l.context.Nvim().GetValueCompletion(l.context.Target().Identifier())

	if err != nil {
		return nil, err
	}

	for _, fieldName := range tableFields {
		err := l.context.ExtractChild(table, fieldName)

		if err != nil {
			return nil, err
		}
	}

	return table, nil
}
