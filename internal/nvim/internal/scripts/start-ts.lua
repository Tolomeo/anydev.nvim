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

	vim.g.ts_ready = true
end

vim.treesitter.start(bufnr, luadoc)
vim.treesitter.start(bufnr, lua)

