package main

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex"
)

const path string = "vim.deepcopy"

func main() {
	extractor, err := lex.NewLexer()

	if err != nil {
		panic(fmt.Errorf("Error extracting %s: %w", path, err))
	}

	err = extractor.Lex(path)

	if err != nil {
		panic(fmt.Errorf("Error extracting %s: %w", path, err))
	}
}
