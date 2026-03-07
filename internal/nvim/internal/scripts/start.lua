---@class anydev_ts_parser
---@field source string
---@field parser vim.treesitter.LanguageTree

---@class anydev_lsp_activity
---@field last integer
---@field activity table<string, any>

---@class anydev
---@field ts_parsers table<string, anydev_ts_parser>
---@field lsp_activity anydev_lsp_activity
_G.Anydev = {
	ts_parsers = {},
	lsp_activity = { last = vim.loop.now() + 500, activity = {} },
}

_G.Anydev.wait_for_lsp_idle = (function()
	vim.api.nvim_create_autocmd("LspProgress", {
		group = vim.api.nvim_create_augroup("LuaLSReady", { clear = true }),
		callback = function(autocmd_args)
			local message = string.format(
				"[Anydev:LspProgress:%s]:%s",
				autocmd_args.data.params.value.kind,
				vim.api.nvim_buf_get_name(0)
			)
			vim.cmd(string.format("echom '%s'", message))

			local value = autocmd_args.data.params.value
			local token = autocmd_args.data.params.token

			if value.kind == "begin" then
				_G.Anydev.lsp_activity.activity[token] = value
			elseif value.kind == "end" then
				_G.Anydev.lsp_activity.activity[token] = nil
			end

			_G.Anydev.lsp_activity.last = vim.loop.now()
		end,
	})

	vim.lsp.enable("lua_ls")

	return function()
		vim.wait(5000, function()
			local has_activity = next(_G.Anydev.lsp_activity.activity) ~= nil

			if has_activity then
				return false
			end

			local is_settled = (vim.loop.now() - _G.Anydev.lsp_activity.last) > 1500

			return is_settled
		end, 100)
	end
end)()

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

---@class anydev_ts_parser_node
---@field type string
---@field range { start: { line: number, character: number }, end: { line: number, character: number }}
---@field text string

---@param node TSNode
---@param node_source string
---@return anydev_ts_parser_node
function _G.Anydev:get_ts_parser_node(node, node_source)
	local nodeType = node:type()
	local startLine, startCharacter, endLine, endCharacter = node:range()
	local text = vim.treesitter.get_node_text(node, node_source)

	local anydev_ts_parser_node = {
		type = nodeType,
		range = {
			start = { line = startLine, character = startCharacter },
			["end"] = { line = endLine, character = endCharacter },
		},
		text = text,
	}

	return anydev_ts_parser_node
end

---@param str string
---@return boolean
---@return string
function _G.Anydev:is_comment(str)
	local rest = str:match("^[%,%;%s]*(.*)")
	return rest:sub(1, 2) == "--", rest
end

---@param char string
function _G.Anydev:is_separator(char)
	return char == "," or char == ";"
end
