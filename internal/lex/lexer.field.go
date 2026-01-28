package lex

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (l *Lexer) lexFieldType(fieldOrigin *origin.FieldOrigin) (symbol.Type, error) {
	// TODO: assign qualifiers taken from annotations
	lexedField, err := l.lexTypeAnnotation(TypeAnnotation{fieldOrigin.Type()})

	if err != nil {
		return nil, err
	}

	return lexedField, nil
}
