package main

import (
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
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

	luacode := `
		local util = require('vim.lsp.util')

		-- Create a buffer with example Lua code
		vim.cmd('enew')
		vim.api.nvim_buf_set_lines(0, 0, -1, false, {
			"local function foo()",
			"  local x = 42",
			"  return x * 2",
			"end",
			"foo()"
		})

		-- Start the LSP client (lua-language-server must be in PATH)
		local client_id = vim.lsp.start({
			cmd = { "lua-language-server" },
			root_dir = vim.fn.getcwd(),
		})

		-- Wait for LSP to attach
		vim.wait(2000, function()
			return next(vim.lsp.get_active_clients()) ~= nil
		end)


		-- Send a textDocument/documentSymbol request
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local result = vim.lsp.buf_request_sync(0, 'textDocument/documentSymbol', { textDocument = textDocumentParams }, 2000)

		-- Convert Lua table result to JSON for Go to decode
		return vim.fn.json_encode(vim.g.test)
	`

	result, err := nvimClient.ExecLua(luacode, []any{})

	if err != nil {
		panic(fmt.Errorf("Error executing lua: %v", err))
	}

	fmt.Println(result)

	if err := nvimClient.Close(); err != nil {
		fmt.Printf("Error closing nvim gracefully: %v", err)
	}
}
