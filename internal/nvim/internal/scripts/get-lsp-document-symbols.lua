local bufnr = 0

local textDocumentParams = vim.lsp.util.make_text_document_params(bufnr)
local lspResponse, err =
	vim.lsp.buf_request_sync(bufnr, "textDocument/documentSymbol", { textDocument = textDocumentParams }, 2000)

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
