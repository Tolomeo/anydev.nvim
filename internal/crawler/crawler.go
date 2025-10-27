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
}

func (c *crawler) Crawl(api string) error {
	err := c.nvim.Open()

	if err != nil {
		return fmt.Errorf("Error opening nvim: %v", err)
	}

	tempFile := c.nvim.Options().Config().Dir() + "anydev.lua"
	err = c.nvim.Edit(tempFile)

	if err != nil {
		return err
	}

	bufferName, err := c.nvim.GetBufferName()

	if err != nil {
		return err
	}

	fmt.Println("BufferName")
	fmt.Println(bufferName)
	fmt.Println("/BufferName")

	err = c.nvim.SetBufferLines([]string{
		"local ref = _G.",
	})

	if err != nil {
		return err
	}

	lines, err := c.nvim.GetBufferLines()

	if err != nil {
		return err
	}

	fmt.Println("Lines")
	fmt.Println(lines)
	fmt.Println("/Lines")

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

	fmt.Println("Symbol")
	fmt.Printf("%#v\n", documentSymbol)
	fmt.Println(documentSymbol.SelectionRange.End.Line, documentSymbol.SelectionRange.End.Character)
	fmt.Println("/Symbol")

	completion, err := c.nvim.GetCompletion(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))

	if err != nil {
		return fmt.Errorf("Error retrieving completion information: %v", err)
	}

	fmt.Println("Completion")
	fmt.Printf("%#v\n", completion)
	fmt.Println("/Completion")

	symbol, ok := slicesx.FindFunc(completion.Result.Items, func(item lsp.CompletionItem) bool {
		return item.Label == api
	})

	if !ok {
		return fmt.Errorf("Error retrieving completion item from completion: %+v", completion)
	}

	fmt.Println("Symbol")
	fmt.Printf("%#v\n", symbol)
	fmt.Println("/Symbol")

	if err := c.nvim.Close(); err != nil {
		fmt.Printf("Error closing nvim gracefully: %v", err)
	}

	return nil
}

func New(nvim *nvim.Nvim) *crawler {
	instance := crawler{
		nvim: nvim,
	}

	return &instance
}
