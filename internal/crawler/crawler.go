package crawler

import (
	// "encoding/json"
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
	"github.com/Tolomeo/anydev.nvim/internal/slicesx"
)

type crawler struct {
	nvim *nvim.Nvim
	// path string
	scratchBuffer string
}

func (c *crawler) openScratchBuffer() error {
	_, err := c.nvim.Open(c.scratchBuffer)

	if err != nil {
		return err
	}

	return nil
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

func (c *crawler) Crawl(subpath string) (Symbol, error) {
	_, err := c.nvim.Open(c.scratchBuffer)

	if err != nil {
		return struct{}{}, err
	}

	return c.crawl("_G", subpath)
}

func (c *crawler) crawl(path string, subpath string) (Symbol, error) {
	fullpath := path + "." + subpath

	fmt.Println(fullpath)

	typeName, err := c.getRuntimeTypeName(path, subpath)

	if err != nil {
		return struct{}{}, err
	}

	switch typeName {
	case "table":
		name := subpath
		crawledTable, err := c.crawlTable(fullpath, name)

		if err != nil {
			return struct{}{}, err
		}

		return crawledTable, nil
	}

	// TODO: return Unrecognised type error
	return struct{}{}, nil
}

func (c *crawler) crawlTable(path string, name string) (ClassSymbol, error) {
	classSymbol := ClassSymbol{
		Name: name,
	}

	documentation, err := c.getSymbolDocumentation(path)

	classSymbol.Documentation = documentation

	if err != nil {
		return classSymbol, err
	}

	children, err := c.getChildren(path)

	if err != nil {
		return classSymbol, err
	}

	if len(children) == 0 {
		return classSymbol, nil
	}

	for _, childpath := range children {
		child, err := c.crawl(path, childpath)

		if err != nil {
			return classSymbol, fmt.Errorf("Error crawling %s.%s: %w", path, childpath, err)
		}

		classSymbol.Children = append(classSymbol.Children, child)
	}

	return classSymbol, nil
}

func (c *crawler) getSymbolDocumentation(path string) (string, error) {
	err := c.openScratchBuffer()

	if err != nil {
		return "", err
	}

	assignment := "local ref = " + path

	err = c.nvim.SetBufferLines([]string{
		assignment,
	})

	if err != nil {
		return "", err
	}

	line, character := uint(0), uint(len(assignment))

	hover, err := c.nvim.GetLSPHover(line, character)

	if err != nil {
		return "", fmt.Errorf("Error retrieving documentation for symbol %s: %w", path, err)
	}

	return hover.Result.Contents.Value, nil
}

func (c *crawler) getRuntimeTypeName(path string, subpath string) (string, error) {
	fullpath := path + "." + subpath

	switch subpath {
	case "and", "function", "or", "repeat", "false", "true":
		fullpath = path + "['" + subpath + "']"
	default:
	}

	luaCode := fmt.Sprintf("return type(%s)", fullpath)

	result, err := c.nvim.ExecLua(luaCode, []any{})

	if err != nil {
		return "", fmt.Errorf("Error getting the type of %s: %w", path, err)
	}

	typeName, ok := result.(string)

	if !ok {
		return "", fmt.Errorf("Error getting the type of %s: Error converting the result to a string", fullpath)
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
		return []string{}, fmt.Errorf("Error getting path %s children: %w", err)
	}

	completionPath := "lua " + path + "."
	getcompletionResult, err := c.nvim.CallFunction("getcompletion", []any{completionPath, "cmdline"})

	if err != nil {
		return []string{}, fmt.Errorf("Error getting path %s children: %w", path, err)
	}

	result, ok := slicesx.AnyToString(getcompletionResult.([]any))

	if !ok {
		return []string{}, fmt.Errorf("Error getting path %s children: %w", path)
	}

	return result, nil
}

func New(nvim *nvim.Nvim) *crawler {
	scratchBuffer := nvim.Options().Config().Dir() + "anydev.lua"

	instance := crawler{
		nvim:          nvim,
		scratchBuffer: scratchBuffer,
	}

	return &instance
}
