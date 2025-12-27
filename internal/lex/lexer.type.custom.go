package lex

import (
	"fmt"
	"regexp"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var aliasQuery = func(aliasName string) treesitter.Query {
	return treesitter.Query{
		Language: "luadoc",
		Query: fmt.Sprintf(`
			(alias_annotation) @alias
			(#match? @alias "\\@alias %s")
		`, regexp.QuoteMeta(aliasName)),
	}
}

func (l *Lexer) lexCustomType(name string) error {
	if _, exists := l.context.Result().Types[name]; exists {
		l.context.Logger.Info("Skipping '%s': lexed type already found")
		return nil
	}

	l.context.Result().Types[name] = struct{}{}

	l.context.Fork(name, func(name string) error {
		fmt.Println(name)
		fmt.Println()
		fmt.Printf("%+v", aliasQuery(name).Query)
		fmt.Println()

		reference, err := l.crawler.SourceType(name)

		if err != nil {
			fmt.Printf("Reference error: %v", err)
		}

		fmt.Printf("\nReference source:\n%+v\n\n", reference.Origin)

		l.context.Result().Types[name] = newUnknownType()

		return nil
	})

	return nil
}
