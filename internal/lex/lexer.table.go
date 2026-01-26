package lex

import (
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

func (l *Lexer) lexTableValue(origin *symbol.TableOrigin) (*symbol.Table, error) {
	table := symbol.NewTable()

	for _, capture := range origin.Captures() {
		switch capture.Id {
		case "table.name":
			table.Name = capture.Node.Text
		}
	}

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
