package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var metaQuery string = `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: [
				(vararg_expression) @assignment.right
			] 
		)
	) @assignment
`

func (l *Lexer) matchMeta(source *symbol.Source) (*symbol.Unknown, error) {
	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	err = buffer.SetLines(source.Origin.DefinitionLines())

	if err != nil {
		return nil, err
	}


	captures, err := l.context.Nvim.TsQueryOne(treesitter.Query{Language: "lua", Query: metaQuery})

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	unknown := newUnknownType()
	unknown.Documentation = source.Origin.DocumentationLines()

	return unknown, nil
}

func (l *Lexer) lexMeta(unknown *symbol.Unknown) (symbol.Symbol, error) {
	annotations, err := l.lexAnnotations(unknown.Documentation)

	if err != nil {
		return nil, err
	}

	if annotations.tipe != nil {
		return annotations.tipe, nil
	}

	l.context.Logger.Warn(fmt.Sprintf("Unknown meta type '%s' received", l.context.Current()))
	return unknown, nil
}
