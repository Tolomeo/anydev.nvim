package main

import (
	"bytes"
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"

	"text/template"
)

func main() {
	nvimClient, err := nvim.New()

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %v", err))
	}

	err = nvimClient.Open()

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %v", err))
	}

	luaTpl, err := template.New("lua").Parse(`
		local util = require('vim.lsp.util')

		local args = {...}
		local file = args[1]

		vim.cmd(string.format("e %s/anydev.lua", "{{.Dir}}"))

		-- Wait for LSP to attach
		vim.wait(2000, function()
			return next(vim.lsp.get_active_clients()) ~= nil
		end)

		-- Send a textDocument/documentSymbol request
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local result = vim.lsp.buf_request_sync(0, 'textDocument/documentSymbol', { textDocument = textDocumentParams }, 2000)

		-- Convert Lua table result to JSON for Go to decode
		return vim.fn.json_encode("{{.InitFile}}")
	`)

	if (err != nil) {
		panic(fmt.Errorf("Error parsing lua code template: %v", err))
	}

	var luaTplResult bytes.Buffer

	err = luaTpl.Execute(&luaTplResult, nvimClient.Options().Config())

	if (err != nil) {
		panic(fmt.Errorf("Error parsing lua code template: %v", err))
	}

	luaCode := luaTplResult.String()

	result, err := nvimClient.ExecLua(luaCode, []any{ "/Users/diegofrattini/Projects/anydev.nvim/.config/nvim/anydev.lua" })

	if err != nil {
		panic(fmt.Errorf("Error executing lua: %v", err))
	}

	fmt.Println(result)

	if err := nvimClient.Close(); err != nil {
		fmt.Printf("Error closing nvim gracefully: %v", err)
	}
}
