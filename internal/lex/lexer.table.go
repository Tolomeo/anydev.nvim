package lex

import (
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
	buffer, err := l.context.nvim.Buffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	err = buffer.Edit(source.Origin().Definition())

	if err != nil {
		return nil, err
	}

	for _, query := range tableQueries {
		captures, err := l.context.nvim.TsQueryOne(nvim.TsQueryConfig{Language: "lua", Query: query})

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

func (l *Lexer) lexTable(table *lexed.Table) error {
	tablePath := l.context.current()
	tableFields, err := l.context.nvim.GetCompletion(tablePath)

	if err != nil {
		return err
	}

	for _, fieldName := range tableFields {
		tableField := lexed.TableField{Name: fieldName}

		err := l.context.provide(fieldName, func(path string) error {
			source, err := l.source(l.context.current())

			if err != nil {
				return err
			}

			annotations, err := l.lexAnnotations(source.Origin().Documentation())

			tableField.Private = annotations.private
			tableField.Protected = annotations.protected
			tableFieldValue, err := l.lex(source)

			if err != nil {
				return err
			}

			tableField.Value = tableFieldValue

			return nil
		})

		if err != nil {
			return err
		}

		table.Fields = append(table.Fields, tableField)
	}

	return nil
}
