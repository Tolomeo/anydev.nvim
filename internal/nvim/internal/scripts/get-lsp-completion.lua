local args = { ... }
local line = args[1]
local character = args[2]
local delay = args[3]
local textDocumentParams = vim.lsp.util.make_text_document_params(0)
local positionParams = { line = line, character = character }
local lspResponse, err = vim.lsp.buf_request_sync(
	0,
	"textDocument/completion",
	{ textDocument = textDocumentParams, position = positionParams },
	delay
)

if err ~= nil then
	error(err)
end

if lspResponse == nil or next(lspResponse) == nil then
	return vim.NIL
end

return vim.fn.json_encode(lspResponse[1])
