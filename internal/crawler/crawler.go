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
	path string
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

	documentSymbols, err := c.nvim.GetDocumentSymbols()

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
	c.nvim.GetHover(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("Declaration")
	c.nvim.GetDeclaration(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("Definition")
	c.nvim.GetDefinition(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("TypeDefinition")
	c.nvim.GetTypeDefinition(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("Implementation")
	c.nvim.GetImplementation(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))

	return nil
}

func (c *crawler) Crawl(subpath string) (Symbol, error) {
	tempFile := c.nvim.Options().Config().Dir() + "anydev.lua"
	err := c.nvim.Edit(tempFile)

	if err != nil {
		return struct{}{}, err
	}

	return c.crawl(c.path, subpath)
}

func (c *crawler) crawl(path string, subpath string) (Symbol, error) {
	fullpath := path + "." + subpath

	fmt.Println(fullpath)

	err := c.Debug(path, subpath)

	if err != nil {
		return struct{}{}, err
	}

	/* statement := "local ref = " + fullpath

	err := c.nvim.SetBufferLines([]string{
		statement,
	})

	if err != nil {
		return struct{}{}, err
	}

	documentSymbols, err := c.nvim.GetDocumentSymbols()

	if err != nil {
		return struct{}{}, err
	}

	documentSymbol, ok := slicesx.FindFunc(documentSymbols.Result, func(s lsp.DocumentSymbol) bool {
		return s.Name == "ref"
	})

	if !ok {
		return struct{}{}, fmt.Errorf("Error retrieving ref from document symbols: %+v", documentSymbols)
	}

	err = c.nvim.GetHover(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))

	if err != nil {
		return struct{}{}, fmt.Errorf("Error reading hover: %w", err)
	} */

	typeName, err := c.getTypeName(path, subpath)

	if err != nil {
		return struct{}{}, err
	}

	if typeName != "table" {
		return struct{}{}, nil
	}

	children, err := c.getChildren(fullpath)

	if err != nil {
		return struct{}{}, err
	}

	if len(children) > 0 {
		classSymbol := ClassSymbol{
			Name: subpath,
		}

		for _, childpath := range children {
			symbol, err := c.crawl(fullpath, childpath)

			if err != nil {
				return struct{}{}, fmt.Errorf("Error crawling %s.%s: %w", fullpath, childpath, err)
			}

			classSymbol.Children = append(classSymbol.Children, symbol)
		}

		return classSymbol, nil
	}

	return struct{}{}, nil
}

func (c *crawler) getTypeName(path string, subpath string) (string, error) {
	var fullpath string

	switch subpath {
	case "and", "function", "or", "repeat", "false", "true":
		fullpath = path + "['" + subpath + "']"
	default:
		fullpath = path + "." + subpath
	}

	luaCode := fmt.Sprintf("return type(%s)", fullpath)

	result, err := c.nvim.ExecLua(luaCode, []any{})

	if err != nil {
		return "", fmt.Errorf("Error getting the type of %s: %w", path, err)
	}

	typeName, ok := result.(string)

	if !ok {
		return "", fmt.Errorf("Error getting the type of %s: Error converting the result to a string")
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
	err := c.nvim.WaitForLSP()

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
	instance := crawler{
		nvim: nvim,
		path: "_G",
	}

	return &instance
}
