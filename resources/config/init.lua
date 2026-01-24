vim.opt.swapfile = false
vim.opt.backup = false
vim.opt.writebackup = false
vim.opt.undofile = false
vim.g._ts_force_sync_parsing = true

local config_dir = vim.fn.expand("<sfile>:p:h")

vim.opt.runtimepath:prepend(config_dir)
vim.opt.packpath:prepend(config_dir)
-- vim.opt.packpath = { config_dir }
-- package.path = package.path .. ';' .. config_dir .. '/lua/?.lua'

-- Configuring lua_ls lsp
require("lazydev").setup()

---@class anydev_ts_parser
---@field source string
---@field parser vim.treesitter.LanguageTree

---@class anydev
---@field ts_parsers table<string, anydev_ts_parser>
_G.Anydev = {
	ts_parsers = {},
}

---@param buffer integer
---@return anydev_ts_parser
function _G.Anydev:get_ts_parser(buffer)
	local buffer_name = vim.api.nvim_buf_get_name(buffer)

	if self.ts_parsers[buffer_name] ~= nil then
		return self.ts_parsers[buffer_name]
	end

	local source = table.concat(vim.api.nvim_buf_get_lines(buffer, 0, -1, false), "\n")
	local parser = vim.treesitter.get_string_parser(source, "lua")
	parser:parse(true)

	local anydev_ts_parser = { source = source, parser = parser }
	self.ts_parsers[buffer_name] = anydev_ts_parser

	local group = vim.api.nvim_create_augroup("TsParserDelete", { clear = true })
	vim.api.nvim_create_autocmd({ "TextChanged", "TextChangedI", "BufDelete" }, {
		group = group,
		buffer = buffer,
		callback = function()
			_G.Anydev.ts_parsers[buffer_name] = nil
		end,
	})

	return anydev_ts_parser
end
