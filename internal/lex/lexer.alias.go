package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (l *Lexer) lexAliasType(o *origin.AliasOrigin) (symbol.Type, error) {
	aliasType, hasType := o.Captures().Find("alias.type")

	if !hasType {
		return nil, fmt.Errorf("No type capture found for alias type '%s'", l.context.Target().Identifier())
	}

	// TODO: attach documentation
	lexedAlias, err := l.lexTypeAnnotation(TypeAnnotation{aliasType.Node.Text})

	if err != nil {
		return nil, err
	}

	return lexedAlias, nil
}

func (l *Lexer) lexAliasEnumeratorType(o *origin.AliasEnumeratorOrigin) (symbol.Type, error) {
	typeCaptures, hasTypes := o.Captures().FindAll("alias.type")

	if !hasTypes {
		return nil, fmt.Errorf("No type members found for alias enumerator '%s'", l.context.Target().Identifier())
	}

	enumeratorTypes, _ := slicesx.MapFunc(typeCaptures, func(capture treesitter.Capture) (string, error) {
		return capture.Node.Text, nil
	})

	enumeratorType, err := l.lexTypeAnnotation(TypeAnnotation{strings.Join(enumeratorTypes, "|")})

	if err != nil {
		return nil, err
	}

	return enumeratorType, nil
}
