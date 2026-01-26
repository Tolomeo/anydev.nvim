package lex

import (
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

func (l *Lexer) lexFieldType(origin *symbol.FieldOrigin) (symbol.Type, error) {
	// TODO: assign qualifiers taken from annotations
	lexedField, err := l.lexTypeAnnotation(TypeAnnotation{origin.GetType()})

	if err != nil {
		return nil, err
	}

	return lexedField, nil
}
