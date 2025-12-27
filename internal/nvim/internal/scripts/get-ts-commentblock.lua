local args = { ... }
local startLine = args[1]
local endLine = args[2]

local bufnr = 0

local function is_comment(ln)
	local rest = ln:match("^%s*(.*)")
	return rest:sub(1, 2) == "--"
end

local previous_line = vim.api.nvim_buf_get_lines(bufnr, startLine - 1, startLine, false)[1]

while previous_line and is_comment(previous_line) do
	startLine = startLine - 1
	previous_line = vim.api.nvim_buf_get_lines(bufnr, startLine - 1, startLine, false)[1]
end

local next_line = vim.api.nvim_buf_get_lines(bufnr, endLine + 1, endLine + 2, false)[1]

while next_line and is_comment(next_line) do
	endLine = endLine + 1
	next_line = vim.api.nvim_buf_get_lines(bufnr, endLine + 1, endLine + 2, false)[1]
end

local comment_lines = vim.api.nvim_buf_get_lines(bufnr, startLine, endLine + 1, true)

if #comment_lines < 1 then
	return vim.NIL
end

return vim.fn.json_encode({
	type = "comment_block",
	range = {
		start = { line = startLine, character = 0 },
		["end"] = { line = endLine, character = #(comment_lines[#comment_lines]) },
	},
	text = table.concat(comment_lines, "\n"),
})
