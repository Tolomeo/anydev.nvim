-- Adding the folder to the runtimepath & packpath

local config_dir = vim.fn.expand("<sfile>:p:h")

vim.opt.runtimepath:prepend(config_dir)
vim.opt.packpath:prepend(config_dir)
-- package.path = package.path .. ';' .. config_dir .. '/lua/?.lua'

-- Configuring lua_ls lsp
require("lazydev").setup()

vim.lsp.config["lua_ls"] = {
	-- Command and arguments to start the server.
	cmd = { "lua-language-server" },
	-- Filetypes to automatically attach to.
	filetypes = { "lua" },
	-- Sets the "workspace" to the directory where any of these files is found.
	-- Files that share a root directory will reuse the LSP server connection.
	-- Nested lists indicate equal priority, see |vim.lsp.Config|.
	root_markers = { ".root" },
	-- Specific settings to send to the server. The schema is server-defined.
	-- Example: https://raw.githubusercontent.com/LuaLS/vscode-lua/master/setting/schema.json
	settings = {
		Lua = {
			runtime = {
				version = "LuaJIT",
			},
		},
	},
}

vim.lsp.enable("lua_ls")
-- Wait for lsp to attach

--vim.wait(2000, function()
--	return next(vim.lsp.get_clients()) ~= nil
--end)
--
