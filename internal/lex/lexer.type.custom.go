package lex

import (
	"fmt"
)

func (l *Lexer) lexCustomType(name string) error {
	if _, exists := l.context.Result().Types[name]; exists {
		l.context.Logger.Info("Skipping '%s': lexed type already found")
		return nil
	}

	l.context.Result().Types[name] = struct{}{}

	l.context.Fork(name, func(name string) error {
		typeSource, err := l.crawler.SourceType(name)

		if err != nil {
			fmt.Printf("Reference error: %v", err)
		}

		fmt.Printf("\nReference source:\n%+v\n\n", typeSource.Origin)

		switch typeSource.Origin.Type() {
		case "alias_annotation":
		case "class_annotation":
			//TODO enum
		}

		fmt.Printf("\nType: %+v\n\n", typeSource.Origin.Definition)

		l.context.Result().Types[name] = newUnknownType()

		return nil
	})

	return nil
}
