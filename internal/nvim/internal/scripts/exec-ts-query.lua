local args = { ... }
local query_language = args[1]
local query = args[2]
local start, stop = args[3], args[4]
local bufnr = 0

local parser = _G.Anydev:get_ts_parser(bufnr)
local parsedQuery = vim.treesitter.query.parse(query_language, query)
local queryResult = {}

parser.parser:for_each_tree(function(tree, language_tree)
	local lang = language_tree:lang()

	if query_language == lang then
		for id, node, _ in parsedQuery:iter_captures(tree:root(), bufnr, start, stop) do
			local capture_id = parsedQuery.captures[id]
			local capture_node = _G.Anydev:get_ts_parser_node(node, parser.source)

			local capture = {
				id = capture_id,
				node = capture_node,
			}

			table.insert(queryResult, capture)
		end
	end
end)

if not next(queryResult) then
	return vim.NIL
end

return vim.fn.json_encode(queryResult)
