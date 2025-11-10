package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
	"github.com/Tolomeo/anydev.nvim/internal/slicesx"
)

var ErrNotFound = errors.New("Not found")

type crawler struct {
	nvim          *nvim.Nvim
	scratchBuffer string
}

func (c *crawler) scratch(lines []string) error {
	_, err := c.nvim.Open(c.scratchBuffer)

	if err != nil {
		return err
	}

	err = c.nvim.SetBufferLines(lines)

	return nil
}

func (c *crawler) Crawl(path string) (Symbol, error) {
	foundSymbol, err := c.findSymbol(path)

	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", foundSymbol)

	switch foundSymbol := foundSymbol.(type) {
	case *Class:
		err := c.crawlClass(foundSymbol)

		if err != nil {
			return nil, err
		}

		return foundSymbol, nil

	}

	return struct{}{}, nil
}

func (c *crawler) findSymbol(path string) (Symbol, error) {
	fmt.Println("findSymbol", path)

	runtimeSymbol, err := c.findRuntimeSymbol(path)

	switch {
	case err == nil:
		return runtimeSymbol, nil
	case errors.Is(err, ErrNotFound):
	default:
		return nil, fmt.Errorf("Error retrieving runtime symbol for %s: %w", path, err)
	}

	lspSymbol, err := c.findLSPSymbol(path)

	switch {
	case err == nil:
		return lspSymbol, nil
	case errors.Is(err, ErrNotFound):
	default:
		return nil, fmt.Errorf("Error retrieving lsp symbol for %s: %w", path, err)
	}

	return nil, ErrNotFound
}

func (c *crawler) findRuntimeSymbol(path string) (Symbol, error) {
	runtimeType, err := c.getRuntimeTypeName(path)

	if err != nil {
		return nil, err
	}

	fmt.Println("findRuntimeSymbol", path, runtimeType)

	switch runtimeType {
	case "table":
		return NewClass(path), nil

	case "function":
		fmt.Println("function", path)
		return struct{}{}, nil
	}

	return struct{}{}, nil
}

func (c *crawler) findLSPSymbol(path string) (Symbol, error) {
	fmt.Println("findLSPSymbol", path)
	return struct{}{}, nil
}

func (c *crawler) crawlClass(symbol *Class) error {
	documentation, err := c.getSymbolDocumentation(symbol.Name)

	if err != nil {
		return err
	}

	symbol.Documentation = documentation

	err = c.getDefinition(symbol.Name)

	if err != nil {
		return err
	}

	children, err := c.getChildren(symbol.Name)

	if err != nil {
		return err
	}

	if len(children) == 0 {
		return nil
	}

	path := symbol.Name

	for _, fieldPath := range children {
		child, err := c.Crawl(path + "." + fieldPath)

		if err != nil {
			return fmt.Errorf("Error crawling %s.%s: %w", path, fieldPath, err)
		}

		symbol.Fields = append(symbol.Fields, child)
	}

	return nil
}

func (c *crawler) getDefinition(path string) error {
	definitionLocation, err := c.follow(path)

	switch {
	case err == nil:
	case errors.Is(err, ErrNotFound):
		fmt.Printf("Definition not found for %s\n", path)
		return nil
	default:
		return err
	}

	url, err := url.Parse(string(definitionLocation.TargetUri))

	fmt.Printf("%v:%v,%v\n", url.Path, definitionLocation.TargetRange.Start.Line, definitionLocation.TargetRange.Start.Character)

	if err != nil {
		return err
	}

	_, err = c.nvim.Open(url.Path)

	if err != nil {
		return err
	}

	cursorPosition := nvim.CursorPosition{
		Line:      uint(definitionLocation.TargetRange.Start.Line),
		Character: uint(definitionLocation.TargetRange.Start.Character),
	}

	_, err = c.nvim.GetAnnotatedAssignmentBufferLines(cursorPosition)

	if err != nil {
		return err
	}

	return nil
}

func (c *crawler) follow(path string) (lsp.DefinitionLocation, error) {
	lines := []string{"local ref = " + path}

	err := c.scratch(lines)

	if err != nil {
		return lsp.DefinitionLocation{}, err
	}

	line, character := uint(0), uint(len(lines[0]))

	definition, err := c.nvim.GetLSPDefinition(line, character)

	if err != nil {
		return lsp.DefinitionLocation{}, err
	}

	if len(definition.Result) == 0 {
		return lsp.DefinitionLocation{}, ErrNotFound
	}

	return definition.Result[0], nil
}

