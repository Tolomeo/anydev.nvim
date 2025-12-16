package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

var tableQueries = map[string]string{
	"tableDeclaration": `
	(variable_declaration
		(assignment_statement
			(variable_list
				name: (identifier)
			) @table.name
			(expression_list
				value: (table_constructor)
			) @table.value
		)
	)
`,
	"tableFieldDeclaration": `
	(assignment_statement
		(variable_list
			name: (dot_index_expression
				table: (_)
				field: (identifier) @table.name
			)
		)
		(expression_list
			value: (table_constructor) @table.value
		)
	)
`,
	"tableIndexFieldDeclaration": `
		(assignment_statement
			(variable_list
				name: (bracket_index_expression
					table: (_)
					field: (string
						content: (string_content) @table.name
					)
				)
			)
			(expression_list
				value: (table_constructor) @table.value
			)
		)
`}

func (l *Lexer) matchTable(source *crawl.Source) (*lexed.Table, error) {
	definitionLines := strings.Split(source.Origin().Definition(), "\n")
	err := l.scratch(definitionLines)

	if err != nil {
		return nil, err
	}

	for _, query := range tableQueries {
		captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "lua", Query: query})

		switch {
		case err != nil:
			return nil, err
		case captures == nil:
			continue
		}

		table := newTableType()

		for _, capture := range *captures {
			switch capture.Id {
			case "table.name":
				table.Name = &capture.Node.Text
			}
		}

		return table, nil
	}

	return nil, nil
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
