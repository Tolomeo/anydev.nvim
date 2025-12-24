package lex

import (
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
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

func (l *Lexer) matchTable(source *symbol.Source) (*symbol.Table, error) {
	buffer, err := l.context.Nvim.Buffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	err = buffer.SetLines(source.Origin.Definition())

	if err != nil {
		return nil, err
	}

	for _, query := range tableQueries {
		captures, err := l.context.Nvim.TsQueryOne(nvim.TsQueryConfig{Language: "lua", Query: query})

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

func (l *Lexer) lexTable(table *symbol.Table) error {
	tablePath := l.context.Current()
	tableFields, err := l.context.Nvim.GetCompletion(tablePath)

	if err != nil {
		return err
	}

	for _, fieldName := range tableFields {
		tableField := symbol.TableField{Name: fieldName}

		err := l.context.Provide(fieldName, func(path string) error {
			source, err := l.source(l.context.Current())

			if err != nil {
				return err
			}

			annotations, err := l.lexAnnotations(source.Origin.Documentation)

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
