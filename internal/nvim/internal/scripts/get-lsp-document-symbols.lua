local textDocumentParams = vim.lsp.util.make_text_document_params(0)
local result = vim.lsp.buf_request_sync(0, "textDocument/documentSymbol", { textDocument = textDocumentParams }, 2000)

if result == nil then
	return vim.NIL
end

return vim.fn.json_encode(result[1])
