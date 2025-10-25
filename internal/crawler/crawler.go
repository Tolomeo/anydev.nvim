package crawler

import (
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type crawler struct {
	nvim *nvim.Nvim
}

func (c *crawler) Crawl() error {
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
		"local vim_api = vim",
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

	fmt.Println(documentSymbols)

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
