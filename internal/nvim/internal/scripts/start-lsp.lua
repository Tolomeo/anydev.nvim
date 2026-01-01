if vim.g.lua_ls_ready == true then
	return
end

local args = { ... }
local delay = args[1]

vim.g.lua_ls_ready = false

local lsp_inflight_events = {}

vim.api.nvim_create_augroup("LuaLSReady", { clear = true })

vim.api.nvim_create_autocmd("LspProgress", {
	group = "LuaLSReady",
	callback = function(autocmd_args)
		local message = string.format("[anydev:LspProgress]:%s", autocmd_args.data.params.value.kind)
		vim.cmd(string.format("echom '%s'", message))

		local value = autocmd_args.data.params.value
		local token = autocmd_args.data.params.token

		if value.kind == "begin" then
			lsp_inflight_events[token] = value
		elseif value.kind == "end" then
			lsp_inflight_events[token] = nil
		end

		vim.g.lua_ls_ready = next(lsp_inflight_events) == nil
	end,
})

vim.lsp.enable("lua_ls")

vim.wait(delay, function()
	return vim.g.lua_ls_ready == true
end)
