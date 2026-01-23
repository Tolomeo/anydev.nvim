local args = { ... }
local query_language = args[1]
local query = args[2]
local start, stop = args[3], args[4]
local bufnr = 0

local lua = "lua"
local luadoc = "luadoc"

local parser = vim.treesitter.get_parser(bufnr, lua)

if not parser then
	error("Error: could not initialize lua parser")
end

parser:add_child(luadoc)

local childParser = parser:children()[luadoc]

if not childParser then
	error("Error: could not initialize luadoc parser")
end

local parsedQuery = vim.treesitter.query.parse(query_language, query)
local queryResult = {}

parser:parse(true)
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
