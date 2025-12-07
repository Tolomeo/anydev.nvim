package nvim

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/ts"
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

var ErrNotFound = errors.New("Not found")

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
				vim.treesitter.start(opts.buf, "lua")
			end,
		})

		vim.treesitter.start(0, 'lua')

		vim.g.lua_ts_ready = true
	`

	_, err := n.ExecLua(luaCode, []any{30000})

	if err != nil {
		return fmt.Errorf("Error starting treesitter lua: %w", err)
	}

	return nil
}

func (n *Nvim) GetLuaTypeName(variable string) (string, error) {
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

type TsQueryConfig struct {
	Language      string
	Query         string
	WithoutErrors bool
}

func (n *Nvim) TsQuery(config TsQueryConfig) ([]ts.Capture, error) {
	err := n.startTS()

	if err != nil {
		return []ts.Capture{}, err
	}

	luaCode := `
		local args = { ... }
		local language = args[1]
		local query = args[2]
		local bufnr = 0

		local parser = vim.treesitter.get_parser(bufnr, language)

		if not parser then
			error("Error: Treesitter parser not found")
		end

		local tree = parser:parse()[1]

		if not tree then
			error("Error: Could not parse the buffer content into a Tree-sitter tree.")
		end

		local parsedQuery = vim.treesitter.query.parse(language, query)

		local queryResult = {}

		for id, node in parsedQuery:iter_captures(tree:root(), bufnr) do
			local captureId = parsedQuery.captures[id]

			local nodeType = node:type()
			local startLine, startCharacter, endLine, endCharacter = node:range()
	    local text = vim.treesitter.get_node_text(node, bufnr)

			local tsNode = {
				type = nodeType,
				range = {
					start = { line = startLine, character = startCharacter },
					["end"] = { line = endLine, character = endCharacter },
				},
				text = text,
			}

			local capture = {
				id = captureId,
				node = tsNode,
			}

			table.insert(queryResult, capture)
		end

		if not next(queryResult) then
			return vim.NIL
		end

		return vim.fn.json_encode(queryResult)
	`

	if config.WithoutErrors {
		result, err := n.ExecLua(luaCode, []any{config.Language, `(ERROR) @error`})

		if err != nil {
			return []ts.Capture{}, fmt.Errorf("Nvim TSQuery error: %w", err)
		}

		if result != nil {
			return []ts.Capture{}, ErrNotFound
		}
	}

	result, err := n.ExecLua(luaCode, []any{config.Language, config.Query})

	if err != nil {
		return []ts.Capture{}, fmt.Errorf("Nvim TSQuery error: %w", err)
	}

	if result == nil {
		return []ts.Capture{}, ErrNotFound
	}

	stringResult, ok := result.(string)

	if !ok {
		return []ts.Capture{}, fmt.Errorf("Error reading tsNodes query result as a string: %v", result)
	}

	var capturedTsNodes []ts.Capture

	if err := json.Unmarshal([]byte(stringResult), &capturedTsNodes); err != nil {
		return []ts.Capture{}, fmt.Errorf("Error decoding tsNodes json response: %w", err)
	}

	return capturedTsNodes, nil
}

func (n *Nvim) ReadTsNode(node ts.TsNode) (string, error) {
	startRow, startCol, endRow, endCol := int(node.Range.Start.Line), int(node.Range.Start.Character), int(node.Range.End.Line), int(node.Range.End.Character)
	textContent, err := n.GetBufferText(startRow, startCol, endRow, endCol)

	if err != nil {
		return "", err
	}

	return strings.Join(textContent, ""), nil
}

func (n *Nvim) GetTSNodeAt(nodeTypes []string, line uint, character uint) (ts.TsNode, error) {
	tsNode := ts.TsNode{}

	err := n.startTS()

	if err != nil {
		return tsNode, err
	}

	luaCode := `
		local args = { ... }
		local ancestorNodeTypes = args[1]
		local line = args[2]
		local character = args[3]
		local bufnr = 0

		local parser = vim.treesitter.get_parser(bufnr, "lua")
		local root = parser:parse()[1]:root()

		local node = root:descendant_for_range(line, character, line, character)

		if node == nil then
			return vim.NIL
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
			return vim.NIL
		end

		local nodeType = node:type()
		local startLine, startCharacter, endLine, endCharacter = targetNode:range(false)
	  local text = vim.treesitter.get_node_text(node, bufnr)

		return vim.fn.json_encode({
			type = nodeType,
			range = {
				start = { line = startLine, character = startCharacter },
				["end"] = { line = endLine, character = endCharacter },
			},
			text = text
		})
	`

	result, err := n.ExecLua(luaCode, []any{nodeTypes, line, character})

	if err != nil {
		return tsNode, err
	}

	if result == nil {
		return tsNode, ErrNotFound
	}

	stringResult, ok := result.(string)

	if !ok {
		return tsNode, fmt.Errorf("Error converting result into string: %v", result)
	}

	err = tsNode.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return tsNode, fmt.Errorf("Error unmarshalling tsnode response: %w", err)
	}

	return tsNode, nil
}

func (n *Nvim) ReadCommentBlockAt(cursorPosition CursorPosition) ([]string, error) {
	tsNode, err := n.GetTSNodeAt([]string{ts.COMMENT}, cursorPosition.Line, cursorPosition.Character)

	if err != nil {
		return []string{}, err
	}

	luaCode := `
		local args = {...}
		local startLine = args[1]
		local endLine = args[2]

		local function is_comment(ln)
			local rest = ln:match("^%s*(.*)")
			return rest:sub(1,2) == "--"
		end

		local previous_line = vim.api.nvim_buf_get_lines(0, startLine -1, startLine, false)[1]

		while previous_line and is_comment(previous_line) do
			startLine = startLine - 1
			previous_line = vim.api.nvim_buf_get_lines(0, startLine -1, startLine, false)[1]
		end

		local next_line = vim.api.nvim_buf_get_lines(0, endLine + 1, endLine + 2, false)[1]

		while next_line and is_comment(next_line) do
			endLine = endLine + 1
			next_line = vim.api.nvim_buf_get_lines(0, endLine + 1, endLine + 2, false)[1]
		end

		return vim.api.nvim_buf_get_lines(0, startLine, endLine + 1, true)
	`

	result, err := n.ExecLua(luaCode, []any{tsNode.Range.Start.Line, tsNode.Range.End.Line})

	if err != nil {
		return []string{}, err
	}

	bufferLines, err := anyx.ToSliceOf[string](result)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer lines return value: %w", err)
	}

	return bufferLines, nil
}

func (n *Nvim) ReadTSNodeAt(cursorPosition CursorPosition, nodeType string, nodeTypes ...string) ([]string, error) {
	tsNodeTypes := append([]string{nodeType}, nodeTypes...)
	tsNode, err := n.GetTSNodeAt(tsNodeTypes, cursorPosition.Line, cursorPosition.Character)

	if err != nil {
		return []string{}, err
	}

	bufferLines, err := n.GetBufferLines(int(tsNode.Range.Start.Line), int(tsNode.Range.End.Line+1))

	if err != nil {
		return []string{}, err
	}

	return bufferLines, nil
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

func (n *Nvim) GetDefinitionLocation(line uint, character uint) ([]Location, error) {
	lspDefinitions, err := n.GetLSPDefinition(line, character)

	switch {
	case err != nil:
		return []Location{}, err
	case len(lspDefinitions) == 0:
		return []Location{}, ErrNotFound
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
		return []Location{}, err
	}

	return locations, nil
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
