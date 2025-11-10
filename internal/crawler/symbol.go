package crawler

import "github.com/Tolomeo/anydev.nvim/internal/crawler/symbol"

type Symbol any

type source struct {
	// docString  []string
	definition []string
}

type Class struct {
	source
	symbol.ClassSymbol
}

func NewClass(path string, definition []string) *Class {
	class := Class{
		source{
			definition: definition,
		},
		symbol.ClassSymbol{
			Name: path,
		},
	}

	return &class
}
