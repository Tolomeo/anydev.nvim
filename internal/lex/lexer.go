package lex

import (
	"fmt"
	"path"
)

type Lexer struct {
	context *lexingContext
}

func (l *Lexer) Lex(paths []string, context *lexingContext) error {
	l.context = context
	defer func() { l.context = nil }()

	for _, path := range paths {
		err := l.context.provide(path, l.lex)

		if err != nil {
			return fmt.Errorf("Error lexing %s: %w", path, err)
		}
	}

	return nil
}

func (l *Lexer) lex(path string) error {
	if _, exists := l.context.result.Runtime[path]; exists {
		return fmt.Errorf("Error lexing '%s': path already found", path)
	}

	l.context.Result().Runtime[path] = struct{}{}

	_, err := l.context.crawler.SourceRuntime(path)

	if err != nil {
		return fmt.Errorf("Error lexing %s: %w", path, err)
	}

	// fmt.Printf("%+v\n", source)
	/* switch v := source.(type) {
	case *crawl.TableSource:
		fmt.Println(v, "table")
		l.context.result.Runtime[path] = struct{}{}
	case *crawl.FunctionSource:
		lexedFunction, err := l.lexFunction(v)

		if err != nil {
			return err
		}

		l.context.Result().Runtime[path] = lexedFunction
		return nil
	case *crawl.VariableSource:
		fmt.Println("variable")
		l.context.result.Runtime[path] = struct{}{}
	default:
		return fmt.Errorf("Error lexing '%s' source: unknown source type", path)
	} */

	return nil
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
