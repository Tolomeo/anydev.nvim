---@type vim.lsp.Config
---@see https://raw.githubusercontent.com/LuaLS/vscode-lua/master/setting/schema.json
return {
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
				path = {
					"lua/?.lua",
					"lua/?/init.lua",
				},
			},
			workspace = {
				checkThirdParty = true,
				library = {
					vim.env.VIMRUNTIME,
					"${3rd}/luv/library",
					"${3rd}/busted/library",
				},
				maxPreload = 999999,
				-- https://github.com/neovim/nvim-lspconfig/issues/3189#issuecomment-3021345989
				-- TODO: config path from current dir
				--[[ library = vim.tbl_filter(function(d)
					return not d:find(vim.fn.stdpath("config"), 1, true)
				end, vim.api.nvim_get_runtime_file("", true)), ]]
			},
			completion = {
				maxSuggestCount = 999999
			},
			diagnostics = {
				enable = false
			}
		},
	},
}
