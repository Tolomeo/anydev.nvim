package nvim

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
	"github.com/Tolomeo/anydev.nvim/internal/utils/anyx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type CursorPosition struct {
	Line      uint
	Character uint
}

type Location struct {
	lsp.DefinitionLocation
	Url string
}

type Nvim struct {
	options options
	cmd     *exec.Cmd
	rpc     *rpc
}

func (n *Nvim) Options() options {
	return n.options
}

func (n *Nvim) Start() error {
	return n.cmd.Start()
}

func (n *Nvim) Quit() error {
	request := requestMessage{
		method: "nvim_command",
		params: []any{"qa!"},
	}
	_, err := n.rpc.Send(request)

	// EOF error expected
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	err = n.cmd.Wait()

	if err != nil {
		return err
	}

	return nil
}

func (n *Nvim) Open(file string) (string, error) {
	request := requestMessage{
		method: "nvim_command",
		params: []any{"edit" + file},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return "", fmt.Errorf("Error opening %s: %v\n", file, err)
	}

	return file, nil
}

func (n *Nvim) Write() error {
	request := requestMessage{
		method: "nvim_command",
		params: []any{"write"},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to write buffer: %v\n", err)
	}

	return nil
}

func (n *Nvim) GetBufferName() (string, error) {
	request := requestMessage{
		method: "nvim_buf_get_name",
		params: []any{0},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return "", fmt.Errorf("Error reading buffer name: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return "", fmt.Errorf("Error executing lua: %v\n", err)
	}

	return result.(string), nil
}

func (n *Nvim) SetBufferLines(lines []string) error {
	request := requestMessage{
		method: "nvim_buf_set_lines",
		params: []any{0, 0, -1, true, lines},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return err
	}

	_, err = response.Result()

	if err != nil {
		return err
	}

	return nil
}

func (n *Nvim) GetBufferText(startRow int, startCol int, endRow int, endCol int) ([]string, error) {
	request := requestMessage{
		method: "nvim_buf_get_text",
		params: []any{0, startRow, startCol, endRow, endCol, struct{}{}},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer text: %w\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer text: %w\n", err)
	}

	bufferText, err := anyx.ToSliceOf[string](result)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer text return value: %w", err)
	}

	return bufferText, nil
}

func (n *Nvim) GetBufferLines(start int, end int) ([]string, error) {
	request := requestMessage{
		method: "nvim_buf_get_lines",
		params: []any{0, start, end, false},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer name: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return []string{}, fmt.Errorf("Error executing lua: %v\n", err)
	}

	bufferLines, err := anyx.ToSliceOf[string](result)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer lines return value: %w", err)
	}

	return bufferLines, nil
}

func (n *Nvim) DeleteBuffer() error {
	request := requestMessage{
		method: "nvim_buf_delete",
		params: []any{0, struct{ force bool }{force: true}}}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to delete buffer: %v\n", err)
	}

	return nil
}

/*
	 func (n *Nvim) ApiInfo() (any, error) {
		request := requestMessage{
			method: "nvim_get_api_info",
			params: []any{},
		}
		response, err := n.rpc.Send(request)

		if err != nil {
			return nil, fmt.Errorf("Error getting API info: %v\n", err)
		}

		result, err := response.Result()

		if err != nil {
			return nil, fmt.Errorf("Error getting API info: %v\n", err)
		}

		return result, nil
	}
*/

func (n *Nvim) CallFunction(function string, functionArgs []any) (any, error) {
	request := requestMessage{
		method: "nvim_call_function",
		params: []any{function, functionArgs},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return nil, fmt.Errorf("Error executing function: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return nil, fmt.Errorf("Error executing function: %v\n", err)
	}

	return result, nil
}

func (n *Nvim) ExecLua(lua string, args []any) (any, error) {
	request := requestMessage{
		method: "nvim_exec_lua",
		params: []any{lua, args},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return nil, fmt.Errorf("Error executing lua: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return nil, fmt.Errorf("Error executing lua: %v\n", err)
	}

	return result, nil
}

func (n *Nvim) startTS() error {
	luaCode := `
		if vim.g.lua_ts_ready == true then return end

		vim.g.lua_ts_ready = false

		local ok = pcall(vim.treesitter.language.add, "lua")

		if not ok then error("Treesitter Lua parser registration failed.") end

		local ok = pcall(vim.treesitter.language.add, "luadoc")

		if not ok then error("Treesitter Luadoc parser registration failed.") end

		vim.api.nvim_create_autocmd("FileType", {
			pattern = { "lua" },
			callback = function(opts)
				vim.treesitter.start(opts.buf, "luadoc")
				vim.treesitter.start(opts.buf, "lua")
			end,
		})

		vim.treesitter.start(0, 'luadoc')
		vim.treesitter.start(0, 'lua')

		vim.g.lua_ts_ready = true
	`

	_, err := n.ExecLua(luaCode, []any{30000})

	if err != nil {
		return fmt.Errorf("Error starting treesitter lua: %w", err)
	}

	return nil
}

func (n *Nvim) GetRuntimeType(variable string) (string, error) {
	runtimePath := variable
	parts := strings.Split(runtimePath, ".")

	switch len(parts) {
	case 1:
	default:
		tail := parts[len(parts)-1]
		// https://www.lua.org/manual/5.1/manual.html#2.1
		switch tail {
		case "and", "break", "do", "else", "elseif", "end", "false", "for", "function", "if", "in", "local", "nil", "not", "or", "repeat", "return", "then", "true", "until", "while":
			head := parts[:len(parts)-1]
			runtimePath = strings.Join(head, ".") + "['" + tail + "']"
		}
	}

	luaCode := fmt.Sprintf("return type(%s)", runtimePath)

	result, err := n.ExecLua(luaCode, []any{})

	if err != nil {
		return "", fmt.Errorf("Error getting the type of %s: %w", variable, err)
	}

	typeName, ok := result.(string)

	if !ok {
		return "", fmt.Errorf("Error getting the type of %s: Error converting the result to a string", runtimePath)
	}

	return typeName, nil
}

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

func New(config Config, opts ...optionProvider) (*Nvim, error) {
	options, err := NewOptions(config, opts...)

	if err != nil {
		return nil, fmt.Errorf("Error getting nvim options: %v", err)
	}

	arguments := []string{"--embed", "--headless", "-u", options.config.InitFile()}
	arguments = append(arguments, options.arguments...)
	cmd := exec.Command(options.command, arguments...)
	rpc, err := NewRpc(cmd)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim rpc: %v", err)
	}

	return &Nvim{
		options: options,
		cmd:     cmd,
		rpc:     rpc,
	}, nil
}
