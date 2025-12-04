package main

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex"
)

const path string = "vim.validate"

func main() {
	extractor, err := lex.NewExtractor()

	if err != nil {
		panic(fmt.Errorf("Error extracting %s: %w", path, err))
	}

	err = extractor.Extract(path)

	if err != nil {
		panic(fmt.Errorf("Error extracting %s: %w", path, err))
	}
}
