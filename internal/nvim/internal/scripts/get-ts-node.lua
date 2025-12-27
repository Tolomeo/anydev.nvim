local args = { ... }
local node_types = args[1]
local line = args[2]
local character = args[3]

local lua = "lua"
local luadoc = "luadoc"
local bufnr = 0

local function get_tsnode() end

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
