local args = {...}
local delay = args[1]

local ok = pcall(vim.treesitter.language.add, "lua")

if not ok then
	error("Treesitter Lua parser registration failed.")
end

ok = pcall(vim.treesitter.language.add, "luadoc")

if not ok then
	error("Treesitter Luadoc parser registration failed.")
end

local bufnr = 0
local lua = "lua"
local luadoc = "luadoc"

vim.api.nvim_create_autocmd("FileType", {
	group = vim.api.nvim_create_augroup("LuaTSReady", { clear = true }),
	pattern = { lua },
	callback = function(opts)
		vim.treesitter.start(opts.buf, luadoc)
		vim.treesitter.start(opts.buf, lua)
	end,
})

vim.treesitter.start(bufnr, luadoc)
vim.treesitter.start(bufnr, lua)

-- Query injections are only recalculated when the buffer changes
-- so we make a (hopefully) inhert change to force their presence
-- by adding an empty line at the end of the buffer text
local keys = vim.api.nvim_replace_termcodes("Go<Esc>", true, false, true)
vim.api.nvim_feedkeys(keys, "n", false)
-- Forcing treesitter to parse all trees
vim.treesitter.get_parser(bufnr, lua):parse(true)
