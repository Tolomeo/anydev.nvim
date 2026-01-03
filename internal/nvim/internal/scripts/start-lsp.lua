local args = { ... }
local delay = args[1]
local settle_time = 500

vim.g.lsp_activity = vim.g.lsp_activity and vim.g.lsp_activity or { last = vim.loop.now(), activity = {} }

vim.api.nvim_create_autocmd("LspProgress", {
	group = vim.api.nvim_create_augroup("LuaLSReady", { clear = true }),
	callback = function(autocmd_args)
		local value = autocmd_args.data.params.value
		local token = autocmd_args.data.params.token

		if value.kind == "begin" then
			vim.g.lsp_activity.activity[token] = value
		elseif value.kind == "end" then
			vim.g.lsp_activity.activity[token] = nil
		end

		vim.g.lsp_activity.last = vim.loop.now()
	end,
})

vim.lsp.enable("lua_ls")

vim.wait(delay, function()
	local has_activity = next(vim.g.lsp_activity.activity) ~= nil

	if has_activity then
		return false
	end

	local is_settled = (vim.loop.now() - vim.g.lsp_activity.last) > settle_time

	return is_settled
end, 100)
