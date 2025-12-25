if vim.g.lua_ts_ready == true then
	return
end

vim.g.lua_ts_ready = false

local ok = pcall(vim.treesitter.language.add, "lua")

if not ok then
	error("Treesitter Lua parser registration failed.")
end

ok = pcall(vim.treesitter.language.add, "luadoc")

if not ok then
	error("Treesitter Luadoc parser registration failed.")
end

vim.api.nvim_create_autocmd("FileType", {
	pattern = { "lua" },
	callback = function(opts)
		vim.treesitter.start(opts.buf, "luadoc")
		vim.treesitter.start(opts.buf, "lua")
	end,
})

vim.treesitter.start(0, "luadoc")
vim.treesitter.start(0, "lua")

vim.g.lua_ts_ready = true
