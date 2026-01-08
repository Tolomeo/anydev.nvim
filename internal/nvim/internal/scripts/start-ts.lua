local args = { ... }
local delay = args[1]

local bufnr = 0
local lua = "lua"
local luadoc = "luadoc"

if vim.g.ts_ready ~= true then
	local ok = pcall(vim.treesitter.language.add, "lua")

	if not ok then
		error("Treesitter Lua parser registration failed.")
	end

	ok = pcall(vim.treesitter.language.add, "luadoc")

	if not ok then
		error("Treesitter Luadoc parser registration failed.")
	end

	vim.g._ts_force_sync_parsing = true
	vim.g.ts_ready = true
end

vim.treesitter.start(bufnr, luadoc)
vim.treesitter.start(bufnr, lua)

-- Query injections are only recalculated when the buffer changes
-- so we make a (hopefully) inhert change to force their presence
-- by adding an empty line at the end of the buffer text
local add_empty_line = vim.api.nvim_replace_termcodes("Go<Esc>", true, false, true)
vim.api.nvim_feedkeys(add_empty_line, "n", false)

vim.wait(1, function()
	return false
end, 1, false)

-- Removing the added line
local remove_line = vim.api.nvim_replace_termcodes("Gdd<Esc>", true, false, true)
vim.api.nvim_feedkeys(remove_line, "n", false)
