local args = { ... }
local line = args[1]
local character = args[2]
local position = { line = line, character = character }
local textDocument = vim.lsp.util.make_text_document_params(0)
local timeout = args[3]

local result, err =
	vim.lsp.buf_request_sync(0, "textDocument/hover", { textDocument = textDocument, position = position }, timeout)

if err ~= nil then
	error(err)
end

if result == nil then
	return vim.NIL
end

return vim.fn.json_encode(result[1])
