package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var aliasQuery = func(aliasName string) treesitter.Query {
	var query *treesitter.QueryBuilder = treesitter.NewQueryBuilder("luadoc")
	query.NewSExpression("alias_annotation").Var("alias", query).End()
	query.NewMatcher("match").Identifier("alias").ValueAsString(fmt.Sprintf(`\\@alias %s`, aliasName))

	return query.Marshal()
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
		fmt.Printf("%+v", string(aliasQuery(name).Language))
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
