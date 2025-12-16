package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

var tableDeclarationQuery string = `
	(variable_declaration
		(assignment_statement
			(variable_list
				name: (identifier))
				(expression_list
					value: (table_constructor) @table
				)
		)
	)
`

func (l *Lexer) matchTable(source *crawl.Source) (*lexed.Table, error) {
	table := newTableType()

	definitionLines := strings.Split(source.Origin().Definition(), "\n")
	err := l.scratch(definitionLines)

	if err != nil {
		return nil, err
	}

	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "lua", Query: tableDeclarationQuery})

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	return table, nil
}

func (l *Lexer) lexTable(table *lexed.Table, source *crawl.Source) error {
	fields, err := l.context.nvim.GetCompletion(source.Path())

	if err != nil {
		return err
	}

	for _, fieldName := range fields {
		fieldPath := fmt.Sprintf("%s.%s", source.Path(), fieldName)
		fieldSymbol, err := l.context.provide(fieldPath, l.lex)

		switch {
		case err != nil:
			return err
		case fieldSymbol == nil:
			return fmt.Errorf("Error lexing table field '%s': no symbol found", fieldPath)
		}

		table.Fields = append(table.Fields, lexed.TableField{Name: fieldName, Value: fieldSymbol})
	}

	return nil
}
