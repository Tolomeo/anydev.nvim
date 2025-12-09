package lex

import "github.com/Tolomeo/anydev.nvim/internal/lex/lexed"

type result struct {
	Runtime map[string]lexed.Symbol `json:"runtime" yaml:"runtime"`
	Types   map[string]lexed.Symbol `json:"types" yaml:"types"`
}

func newResult() *result {
	return &result{
		Runtime: map[string]lexed.Symbol{},
		Types:   map[string]lexed.Symbol{},
	}
}
