package nvim

import (
	"fmt"
	"net/url"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
	"github.com/Tolomeo/anydev.nvim/internal/utils/anyx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (n *Nvim) startLSP() error {
	luaCode := `
		if vim.g.lua_ls_ready == true then return end

		local args = {...}
		local delay = args[1]

		vim.g.lua_ls_ready = false

		local lsp_inflight_events = {}

		vim.api.nvim_create_augroup("LuaLSReady", { clear = true })

		vim.api.nvim_create_autocmd("LspProgress", {
			group = "LuaLSReady",
			callback = function(args)
				local value = args.data.params.value
				local token = args.data.params.token

				if value.kind == "begin" then
					lsp_inflight_events[token] = value
				elseif value.kind == "end" then
					lsp_inflight_events[token] = nil
				end

				vim.g.lua_ls_ready = next(lsp_inflight_events) == nil
			end,
		})

		vim.lsp.enable("lua_ls")

		vim.wait(delay, function()
			return vim.g.lua_ls_ready == true
		end)
	`

	_, err := n.ExecLua(luaCode, []any{30000})

	if err != nil {
		return fmt.Errorf("Error starting lua lsp: %v", err)
	}

	return nil
}

func (n *Nvim) GetDocumentSymbols() (lsp.TextDocumentDocumentSymbolResponse, error) {
	documentSymbols := lsp.TextDocumentDocumentSymbolResponse{}

	err := n.startLSP()

	if err != nil {
		return documentSymbols, err
	}

	luaCode := `
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local result = vim.lsp.buf_request_sync(0, 'textDocument/documentSymbol', { textDocument = textDocumentParams }, 2000)
		return vim.fn.json_encode(result[1])
	`

	result, err := n.ExecLua(luaCode, []any{})

	if err != nil {
		return documentSymbols, fmt.Errorf("Error getting document symbols: %v", err)
	}

	stringResult, ok := result.(string)

	if !ok {
		return documentSymbols, fmt.Errorf("Error reading document symbols response: %v", result)
	}

	err = documentSymbols.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return documentSymbols, fmt.Errorf("Error unmarshalling document symbols response: %v", err)
	}

	return documentSymbols, nil
}

func (n *Nvim) GetHover(line uint, character uint) (lsp.TextDocumentHoverResponse, error) {
	hover := lsp.TextDocumentHoverResponse{}

	err := n.startLSP()

	if err != nil {
		return hover, err
	}

	luaCode := `
		local args = { ... }
		local line = args[1]
		local character = args[2]
		local position = { line = line, character = character }
		local textDocument = vim.lsp.util.make_text_document_params(0)
		local timeout = args[3]

		local result = vim.lsp.buf_request_sync(
			0,
			"textDocument/hover",
			{ textDocument = textDocument, position = position },
			timeout
		)

		return vim.fn.json_encode(result[1])
	`

	result, err := n.ExecLua(luaCode, []any{line, character, 15000})

	if err != nil {
		return hover, fmt.Errorf("Error getting lsp hover response: %v", err)
	}

	stringResult, ok := result.(string)

	if !ok {
		return hover, fmt.Errorf("Error reading lsp hover response: %v", result)
	}

	// fmt.Println(stringResult)

	err = hover.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return hover, fmt.Errorf("Error unmarshalling lsp hover response: %v", err)
	}

	return hover, nil
}

func (n *Nvim) GetLSPDefinition(line uint, character uint) ([]lsp.DefinitionLocation, error) {
	err := n.startLSP()

	if err != nil {
		return []lsp.DefinitionLocation{}, err
	}

	luaCode := `
		local args = {...}
		local line = args[1]
		local character = args[2]
		local delay = args[3]
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local positionParams = {line = line, character = character}
		local lspResponse, err = vim.lsp.buf_request_sync(0, 'textDocument/definition', { textDocument = textDocumentParams, position = positionParams }, delay)

		if err ~= nil then
			error(err)
		end

		local result = next(lspResponse[1]) and lspResponse[1] or { result = {} }

		return vim.fn.json_encode(result)
	`

	result, err := n.ExecLua(luaCode, []any{line, character, 15000})

	if err != nil {
		return []lsp.DefinitionLocation{}, fmt.Errorf("Error getting lsp definition: %v", err)
	}

	stringResult, ok := result.(string)

	if !ok {
		return []lsp.DefinitionLocation{}, fmt.Errorf("Error reading lsp definition response: %v", result)
	}

	response := lsp.TextDocumentDefinitionResponse{}
	err = response.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return []lsp.DefinitionLocation{}, fmt.Errorf("Error unmarshalling lsp definition response: %w", err)
	}

	return response.Result, nil
}

func (n *Nvim) GetDefinitionLocation(line uint, character uint) (*[]Location, error) {
	lspDefinitions, err := n.GetLSPDefinition(line, character)

	switch {
	case err != nil:
		return nil, err
	case len(lspDefinitions) == 0:
		return nil, nil
	}

	locations, err := slicesx.MapFunc(lspDefinitions, func(lspLocation lsp.DefinitionLocation) (Location, error) {
		location := Location{
			DefinitionLocation: lspLocation,
		}

		url, err := url.Parse(string(location.DefinitionLocation.TargetUri))

		if err != nil {
			return Location{}, err
		}

		location.Url = url.Path

		return location, nil
	})

	if err != nil {
		return nil, err
	}

	return &locations, nil
}

func (n *Nvim) GetCompletion(head string) ([]string, error) {
	err := n.startLSP()

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", head, err)
	}

	cmd := "lua " + head + "."
	getcompletionResult, err := n.CallFunction("getcompletion", []any{cmd, "cmdline"})

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", head, err)
	}

	result, err := anyx.ToSliceOf[string](getcompletionResult)

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", head, err)
	}

	return result, nil
}