func (c *crawler) getSymbolDocumentation(path string) (string, error) {
	lines := []string{"local ref = " + path}

	err := c.scratch(lines)

	if err != nil {
		return "", err
	}

	line, character := uint(0), uint(len(lines[0]))

	hover, err := c.nvim.GetLSPHover(line, character)

	if err != nil {
		return "", fmt.Errorf("Error retrieving documentation for symbol %s: %w", path, err)
	}

	return hover.Result.Contents.Value, nil
}

func (c *crawler) getRuntimeTypeName(path string) (string, error) {
	runtimePath := path
	parts := strings.Split(runtimePath, ".")

	switch len(parts) {
	case 1:
	default:
		tail := parts[len(parts)-1]
		// https://www.lua.org/manual/5.1/manual.html#2.1
		switch tail {
		case "and", "break", "do", "else", "elseif", "end", "false", "for", "function", "if", "in", "local", "nil", "not", "or", "repeat", "return", "then", "true", "until", "while":
			head := parts[:len(parts)-1]
			runtimePath = strings.Join(head, ".") + "['" + tail + "']"
		}
	}

	luaCode := fmt.Sprintf("return type(%s)", runtimePath)

	result, err := c.nvim.ExecLua(luaCode, []any{})

	if err != nil {
		return "", fmt.Errorf("Error getting the type of %s: %w", path, err)
	}

	typeName, ok := result.(string)

	if !ok {
		return "", fmt.Errorf("Error getting the type of %s: Error converting the result to a string", runtimePath)
	}

	return typeName, nil
}

/* func (c *crawler) crawlClass(symbol *ClassSymbol, subpath string) error {
	children, err := c.getChildren(subpath)

	if err != nil {
		return err
	}

	for _, child := range children {
		if child.Label == "__index" {
			continue
		}

		if child.Kind == nil {
			return fmt.Errorf("Error reading nil completion item kind: %#v", child)
		}

		childSubpath := subpath + "." + child.Label

		if child.InsertText != nil {
			childSubpath = subpath + "." + *(child.InsertText)
		}

		fmt.Println(childSubpath, child.InsertText, *child.Kind)
		fmt.Printf("%+v\n\n", child)

		switch *(child.Kind) {
		case 5:
			childClassSymbol := ClassSymbol{
				Name:       child.Label,
				Deprecated: child.Deprecated,
				Detail:     child.Detail,
			}

			if child.InsertText != nil {
				childClassSymbol.Name = *(child.InsertText)
			}

			err = c.crawlClass(&childClassSymbol, childSubpath)

			if err != nil {
				return fmt.Errorf("Error crawling %s: %w", childSubpath, err)
			}

			symbol.Children = append(symbol.Children, childClassSymbol)
		}
	}

	return nil
} */

func (c *crawler) getChildren(path string) ([]string, error) {
	err := c.nvim.StartLSP()

	if err != nil {
		return []string{}, fmt.Errorf("Error getting path %s children: %w", path, err)
	}

	completionPath := "lua " + path + "."
	getcompletionResult, err := c.nvim.CallFunction("getcompletion", []any{completionPath, "cmdline"})

	if err != nil {
		return []string{}, fmt.Errorf("Error getting path %s children: %w", path, err)
	}

	result, ok := slicesx.AnyToString(getcompletionResult.([]any))

	if !ok {
		return []string{}, fmt.Errorf("Error getting path %s children", path)
	}

	return result, nil
}

func (c *crawler) Debug(path string, subpath string) error {
	statement := "local ref = " + path + "." + subpath

	err := c.nvim.SetBufferLines([]string{
		statement,
	})

	if err != nil {
		return err
	}

	lines, err := c.nvim.GetBufferLines()

	if err != nil {
		return err
	}

	fmt.Println(lines)

	documentSymbols, err := c.nvim.GetLSPDocumentSymbols()

	if err != nil {
		return err
	}

	documentSymbol, ok := slicesx.FindFunc(documentSymbols.Result, func(s lsp.DocumentSymbol) bool {
		return s.Name == "ref"
	})

	if !ok {
		return fmt.Errorf("Error retrieving ref from document symbols: %+v", documentSymbols)
	}

	fmt.Println("Hover")
	c.nvim.GetLSPHover(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("Declaration")
	c.nvim.GetLSPDeclaration(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("Definition")
	c.nvim.GetLSPDefinition(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("TypeDefinition")
	c.nvim.GetLSPTypeDefinition(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("Implementation")
	c.nvim.GetLSPImplementation(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))

	return nil
}

func New(nvim *nvim.Nvim) *crawler {
	scratchBuffer := nvim.Options().Config().Dir() + "anydev.lua"

	instance := crawler{
		nvim:          nvim,
		scratchBuffer: scratchBuffer,
	}

	return &instance
}
