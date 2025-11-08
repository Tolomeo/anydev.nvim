package nvim

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
)

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
		params: []any{0, 0, -1, false, lines},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying set buffer lines: %v\n", err)
	}

	return nil

}

func (n *Nvim) GetBufferLines() ([]string, error) {
	request := requestMessage{
		method: "nvim_buf_get_lines",
		params: []any{0, 0, -1, false},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer name: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return []string{}, fmt.Errorf("Error executing lua: %v\n", err)
	}

	sliceOfAny, ok := result.([]any)

	if !ok {
		return []string{}, fmt.Errorf("Error reading buffer lines return value: $v")
	}

	sliceOfStrings := make([]string, len(sliceOfAny))

	for index, value := range sliceOfAny {
		str, ok := value.(string)

		if !ok {
			return []string{}, fmt.Errorf("Error reading buffer lines return value item %d: %s", index, value)
		}

		sliceOfStrings[index] = str
	}

	return sliceOfStrings, nil
}

func (n *Nvim) DeleteBuffer() error {
	request := requestMessage{
		method: "nvim_buf_delete",
		params: []any{0, struct{ force bool }{force: true}}}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to write buffer: %v\n", err)
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

type CursorPosition struct {
	Line      uint
	Character uint
}

func (n *Nvim) GetTSNodeAncestorAt(ancestorType []string, cursorPosition CursorPosition) error {
	err := n.StartTS()

	if err != nil {
		return err
	}

	luaCode := `
		local args = {...}
		local ancestorNodeTypes = args[1]
		local line = args[2]
		local character = args[3]

		local parser = vim.treesitter.get_parser(0, "lua")
		local root = parser:parse()[1]:root()

		local node = root:descendant_for_range(line, character, line, character)

		if node == nil then
			error("No node found")
		end

		local targetNode = nil

		while not node:equal(root) do
			if vim.tbl_contains(ancestorNodeTypes, node:type()) then
				targetNode = node
				break
			end

			node = node:parent()
		end

		if targetNode == nil then
			error("No node found")
		end

		local startLine, _, endLine, _ = targetNode:range()
		local previous_line = vim.api.nvim_buf_get_lines(0, startLine -1, startLine, false)[1]

		while previous_line and #previous_line > 0 and previous_line:find("^%s*--") do
			startLine = startLine -1
			previous_line = vim.api.nvim_buf_get_lines(0, startLine -1, startLine, false)[1]
		end

		return vim.fn.json_encode({ result = vim.api.nvim_buf_get_lines(0, startLine, endLine + 1, true) })
	`

	result, err := n.ExecLua(luaCode, []any{ancestorType, cursorPosition.Line, cursorPosition.Character})

	if err != nil {
		return err
	}

	stringResult, ok := result.(string)

	if !ok {
		return fmt.Errorf("Error reading tsparent response: %v", result)
	}

	fmt.Println("TSParent")
	fmt.Println(stringResult)

	return nil
}

func (n *Nvim) GetAnnotatedFunctionBufferLinesAt(file string, line uint, character uint) error {
	err := n.StartTS()

	if err != nil {
		return err
	}

	_, err = n.Open(file)

	if err != nil {
		return err
	}

	luaCode := `
		local args = {...}
		local line = args[1]
		local column = args[2]

		local parser = vim.treesitter.get_parser(0, 'lua')
		local tree = parser:parse()[1]
		local root = tree:root()

		local node = vim.treesitter.get_node({line, column}, 0)

		return node:type()
	`

	result, err := n.ExecLua(luaCode, []any{line, character})

	if err != nil {
		return err
	}

	stringResult, ok := result.(string)

	if !ok {
		return fmt.Errorf("Error reading function buffer lines response: %v", result)
	}

	fmt.Println(stringResult)

	return nil
}

func (n *Nvim) StartTS() error {
	luaCode := `
		if vim.g.lua_ts_ready == true then return end

		vim.g.lua_ts_ready = false

		local ok = pcall(vim.treesitter.language.add, "lua")

		if not ok then error("Treesitter Lua parser registration failed.") end

		vim.api.nvim_create_autocmd("FileType", {
			pattern = { "lua" },
			callback = function(opts)
				vim.treesitter.start(opts.buf, "lua")
			end,
		})

		vim.g.lua_ts_ready = true
	`

	_, err := n.ExecLua(luaCode, []any{30000})

	if err != nil {
		return fmt.Errorf("Error starting treesitter lua: %w", err)
	}

	return nil
}

func (n *Nvim) StartLSP() error {
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

func (n *Nvim) GetLSPDocumentSymbols() (lsp.TextDocumentDocumentSymbolResponse, error) {
	documentSymbols := lsp.TextDocumentDocumentSymbolResponse{}

	err := n.StartLSP()

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

func (n *Nvim) GetLSPHover(line uint, character uint) (lsp.TextDocumentHoverResponse, error) {
	hover := lsp.TextDocumentHoverResponse{}

	err := n.StartLSP()

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

	result, err := n.ExecLua(luaCode, []any{line, character, 2000})

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

func (n *Nvim) GetLSPDeclaration(line uint, character uint) error {
	err := n.StartLSP()

	if err != nil {
		return err
	}

	luaCode := `
		local args = {...}

		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local positionParams = {line = args[1], character = args[2]}
		local result = vim.lsp.buf_request_sync(0, 'textDocument/declaration', { textDocument = textDocumentParams, position = positionParams }, 2000)
		return vim.fn.json_encode(result[1])
	`

	result, err := n.ExecLua(luaCode, []any{line, character})

	if err != nil {
		return fmt.Errorf("Error getting completion: %v", err)
	}

	stringResult, ok := result.(string)

	if !ok {
		return fmt.Errorf("Error reading completion response: %v", result)
	}

	fmt.Println(stringResult)

	return nil
}

func (n *Nvim) GetLSPDefinition(line uint, character uint) (lsp.TextDocumentDefinitionResponse, error) {
	definition := lsp.TextDocumentDefinitionResponse{}

	err := n.StartLSP()

	if err != nil {
		return definition, err
	}

	luaCode := `
		local args = {...}
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local positionParams = {line = args[1], character = args[2]}
		local result = vim.lsp.buf_request_sync(0, 'textDocument/definition', { textDocument = textDocumentParams, position = positionParams }, 2000)
		return vim.fn.json_encode(result[1])
	`

	result, err := n.ExecLua(luaCode, []any{line, character})

	if err != nil {
		return definition, fmt.Errorf("Error getting completion: %v", err)
	}

	stringResult, ok := result.(string)

	if !ok {
		return definition, fmt.Errorf("Error reading completion response: %v", result)
	}

	err = definition.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return definition, fmt.Errorf("Error unmarshalling definition response: %w", err)
	}

	return definition, nil
}

func (n *Nvim) GetLSPImplementation(line uint, character uint) error {
	err := n.StartLSP()

	if err != nil {
		return err
	}

	luaCode := `
		local args = {...}
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local positionParams = {line = args[1], character = args[2]}
		local result = vim.lsp.buf_request_sync(0, 'textDocument/implementation', { textDocument = textDocumentParams, position = positionParams }, 2000)
		return vim.fn.json_encode(result[1])
	`

	result, err := n.ExecLua(luaCode, []any{line, character})

	if err != nil {
		return fmt.Errorf("Error getting completion: %v", err)
	}

	stringResult, ok := result.(string)

	if !ok {
		return fmt.Errorf("Error reading completion response: %v", result)
	}

	fmt.Println(stringResult)

	return nil
}

func (n *Nvim) GetLSPTypeDefinition(line uint, character uint) error {
	err := n.StartLSP()

	if err != nil {
		return err
	}

	luaCode := `
		local args = {...}
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local positionParams = {line = args[1], character = args[2]}
		local result = vim.lsp.buf_request_sync(0, 'textDocument/typeDefinition', { textDocument = textDocumentParams, position = positionParams }, 2000)
		return vim.fn.json_encode(result[1])
	`

	result, err := n.ExecLua(luaCode, []any{line, character})

	if err != nil {
		return fmt.Errorf("Error getting completion: %v", err)
	}

	stringResult, ok := result.(string)

	if !ok {
		return fmt.Errorf("Error reading completion response: %v", result)
	}

	fmt.Println(stringResult)

	/* err = completion.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return completion, fmt.Errorf("Error unmarshalling completion response: %v", err)
	}

	return completion, nil */
	return nil
}

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

func (n *Nvim) GetLSPCompletion(line uint, character uint) (lsp.TextDocumentCompletionResponse, error) {
	completion := lsp.TextDocumentCompletionResponse{}

	err := n.StartLSP()

	if err != nil {
		return completion, err
	}

	luaCode := `
		local args = {...}
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local positionParams = {line = args[1], character = args[2]}
		local result = vim.lsp.buf_request_sync(0, 'textDocument/completion', { textDocument = textDocumentParams, position = positionParams }, 2000)
		return vim.fn.json_encode(result[1])
	`

	result, err := n.ExecLua(luaCode, []any{line, character})

	if err != nil {
		return completion, fmt.Errorf("Error getting completion: %v", err)
	}

	stringResult, ok := result.(string)

	if !ok {
		return completion, fmt.Errorf("Error reading completion response: %v", result)
	}

	err = completion.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return completion, fmt.Errorf("Error unmarshalling completion response: %v", err)
	}

	return completion, nil
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

func New(opts ...optionProvider) (*Nvim, error) {
	options, err := NewOptions(opts...)

	if err != nil {
		return nil, fmt.Errorf("Error getting nvim options: %v", err)
	}

	arguments := []string{"--embed", "--headless", "-u", options.config.InitFile()}
	arguments = append(arguments, options.arguments...)
	cmd := exec.Command(options.cmd, arguments...)
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
