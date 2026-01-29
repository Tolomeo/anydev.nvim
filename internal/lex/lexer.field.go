package lex

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	type_ "github.com/Tolomeo/anydev.nvim/internal/domain/type"
)

func (l *Lexer) lexFieldType(fieldOrigin *origin.FieldOrigin) (type_.Type, error) {
	// TODO: assign qualifiers taken from annotations
	lexedField, err := l.lexTypeAnnotation(TypeAnnotation{fieldOrigin.Type()})

	if err != nil {
		return nil, err
	}

	return lexedField, nil
}
