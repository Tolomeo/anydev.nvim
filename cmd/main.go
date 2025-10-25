package main

import (
	// "bytes"
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	// "text/template"
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

	tempFile := nvimClient.Options().Config().Dir() + "anydev.lua"
	err = nvimClient.Edit(tempFile)

	if err != nil {
		panic(err)
	}

	bufferName, err := nvimClient.GetBufferName()

	if err != nil {
		panic(err)
	}

	fmt.Println(bufferName)

	luaCode := `
		vim.wait(2000, function()
			return next(vim.lsp.get_active_clients()) ~= nil
		end)
	`

	_, err = nvimClient.ExecLua(luaCode, []any{})

	if err != nil {
		panic(err)
	}

	err = nvimClient.SetBufferLines([]string{
		"local vim_api = vim",
	})

	if err != nil {
		panic(err)
	}

	lines, err := nvimClient.GetBufferLines()

	if err != nil {
		panic(err)
	}

	fmt.Println(lines)
	/* luaTpl, err := template.New("lua").Parse(`
		local util = require('vim.lsp.util')

		local args = {...}

		vim.cmd(string.format("e %s/anydev.lua", "{{.Dir}}"))

		-- Wait for LSP to attach
		vim.wait(2000, function()
			return next(vim.lsp.get_active_clients()) ~= nil
		end)

		vim.api.nvim_buf_set_lines(0, 0, -1, false, {
			"local vim_api = vim",
		})

		-- Send a textDocument/documentSymbol request
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local result = vim.lsp.buf_request_sync(0, 'textDocument/documentSymbol', { textDocument = textDocumentParams }, 2000)

		-- Convert Lua table result to JSON for Go to decode
		return vim.fn.json_encode(1)
	`)

	if (err != nil) {
		panic(fmt.Errorf("Error parsing lua code template: %v", err))
	}

	var luaTplResult bytes.Buffer

	err = luaTpl.Execute(&luaTplResult, nvimClient.Options().Config())

	if (err != nil) {
		panic(fmt.Errorf("Error parsing lua code template: %v", err))
	}

	luaCode := luaTplResult.String() */

	luaCode = `return vim.fn.json_encode({ "test" })`

	result, err := nvimClient.ExecLua(luaCode, []any{})

	if err != nil {
		panic(fmt.Errorf("Error executing lua: %v", err))
	}

	fmt.Println(result)

	err = nvimClient.DeleteBuffer()

	if err != nil {
		panic(err)
	}

	if err := nvimClient.Close(); err != nil {
		fmt.Printf("Error closing nvim gracefully: %v", err)
	}
}
