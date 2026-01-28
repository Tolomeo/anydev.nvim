package lex

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (l *Lexer) lexAliasType(aliasOrigin *origin.AliasOrigin) (symbol.Type, error) {
	// TODO: attach documentation
	lexedAlias, err := l.lexTypeAnnotation(TypeAnnotation{aliasOrigin.Type()})

	if err != nil {
		return nil, err
	}

	return lexedAlias, nil
}

func (l *Lexer) lexAliasEnumeratorType(aliasEnumeratorOrigin *origin.AliasEnumeratorOrigin) (symbol.Type, error) {
	enumeratorType, err := l.lexTypeAnnotation(TypeAnnotation{strings.Join(aliasEnumeratorOrigin.Types(), "|")})

	if err != nil {
		return nil, err
	}

	return enumeratorType, nil
}
