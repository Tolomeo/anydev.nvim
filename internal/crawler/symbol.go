package crawler

import "github.com/Tolomeo/anydev.nvim/internal/crawler/symbol"

type Symbol any

type Class struct {
	symbol.ClassSymbol
}

func NewClass(path string) *Class {
	class := Class{
		symbol.ClassSymbol {
			Name: path,
		},
	}

	return &class
}
