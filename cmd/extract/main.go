package main

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/extract"
)

const path string = "vim.fs"

func main() {
	extractor, err := extract.NewExtractor()

	if err != nil {
		panic(fmt.Errorf("Error creating extractor: %w", err))
	}

	extractor.Extract(path)
}
