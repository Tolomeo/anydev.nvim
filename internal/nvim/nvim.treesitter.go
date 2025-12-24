package nvim

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type TsQueryConfig struct {
	Language string
	Query    string
	Range    *treesitter.LineRange
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

	// Query injections are only recalculated when the buffer changes
	// so we make a (hopefully) inhert change to force their presence
	// by adding an empty line at the end of the buffer text
	luaCode = `
		local keys = vim.api.nvim_replace_termcodes("Go<Esc>", true, false, true)
		vim.api.nvim_feedkeys(keys, "n", false)
	`

	_, err = n.ExecLua(luaCode, []any{})

	if err != nil {
		return err
	}

	return nil
}

func (n *Nvim) execTsQuery(config TsQueryConfig) (*[]treesitter.Capture, error) {
	err := n.startTS()

	if err != nil {
		return nil, err
	}

	luaCode := `
		local args = { ... }
		local query_language = args[1]
		local query = args[2]
		local start, stop = args[3], args[4]

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
		luaArgs = append(luaArgs, config.Range.Start, config.Range.End+1)
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

	var captures []treesitter.Capture

	if err := json.Unmarshal([]byte(stringResult), &captures); err != nil {
		return nil, fmt.Errorf("Error decoding tsNodes json response: %w", err)
	}

	return &captures, nil
}

type TsQueryMatch []treesitter.Capture

func (n *Nvim) TsQueryAll(config TsQueryConfig) (*[]TsQueryMatch, error) {
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

	queryCaptures, _ := slicesx.FilterFunc(*captures, func(capture treesitter.Capture) (bool, error) {
		return (capture.Id == "tsquery.match"), nil
	})

	queryMatches, _ := slicesx.MapFunc(queryCaptures, func(queryCapture treesitter.Capture) (TsQueryMatch, error) {
		return slicesx.FilterFunc(*captures, func(capture treesitter.Capture) (bool, error) {
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

func (n *Nvim) TsQueryOne(config TsQueryConfig) (*TsQueryMatch, error) {
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

func capturesLineRange(captures []treesitter.Capture) *treesitter.LineRange {
	startLines, _ := slicesx.MapFunc(captures, func(capture treesitter.Capture) (float64, error) {
		return capture.Node.Range.Start.Line, nil
	})
	endLines, _ := slicesx.MapFunc(captures, func(capture treesitter.Capture) (float64, error) {
		return capture.Node.Range.End.Line, nil
	})

	return &treesitter.LineRange{
		Start: slices.Min(startLines),
		End:   slices.Max(endLines),
	}
}

type SafeTsQueryResult struct {
	HasError bool
	Captures TsQueryMatch
}

func (n *Nvim) SafeTsQueryOne(config TsQueryConfig) (*SafeTsQueryResult, error) {
	captures, err := n.TsQueryOne(config)

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	errorCaptures, err := n.TsQueryAll(TsQueryConfig{
		Language: config.Language,
		Query:    `(ERROR) @tsquery.error`,
		Range:    capturesLineRange(*captures),
	})

	switch {
	case err != nil:
		return nil, err
	case errorCaptures != nil:
		return &SafeTsQueryResult{
			HasError: true,
			Captures: *captures,
		}, nil
	default:
		return &SafeTsQueryResult{
			HasError: false,
			Captures: *captures,
		}, nil
	}
}

func (n *Nvim) SafeTsQueryAll(config TsQueryConfig) (*[]SafeTsQueryResult, error) {
	matches, err := n.TsQueryAll(config)

	switch {
	case err != nil:
		return nil, err
	case matches == nil:
		return nil, nil
	}

	results := []SafeTsQueryResult{}

	for _, matchCaptures := range *matches {
		errorCaptures, err := n.TsQueryAll(TsQueryConfig{
			Language: config.Language,
			Query:    `(ERROR) @tsquery.error`,
			Range:    capturesLineRange(matchCaptures),
		})

		// fmt.Printf("\n%+v\n\n", errorCaptures)

		switch {
		case err != nil:
			return nil, err
		case errorCaptures != nil:
			results = append(results, SafeTsQueryResult{
				HasError: true,
				Captures: matchCaptures,
			})
		default:
			results = append(results, SafeTsQueryResult{
				HasError: false,
				Captures: matchCaptures,
			})
		}
	}

	return &results, nil
}

func (n *Nvim) GetTSNodeAt(nodeTypes []string, line uint, character uint) (*treesitter.TsNode, error) {
	err := n.startTS()

	if err != nil {
		return nil, err
	}

	luaCode := `
		local args = { ... }
		local ancestorNodeTypes = args[1]
		local line = args[2]
		local character = args[3]

		local lua = "lua"
		local luadoc = "luadoc"
		local bufnr = 0

		local root_parser = vim.treesitter.get_parser(bufnr, lua)

		if not root_parser then
			error("Error: could not initialize lua parser")
		end

		root_parser:add_child(luadoc)

		local childParser = root_parser:children()[luadoc]

		if not childParser then
			error("Error: could not initialize luadoc parser")
		end

		root_parser:parse(true)

		local node = nil

		root_parser:for_each_tree(function(_, parser)
			if node ~= nil then
				return
			end

			local tree_node = parser:named_node_for_range({ line, character, line, character })

			while tree_node ~= nil do
				if vim.tbl_contains(ancestorNodeTypes, tree_node:type()) then
					node = tree_node
					break
				end

				tree_node = tree_node:parent()
			end
		end)

		if node == nil then
			return vim.NIL
		end

		local nodeType = node:type()
		local startLine, startCharacter, endLine, endCharacter = node:range(false)
		local text = vim.treesitter.get_node_text(node, bufnr)

		return vim.fn.json_encode({
			type = nodeType,
			range = {
				start = { line = startLine, character = startCharacter },
				["end"] = { line = endLine, character = endCharacter },
			},
			text = text,
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

	tsNode := treesitter.TsNode{}
	err = tsNode.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling tsnode response: %w", err)
	}

	return &tsNode, nil
}
