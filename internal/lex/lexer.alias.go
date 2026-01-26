package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (l *Lexer) lexAliasType(origin *symbol.AliasOrigin) (symbol.Type, error) {
	aliasType, hasType := origin.Captures().Find("alias.type")

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

func (l *Lexer) lexAliasEnumeratorType(origin *symbol.AliasEnumeratorOrigin) (symbol.Type, error) {
	l.context.Logger().Debugf("Alias enumerator: %+v", origin)

	typeCaptures, hasTypes := origin.Captures().FindAll("alias.type")

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

	l.context.Logger().Debugf("%+v", origin)

	return enumeratorType, nil

}
