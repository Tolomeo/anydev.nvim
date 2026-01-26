package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

func (l *Lexer) lexMetaValue(origin *symbol.MetaOrigin) (symbol.Type, error) {
	unknown := symbol.NewUnknown()
	unknown.Documentation = origin.Documentation()

	annotations, err := l.lexAtAnnotations(origin.Documentation())

	switch {
	case err != nil:
		return nil, err
	case annotations.AtType == nil:
		l.context.Logger().Warn(fmt.Sprintf("Unknown meta type '%s' received", l.context.Target().Name()))
		return unknown, nil
	case len(annotations.AtType.Types) < 1:
		return nil, fmt.Errorf("Error lexing @type annotations for meta type '%s': no type annotations found", l.context.Target().Name())
	}

	lexedType, err := l.lexTypeAnnotation(annotations.AtType.Types[0])

	if err != nil {
		return nil, err
	}

	return lexedType, nil
}
