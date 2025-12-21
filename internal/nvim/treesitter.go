package nvim

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/ts"
	"github.com/Tolomeo/anydev.nvim/internal/utils/anyx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type TsQueryConfig struct {
	Language string
	Query    string
	Range    *ts.LineRange
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

func (n *Nvim) execTsQuery(config TsQueryConfig) (*[]ts.Capture, error) {
	err := n.startTS()

	if err != nil {
		return nil, err
	}

	// Query injections are only recalculated when the buffer changes
	// so we make a (hopefully) inhert change to force their presence
	// by adding an empty line at the end of the buffer text
	luaCode := `
		local keys = vim.api.nvim_replace_termcodes("Go<Esc>", true, false, true)
		vim.api.nvim_feedkeys(keys, "n", false)
	`

	_, err = n.ExecLua(luaCode, []any{})

	if err != nil {
		return nil, err
	}

	luaCode = `
		local args = { ... }
		local query_language = args[1]
		local query = args[2]
		local start, stop = arg[3], arg[4]

		local lua = "lua"
		local luadoc = "luadoc"
		local bufnr = 0

		local parser = vim.treesitter.get_parser(bufnr, lua)

		if not parser then
			error("Error: could not initialize lua parser")
		end

		parser:add_child(luadoc)

		local childParser = parser:children()[luadoc]

		if not childParser then 
			error("Error: could not initialize luadoc parser")
		end

		parser:parse(true)

		local parsedQuery = vim.treesitter.query.parse(query_language, query)

		local queryResult = {}

		parser:for_each_tree(function(tree, language_tree)
			local lang = language_tree:lang()

			if query_language == lang then
				for id, node, _ in parsedQuery:iter_captures(tree:root(), bufnr, start, stop) do
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
			end
		end)

		if not next(queryResult) then
			return vim.NIL
		end

		return vim.fn.json_encode(queryResult)
	`

	luaArgs := []any{config.Language, config.Query}

	if config.Range != nil {
		luaArgs = append(luaArgs, config.Range.Start, config.Range.End)
	}

	result, err := n.ExecLua(luaCode, luaArgs)

	// fmt.Printf("\nQuery: \n%v\n%v\n%v\n", config.Query, result, err)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Nvim TSQuery error: %w", err)
	case result == nil:
		return nil, nil
	}

	stringResult, ok := result.(string)

	if !ok {
		return nil, fmt.Errorf("Error reading tsNodes query result as a string: %v", result)
	}

	var captures []ts.Capture

	if err := json.Unmarshal([]byte(stringResult), &captures); err != nil {
		return nil, fmt.Errorf("Error decoding tsNodes json response: %w", err)
	}

	return &captures, nil
}

func (n *Nvim) TsQueryAll(config TsQueryConfig) (*[][]ts.Capture, error) {
	queryAllConfig := TsQueryConfig{
		Language: config.Language,
		Query:    fmt.Sprintf("(%s) @tsquery.match", config.Query),
		Range:    config.Range,
	}

	captures, err := n.execTsQuery(queryAllConfig)

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	/* fmt.Println("query all captures")
	fmt.Printf("\n\n%+v\n\n", captures) */

	queryCaptures, _ := slicesx.FilterFunc(*captures, func(capture ts.Capture) (bool, error) {
		return (capture.Id == "tsquery.match"), nil
	})

	queryMatches, _ := slicesx.MapFunc(queryCaptures, func(queryCapture ts.Capture) ([]ts.Capture, error) {
		return slicesx.FilterFunc(*captures, func(capture ts.Capture) (bool, error) {
			if capture.Id == queryCapture.Id {
				return false, nil
			}

			return queryCapture.Node.Contains(capture.Node), nil
		})
	})

	/* fmt.Println("query all matches")
	fmt.Printf("\n\n%+v\n\n", queryMatches) */

	return &queryMatches, nil
}

func (n *Nvim) TsQueryOne(config TsQueryConfig) (*[]ts.Capture, error) {
	matches, err := n.TsQueryAll(config)

	switch {
	case err != nil:
		return nil, err
	case matches == nil:
		return nil, nil
	case len(*matches) > 1:
		return nil, fmt.Errorf("Error executing TSQuery: too many matches, expected 1 but received %d", len(*matches))
	case len(*matches) < 1:
		return nil, nil
	}

	/* fmt.Println("query one")
	fmt.Printf("\n\n%+v\n\n", matches) */

	match := (*matches)[0]

	return &match, nil
}

var ErrSafeTSQueryNoMatch = errors.New("The parsed language tree contains errors")

func (n *Nvim) SafeTsQuery(config TsQueryConfig) (*[]ts.Capture, error) {
	captures, err := n.TsQueryOne(config)

	switch {
	case err != nil:
		return captures, err
	case captures == nil:
		return nil, nil
	}

	// TODO: pass range from capture
	errorCaptures, err := n.TsQueryOne(TsQueryConfig{Language: config.Language, Query: `(ERROR) @syntax.error`})

	switch {
	case err != nil:
		return nil, err
	case errorCaptures != nil:
		return nil, ErrSafeTSQueryNoMatch
	}

	return captures, nil
}

func (n *Nvim) GetTSNodeAt(nodeTypes []string, line uint, character uint) (*ts.TsNode, error) {
	err := n.startTS()

	if err != nil {
		return nil, err
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

	switch {
	case err != nil:
		return nil, err
	case result == nil:
		return nil, nil
	}

	stringResult, ok := result.(string)

	if !ok {
		return nil, fmt.Errorf("Error converting result into string: %v", result)
	}

	tsNode := ts.TsNode{}
	err = tsNode.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling tsnode response: %w", err)
	}

	return &tsNode, nil
}

func (n *Nvim) GetCommentBlockAt(line uint, character uint) (*[]string, error) {
	lines, err := n.GetBufferLines(int(line)-1, int(line))

	switch {
	case err != nil:
		return nil, err
	case len(lines) < 1:
		return nil, nil
	}

	// clamping the received character to be inside the line
	character = max(0, min(character, uint(len(lines[0])-1)))

	node, err := n.GetTSNodeAt([]string{ts.COMMENT}, line, character)

	switch {
	case err != nil:
		return nil, err
	case node == nil:
		return nil, nil
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

	result, err := n.ExecLua(luaCode, []any{node.Range.Start.Line, node.Range.End.Line})

	if err != nil {
		return nil, err
	}

	bufferLines, err := anyx.ToSliceOf[string](result)

	if err != nil {
		return nil, fmt.Errorf("Error reading buffer lines return value: %w", err)
	}

	return &bufferLines, nil
}
