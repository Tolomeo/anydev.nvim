package main

import (
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/lsp"
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

	luaCode = `
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local result = vim.lsp.buf_request_sync(0, 'textDocument/documentSymbol', { textDocument = textDocumentParams }, 2000)
		return vim.fn.json_encode(result[1])
	`

	result, err := nvimClient.ExecLua(luaCode, []any{})

	if err != nil {
		panic(err)
	}

	stringResult, ok := result.(string)

	if !ok {
		panic(fmt.Errorf("Error reading documentSymbol result"))
	}

	documentSymbols := lsp.TextDocumentDocumentSymbolResponse{}
	err = documentSymbols.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		panic(fmt.Errorf("Error marshalling documentSymbol response: %v", err))
	}

	fmt.Println(documentSymbols)

	if err := nvimClient.Close(); err != nil {
		fmt.Printf("Error closing nvim gracefully: %v", err)
	}
}
