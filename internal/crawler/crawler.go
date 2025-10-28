package crawler

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
	"github.com/Tolomeo/anydev.nvim/internal/slicesx"
	// "github.com/Tolomeo/anydev.nvim/internal/dumper"
)

type crawler struct {
	nvim *nvim.Nvim
	path string
}

func (c *crawler) Crawl(subpath string) error {
	err := c.nvim.Open()

	if err != nil {
		return fmt.Errorf("Error opening nvim: %v", err)
	}

	tempFile := c.nvim.Options().Config().Dir() + "anydev.lua"
	err = c.nvim.Edit(tempFile)

	if err != nil {
		return err
	}

	children, err := c.getChildren(c.path)

	if err != nil {
		return err
	}

	child, ok := slicesx.FindFunc(children, func(item lsp.CompletionItem) bool {
		return item.Label == subpath
	})

	if !ok {
		return fmt.Errorf("Error retrieving completion item from completion: %+w", children)
	}

	fmt.Println("Symbol")
	fmt.Printf("%#v\n", child)
	fmt.Println("/Symbol")

	completionItemKind := child.Kind

	if completionItemKind == nil {
		return fmt.Errorf("Error reading nil completion item kind: %#v", child)
	}

	switch *completionItemKind {
	case 5:
		classSymbol := ClassSymbol{
			Name:       child.Label,
			Deprecated: child.Deprecated,
			Detail:     child.Detail,
		}

		c.crawlClass(&classSymbol, subpath)

		fmt.Println("Class")
		fmt.Printf("%+v\n\n", classSymbol)
		fmt.Println("/Class")
	default:
		fmt.Println("not a class")
	}

	if err := c.nvim.Close(); err != nil {
		fmt.Printf("Error closing nvim gracefully: %v", err)
	}

	return nil
}

func (c *crawler) crawlClass(symbol *ClassSymbol, subpath string) {

}

func (c *crawler) getChildren(path string) ([]lsp.CompletionItem, error) {
	err := c.nvim.SetBufferLines([]string{
		"local ref = " + path + ".",
	})

	if err != nil {
		return []lsp.CompletionItem{}, err
	}

	documentSymbols, err := c.nvim.GetDocumentSymbols()

	if err != nil {
		return []lsp.CompletionItem{}, err
	}

	documentSymbol, ok := slicesx.FindFunc(documentSymbols.Result, func(s lsp.DocumentSymbol) bool {
		return s.Name == "ref"
	})

	if !ok {
		return []lsp.CompletionItem{}, fmt.Errorf("Error retrieving ref from document symbols: %+v", documentSymbols)
	}

	/* fmt.Println("Symbol")
	fmt.Printf("%#v\n", documentSymbol)
	fmt.Println("/Symbol") */

	completion, err := c.nvim.GetCompletion(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))

	if err != nil {
		return []lsp.CompletionItem{}, fmt.Errorf("Error retrieving completion information: %w", err)
	}

	/* fmt.Println("Completion")
	fmt.Printf("%#v\n", completion)
	fmt.Println("/Completion") */

	return completion.Result.Items, nil
}

func New(nvim *nvim.Nvim) *crawler {
	instance := crawler{
		nvim: nvim,
		path: "_G",
	}

	return &instance
}
