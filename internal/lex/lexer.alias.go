package lex

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	type_ "github.com/Tolomeo/anydev.nvim/internal/domain/type"
)

func (l *Lexer) lexAliasType(aliasOrigin *origin.AliasOrigin) (type_.Type, error) {
	// TODO: attach documentation
	lexedAlias, err := l.lexTypeAnnotation(TypeAnnotation{aliasOrigin.Type()})

	if err != nil {
		return nil, err
	}

	return lexedAlias, nil
}

func (l *Lexer) lexAliasEnumeratorType(aliasEnumeratorOrigin *origin.AliasEnumeratorOrigin) (type_.Type, error) {
	enumeratorType, err := l.lexTypeAnnotation(TypeAnnotation{strings.Join(aliasEnumeratorOrigin.Types(), "|")})

	if err != nil {
		return nil, err
	}

	return enumeratorType, nil
}
