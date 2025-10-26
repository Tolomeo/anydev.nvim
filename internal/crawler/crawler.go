package crawler

import (
	"fmt"
	"slices"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
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

	fmt.Println(bufferName)

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

	fmt.Println(lines)

	documentSymbols, err := c.nvim.GetDocumentSymbols()

	if err != nil {
		return err
	}

	fmt.Printf("%+v", documentSymbols)

	ref := slices.IndexFunc(documentSymbols.Result, func(s lsp.DocumentSymbol) bool {
		return s.Name == "ref"
	})

	if ref == -1 {
		return fmt.Errorf("Error retrieving ref from document symbols: %+v", documentSymbols)
	}

	symbol := documentSymbols.Result[ref]

	fmt.Println("Symbol")
	fmt.Printf("%+v", symbol)

	completion, err := c.nvim.GetCompletion(uint(symbol.SelectionRange.End.Line), uint(symbol.SelectionRange.End.Character))

	if err != nil {
		return fmt.Errorf("Error retrieving completion information: %v", err)
	}

	fmt.Println("Completion")
	fmt.Printf("%+v", completion)

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
