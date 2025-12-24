package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
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

func (l *Lexer) matchMeta(source *crawl.Source) (*lexed.Unknown, error) {
	buffer, err := l.context.Nvim.Buffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	err = buffer.SetLines(source.Origin().Definition())

	if err != nil {
		return nil, err
	}

	captures, err := l.context.Nvim.TsQueryOne(nvim.TsQueryConfig{Language: "lua", Query: metaQuery})

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	unknown := newUnknownType()
	unknown.Documentation = source.Origin().Documentation()

	return unknown, nil
}

func (l *Lexer) lexMeta(unknown *lexed.Unknown) (lexed.Symbol, error) {
	annotations, err := l.lexAnnotations(unknown.Documentation)

	if err != nil {
		return nil, err
	}

	if annotations.Type != nil {
		return annotations.Type, nil
	}

	l.context.Logger.Warn(fmt.Sprintf("Unknown meta type '%s' received", l.context.Current()))
	return unknown, nil
}
