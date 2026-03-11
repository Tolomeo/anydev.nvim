package main

import (
	"github.com/Tolomeo/anydev.nvim/internal/extract"
)

const debug = true

var values = []string{
	"vim",
}

var override = extract.Override{
	Definition: map[string][]string{
		"vim":        {"vim = {}"},
		"vim.base64": {"vim.base64 = {}"},
		"vim.cmd":    {"---@type fun(command: string|table)|table<string,fun(...:any)>", "vim.cmd = ..."},
		"vim.env":    {"---@type table<string, string>", "vim.env = ..."},
		// TODO: improve module export query
		"vim.iter": {"---@type IterMod", "vim.iter = ..."},
	},
}

func main() {
	extractor, err := extract.NewExtractor(extract.Options{Debug: debug, Override: override})

	if err != nil {
		panic(err)
	}

	for _, value := range values {
		err := extractor.Extract("value", value)

		if err != nil {
			panic(err)
		}
	}
}
