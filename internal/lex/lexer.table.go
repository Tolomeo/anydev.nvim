package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	// "github.com/Tolomeo/anydev.nvim/internal/nvim"
)

func (l *Lexer) lexTable(source *crawl.Source) (lexed.Symbol, error) {
	table := newTableType()

	fields, err := l.context.nvim.GetCompletion(source.Path())

	if err != nil {
		return nil, err
	}

	for _, fieldName := range fields {
		fieldPath := fmt.Sprintf("%s.%s", source.Path(), fieldName)
		fieldSymbol, err := l.context.provide(fieldPath, l.lex)

		switch {
		case err != nil:
			return nil, err
		case fieldSymbol == nil:
			return nil, fmt.Errorf("Error lexing table field '%s': no symbol found", fieldPath)
		}

		table.Fields = append(table.Fields, lexed.TableField{Name: fieldName, Value: fieldSymbol})
	}

	return &table, nil
}
