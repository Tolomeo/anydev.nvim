-- Adding the folder to the runtimepath & packpath

local config_dir = vim.fn.expand('<sfile>:p:h')

vim.opt.runtimepath:prepend(config_dir)
vim.opt.packpath:prepend(config_dir)
-- package.path = package.path .. ';' .. config_dir .. '/lua/?.lua'

-- Configuring lua_ls lsp
require('lazydev').setup()
vim.lsp.enable('lua_ls')
