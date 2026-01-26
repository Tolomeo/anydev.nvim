package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

type metaAtAnnotations struct {
	AtType *AtTypeAnnotation
}

func (l *Lexer) lexMetaAnnotations(docblock []string) (metaAtAnnotations, error) {
	annotations := metaAtAnnotations{
		AtType: nil,
	}

	buffer, err := l.context.Nvim().NewBuffer()

	if err != nil {
		return annotations, err
	}

	defer buffer.Close()

	err = buffer.SetLines(docblock)

	if err != nil {
		return annotations, err
	}

	atType, err := l.lexAtTypeAnnotations(buffer)

	if err != nil {
		return annotations, fmt.Errorf("Error lexing type annotation: %w", err)
	}

	annotations.AtType = atType
	return annotations, nil
}

func (l *Lexer) lexMetaValue(origin *symbol.MetaOrigin) (symbol.Type, error) {
	unknown := symbol.NewUnknown()
	unknown.Documentation = origin.Documentation()

	annotations, err := l.lexMetaAnnotations(origin.Documentation())

	if err != nil {
		return nil, err
	}

	if annotations.AtType == nil {
		l.context.Logger().Warn(fmt.Sprintf("Unknown meta type '%s' received", l.context.Target().Name()))
		return unknown, nil
	}

	if len(annotations.AtType.Types) < 1 {
		return nil, fmt.Errorf("Error lexing @type annotations for meta type '%s': no type annotations found", l.context.Target().Name())
	}

	lexedType, err := l.lexTypeAnnotation(annotations.AtType.Types[0])

	if err != nil {
		return nil, err
	}

	return lexedType, nil
}
