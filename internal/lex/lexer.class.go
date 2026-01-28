package lex

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (l *Lexer) lexClassType(classOrigin *origin.ClassOrigin) (*symbol.Table, error) {
	class := symbol.NewTable()
	class.Name = classOrigin.Name()
	classFields, err := l.context.Nvim().GetTypeCompletion(class.Name)

	if err != nil {
		return nil, err
	}

	for _, fieldName := range classFields {
		err := l.context.ExtractChild(class, fieldName)

		if err != nil {
			return nil, err
		}
	}

	return class, nil
}
