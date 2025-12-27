vim.opt.swapfile = false
vim.opt.backup = false
vim.opt.writebackup = false
vim.opt.undofile = false

local config_dir = vim.fn.expand("<sfile>:p:h")

vim.opt.runtimepath:prepend(config_dir)
vim.opt.packpath:prepend(config_dir)
-- vim.opt.packpath = { config_dir }
-- package.path = package.path .. ';' .. config_dir .. '/lua/?.lua'

-- Configuring lua_ls lsp
require("lazydev").setup()
