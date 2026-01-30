local args = { ... }
local line = args[1]
local character = args[2]
local timeout = args[3]

local bufnr = 0
local position = { line = line, character = character }
local textDocument = vim.lsp.util.make_text_document_params(bufnr)

local lspResponse, err =
	vim.lsp.buf_request_sync(bufnr, "textDocument/hover", { textDocument = textDocument, position = position }, timeout)

if err ~= nil then
	error(err)
end

if lspResponse == nil then
	return vim.NIL
end

local _, lspClientResponse = next(lspResponse)

if lspClientResponse == nil or next(lspClientResponse) == nil then
	return vim.NIL
end

return vim.fn.json_encode(lspClientResponse)
