local args = { ... }
---@type anydev_ts_parser_node
local node = vim.fn.json_decode(args[1])
local bufnr = 0

---@type string[]
local annotations = {}

local current_line = node.range.start.line
local previous_line = vim.api.nvim_buf_get_lines(bufnr, current_line - 1, current_line, false)[1]
while previous_line and _G.Anydev:is_comment(previous_line) do
	table.insert(annotations, 1, previous_line)
	current_line = current_line - 1
	previous_line = vim.api.nvim_buf_get_lines(bufnr, current_line - 1, current_line, false)[1]
end

local trailing_characters = vim.api.nvim_buf_get_text(
	bufnr,
	node.range["end"].line,
	node.range["end"].character,
	node.range["end"].line,
	-1,
	{}
)[1]
local is_comment, comment = _G.Anydev:is_comment(trailing_characters)
if is_comment then
	table.insert(annotations, comment)
end

return vim.fn.json_encode(annotations)
