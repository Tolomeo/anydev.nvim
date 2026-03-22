local args = { ... }
local line = args[1]
local character = args[2]
local delay = args[3]
local bufnr = 0
local textDocumentParams = vim.lsp.util.make_text_document_params(bufnr)
local positionParams = { line = line, character = character }

_G.Anydev:wait_for_lsp_idle()

local lspResponse, err = vim.lsp.buf_request_sync(
	bufnr,
	"textDocument/definition",
	{ textDocument = textDocumentParams, position = positionParams },
	delay
)

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
