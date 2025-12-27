local args = { ... }
local startLine = args[1]
local endLine = args[2]

local function is_comment(ln)
	local rest = ln:match("^%s*(.*)")
	return rest:sub(1, 2) == "--"
end

local previous_line = vim.api.nvim_buf_get_lines(0, startLine - 1, startLine, false)[1]

while previous_line and is_comment(previous_line) do
	startLine = startLine - 1
	previous_line = vim.api.nvim_buf_get_lines(0, startLine - 1, startLine, false)[1]
end

local next_line = vim.api.nvim_buf_get_lines(0, endLine + 1, endLine + 2, false)[1]

while next_line and is_comment(next_line) do
	endLine = endLine + 1
	next_line = vim.api.nvim_buf_get_lines(0, endLine + 1, endLine + 2, false)[1]
end

return vim.api.nvim_buf_get_lines(0, startLine, endLine + 1, true)
