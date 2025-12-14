package lex

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type Lexer struct {
	context *lexingContext
}

func (l *Lexer) Lex(paths []string, context *lexingContext) error {
	l.context = context
	defer func() { l.context = nil }()

	for _, path := range paths {
		if _, exists := l.context.result.Runtime[path]; exists {
			return fmt.Errorf("Error lexing '%s': path already found", path)
		}

		l.context.Result().Runtime[path] = struct{}{}

		symbol, err := l.context.provide(path, l.lex)

		if err != nil {
			return fmt.Errorf("Error lexing %s: %w", path, err)
		}

		l.context.result.Runtime[path] = symbol
	}

	return nil
}

func (l *Lexer) lex(path string) (lexed.Symbol, error) {
	source, err := l.context.crawler.SourceRuntime(path)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing %s: %w", path, err)
	case source == nil:
		return nil, fmt.Errorf("Error lexing %s: No source found", path)
	}

	sourceOrigin := source.Origin()

	if sourceOrigin == nil {
		symbol := newUnknownType()
		l.context.logger.Warn(fmt.Sprintf("Using '%v' for symbol '%s' without origin", symbol, path))
		l.context.result.Runtime[path] = symbol
		return symbol, nil
	}

	sourceType := sourceOrigin.Type()
	sourceDefinition := sourceOrigin.Definition()

	l.scratch([]string{sourceDefinition})

	switch sourceType {
	case "variable_declaration":
		symbol, err := l.lexVariableDeclaration(source)

		switch {
		case err != nil:
			return nil, err
		case symbol == nil:
			return nil, fmt.Errorf("Error lexing variable declaration: no symbol found in '%s'", sourceDefinition)
		}

		return symbol, nil

	case "function_declaration":
		l.lexFunctionDeclaration()
		return struct{}{}, nil
	default:
		return nil, fmt.Errorf("Error lexing %s: Unkonw origin type '%s'", path, sourceType)
	}
}

func (l *Lexer) lexVariableDeclaration(source *crawl.Source) (lexed.Symbol, error) {
	fmt.Println(l.context.path)
	fmt.Println("variable_declaration")

	table, err := l.lexTableDeclaration(source)

	switch {
	case err != nil:
		return nil, err
	case table == nil:
		return nil, nil
	}

	return table, nil
}

var tableDeclarationQuery string = `
	(variable_declaration
		(assignment_statement
			(variable_list
				name: (identifier))
				(expression_list
					value: (table_constructor) @table
				)
		)
	)
`

func (l *Lexer) lexTableDeclaration(source *crawl.Source) (lexed.Symbol, error) {
	captures, err := l.context.nvim.TsQuery(nvim.TsQueryConfig{Language: "lua", Query: tableDeclarationQuery})

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	return l.lexTable(source)
}

func (l *Lexer) lexFunctionDeclaration() {
	fmt.Println(l.context.path)
	fmt.Println("function_declaration")
}

func (l *Lexer) scratch(lines []string) error {
	buffer := path.Join(l.context.nvim.Options().Config().Dir(), "anydev.lexer.lua")

	_, err := l.context.nvim.Open(buffer)

	if err != nil {
		return err
	}

	err = l.context.nvim.SetBufferLines(lines)

	if err != nil {
		return err
	}

	return nil
}

func NewLexer() *Lexer {
	return &Lexer{}
}
