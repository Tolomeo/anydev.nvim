package main

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/extract"
)

const path string = "vim.validate"

func main() {
	extractor := extract.NewExtractor()

	err := extractor.Extract(path)

	if err != nil {
		panic(fmt.Errorf("Error extracting %s: %w", path, err))
	}
}
