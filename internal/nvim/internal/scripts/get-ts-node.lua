local args = { ... }
local node_types = args[1]
local line = args[2]
local character = args[3]
local bufnr = 0

local parser = _G.Anydev:get_ts_parser(bufnr)
local node = nil

parser.parser:for_each_tree(function(_, language_tree)
	if node ~= nil then
		return
	end

	local tree_node = language_tree:named_node_for_range({ line, character, line, character })

	while tree_node ~= nil do
		if vim.tbl_contains(node_types, tree_node:type()) then
			node = tree_node
			break
		end

		tree_node = tree_node:parent()
	end
end)

if node == nil then
	return vim.NIL
end

local result_node = _G.Anydev:get_ts_parser_node(node, parser.source)
return vim.fn.json_encode(result_node)
