---@diagnostic disable: undefined-global

local fs = require("bee.filesystem")
local furi = require("file-uri")

local debug = false

function export.inspect(value, maxDepth, depth, seen)
	maxDepth = maxDepth or 3
	depth = depth or 0
	seen = seen or {}

	local kind = type(value)

	if kind ~= "table" then
		if kind == "string" then
			local maxLength = 150

			if #value > maxLength then
				value = value:sub(1, maxLength) .. "..."
			end

			return string.format("%q", value)
		end

		return tostring(value)
	end

	if seen[value] then
		return "<recursive>"
	end

	if depth >= maxDepth then
		return "<table>"
	end

	seen[value] = true

	local parts = {}

	for key, child in pairs(value) do
		parts[#parts + 1] = {
			key = tostring(key),
			value = child,
		}
	end

	table.sort(parts, function(a, b)
		return a.key < b.key
	end)

	local lines = {
		"{",
	}

	for _, entry in ipairs(parts) do
		lines[#lines + 1] = string.rep("  ", depth + 1)
			.. entry.key
			.. " = "
			.. export.inspect(entry.value, maxDepth, depth + 1, seen)
			.. ","
	end

	lines[#lines + 1] = string.rep("  ", depth) .. "}"

	seen[value] = nil

	return table.concat(lines, "\n")
end

function export.getLocalPath(uri)
	local file_canonical = fs.canonical(furi.decode(uri)):string()
	local doc_canonical = fs.canonical(DOC):string()
	local relativePath = fs.relative(file_canonical, doc_canonical):string()
	if relativePath == "" or relativePath:sub(1, 2) == ".." then
		-- not under project directory
		return "[FOREIGN] " .. file_canonical
	end
	return relativePath
end

function export.positionOf(rowcol)
	return type(rowcol) == "table" and guide.positionOf(rowcol[1], rowcol[2]) or -1
end

function export.sortDoc(a, b)
	if a.name and b.name and a.name ~= b.name then
		return a.name < b.name
	end

	if a.location.file ~= b.location.file then
		return a.location.file < b.location.file
	end

	return export.positionOf(a.location.start) < export.positionOf(b.location.start)
end

--- recursively generate documentation all parser objects downstream of `source`
function export.documentObject(source, has_seen)
	--is this a primative type? then we dont need to process it.
	if type(source) ~= "table" then
		return source
	end

	--set up/check recursion
	if not has_seen then
		has_seen = {}
	end
	if has_seen[source] then
		return nil
	end
	has_seen[source] = true

	--is this an array type? then process each array item and collect it
	if #source > 0 and next(source, #source) == nil then
		local objs = {} --make a pure numerical array
		for i, child in ipairs(source) do
			objs[i] = export.documentObject(child, has_seen)
		end
		return objs
	end

	--if neither, then this is a singular docUnion
	local obj = export.makeDocObject["INIT"](source, has_seen)

	-- print(obj.type)
	local res = export.makeDocObject[obj.type](source, obj, has_seen)

	return res or obj
end

---Switch statement table. functions can be overriden by user file.
---@table
export.makeDocObject = setmetatable({}, {
	__index = function(t, k)
		return function()
			print('DocError: no type "' .. k .. '"')
		end
	end,
})

export.makeDocObject["INIT"] = function(source, has_seen)
	local ok, desc = pcall(getDesc, source)

	return {
		category = source.cate,
		type = source.type,
		name = export.documentObject((source.getCodeName and source:getCodeName()) or source.name, has_seen),
		documentation = ok and desc or "",
		location = {
			file = export.getLocalPath(guide.getUri(source)),
			start = source.start and { guide.rowColOf(source.start) },
			finish = source.finish and { guide.rowColOf(source.finish) },
		},
		-- types = export.documentObject(source.types, has_seen),
	}
end

export.makeDocObject["doc.alias"] = function(source, obj, has_seen)
	obj.file = export.getLocalPath(guide.getUri(source))
end

export.makeDocObject["doc.enum"] = function(source, obj, has_seen)
	obj.file = export.getLocalPath(guide.getUri(source))
end

export.makeDocObject["doc.class"] = function(source, obj, has_seen)
	local extends = source.extends or source.value --doc.class or other
	local field = source.field or source.method
	obj.name = type(field) == "table" and field[1] or nil
	obj.file = export.getLocalPath(guide.getUri(source))
	obj.extends = extends and export.documentObject(extends, has_seen)
	obj.async = vm.isAsync(source, true) and true or false
	obj.deprecated = vm.getDeprecated(source) and true or false
	obj.visible = vm.getVisibleType(source)
end

export.makeDocObject["doc.field.name"] = function(source, obj, has_seen)
	obj["[1]"] = export.documentObject(source[1], has_seen)
	obj.view = source[1]
end

export.makeDocObject["doc.type.arg.name"] = export.makeDocObject["doc.field.name"]

export.makeDocObject["doc.type.function"] = function(source, obj, has_seen)
	obj.args = export.documentObject(source.args, has_seen)
	obj.returns = export.documentObject(source.returns, has_seen)
end

export.makeDocObject["doc.type.table"] = function(source, obj, has_seen)
	obj.fields = export.documentObject(source.fields, has_seen)
end

export.makeDocObject["funcargs"] = function(source, obj, has_seen)
	local objs = {} --make a pure numerical array
	for i, child in ipairs(source) do
		objs[i] = export.documentObject(child, has_seen)
	end
	return objs
end

export.makeDocObject["function"] = function(source, obj, has_seen)
	obj.args = export.documentObject(source.args, has_seen)
	obj.view = getLabel(source, source.parent.type == "setmethod", 1)
	local _, _, max = vm.countReturnsOfFunction(source)
	if max > 0 then
		obj.returns = {}
	end
	for i = 1, max do
		obj.returns[i] = export.documentObject(vm.getReturnOfFunction(source, i), has_seen) --check if bug?
	end
end

export.makeDocObject["function.return"] = function(source, obj, has_seen)
	obj.desc = source.comment and getDesc(source.comment)
	obj.rawdesc = source.comment and getDesc(source.comment, true)
end

export.makeDocObject["local"] = function(source, obj, has_seen)
	obj.name = source[1]
end

export.makeDocObject["self"] = export.makeDocObject["local"]

export.makeDocObject["setglobal"] = export.makeDocObject["doc.class"]

export.makeDocObject["setindex"] = export.makeDocObject["doc.class"]

export.makeDocObject["tableindex"] = function(source, obj, has_seen)
	obj.name = source.index[1]
end

--[[ export.makeDocObject["type"] = function(source, obj, has_seen)
	if export.makeDocObject["variable"](source, obj, has_seen) == false then
		return false
	end
	obj.fields = {}
	vm.getClassFields(ws.rootUri, source, vm.ANY, function(next_source, mark)
		if
			next_source.type == "doc.field"
			or next_source.type == "setfield"
			or next_source.type == "setmethod"
			or next_source.type == "tableindex"
		then
			table.insert(obj.fields, export.documentObject(next_source, has_seen))
		end
	end)
	table.sort(obj.fields, export.sortDoc)
end ]]

local vm = require("vm")
local guide = require("parser.guide")

local function isInMeta(source)
	local uri = guide.getUri(source)
	if not uri or uri == "" then
		return false
	end

	local path = fs.canonical(furi.decode(uri)):string()
	local meta = fs.canonical(METAPATH):string()
	local relative = fs.relative(path, meta):string()

	return relative ~= ".." and relative:sub(1, 3) ~= "../" and relative:sub(1, 3) ~= "..\\"
end

local viewNodeSwitch
viewNodeSwitch = util
	.switch()
	:case("nil")
	:case("boolean")
	:case("string")
	:case("integer")
	:case("number")
	:call(function(source, _infer)
		return {
			kind = "builtin",
			value = source.type,
		}
		-- return source.type
	end)
	:case("global")
	:call(function(source, infer)
		if source.cate ~= "type" then
			return
		end

		--[[ if source.name == "vim.treesitter.highlighter" then
			print(export.inspect(source))
		end ]]

		if isInMeta(source) then
			return {
				kind = "builtin",
				value = source.name,
			}
		else
			return {
				kind = "typereference",
				value = source.name,
			}
		end
	end)
	:case("table")
	:call(function(source, infer, uri)
		local docs = source.bindDocs

		if docs then
			for _, doc in ipairs(docs) do
				if doc.type == "doc.enum" then
					-- TODO: handle enums
					-- return "enum " .. doc.enum[1]
					return {
						kind = "enum",
					}
				end
			end
		end

		-- TODO: what is this
		if #source == 1 and source[1].type == "varargs" then
			local node = vm.getInfer(source[1]):view(uri)
			return ("%s[]"):format(node)
		end

		return {
			kind = "table",
		}
	end)
	:case("function")
	:call(function(source, infer)
		-- TODO: what is this?
		--[[ local parent = source.parent
		if guide.isAssign(parent) then
			infer._hasFunctionDef = true
		end ]]

		local generics = {}

		--Forcing sign calculation
		local sign = vm.getSign(source)

		if sign then
			for _, docGeneric in ipairs(sign.docGeneric) do
				for _, generic in ipairs(docGeneric.generics) do
					local name = generic.generic[1]
					local constraints = {}

					if generic.extends then
						for _, genericTypeConstraint in ipairs(generic.extends.types) do
							local constraintType = export.documentObjectType(genericTypeConstraint)
							table.insert(constraints, constraintType)
						end
					end

					table.insert(generics, { name = name, constraints = constraints })
				end
			end
		end

		local arguments = {}

		for _, arg in ipairs(source.args) do
			local name = arg[1]
			local documentation = arg.bindDocs
					and arg.bindDocs[1]
					and arg.bindDocs[1].comment
					and getDesc(arg.bindDocs[1].comment)
				or ""
			local argType, argInfer = export.documentObjectType(arg)
			local optional = argInfer.node:isOptional()

			table.insert(
				arguments,
				{ name = name, documentation = documentation, infer = argType, optional = optional }
			)
		end

		local returns = {}

		local returnsCount = 0
		for _, returnStatement in ipairs(source.returns or {}) do
			returnsCount = math.max(returnsCount, #returnStatement)
		end

		for i = 1, returnsCount do
			local ret = vm.getReturnOfFunction(source, i)

			assert(ret ~= nil, "ERROR: could not infer return type")

			--TODO: insert unknown when ret is nil
			local name = ret.name or ""
			local documentation = ret.comment and getDesc(ret.comment) or ""
			local retType, retInfer = export.documentObjectType(ret)
			local optional = retInfer.node:isOptional()

			table.insert(returns, { name = name, documentation = documentation, infer = retType, optional = optional })
		end

		return {
			kind = "function",
			generics = generics,
			arguments = arguments,
			returns = returns,
		}
	end)
	--[[ :case("local")
	:call(function(source, infer)
		if source.parent == "funcargs" then
			infer._isParam = true
		else
			infer._isLocal = true
		end
	end) ]]
	:case(
		"doc.type"
	)
	:call(function(source, infer, uri)
		local buf = {}
		for _, tp in ipairs(source.types) do
			buf[#buf + 1] = viewNodeSwitch(tp.type, tp, infer, uri)
		end
		return table.concat(buf, "|")
	end)
	:case("doc.type.name")
	:call(function(source, _infer, uri)
		if source.signs then
			local buf = {}
			for i, sign in ipairs(source.signs) do
				buf[i] = vm.getInfer(sign):view(uri)
			end
			return ("%s<%s>"):format(source[1], table.concat(buf, ", "))
		else
			return source[1]
		end
	end)
	:case("generic")
	:call(function(source, _infer, uri)
		return export.documentObjectType(source.proto)
	end)
	:case("doc.generic.name")
	:call(function(source, _infer, uri)
		return {
			kind = "genericreference",
			value = source[1],
		}
	end)
	:case("doc.type.array")
	:call(function(source, infer, uri)
		return {
			kind = "array",
			items = export.documentObjectType(source.node),
		}
	end)
	:case("doc.type.sign")
	:call(function(source, infer, uri)
		infer._hasClass = true
		local buf = {}
		for i, sign in ipairs(source.signs) do
			local view = vm.getInfer(sign):view(uri)
			-- Strip outer <> from generic names since the sign
			-- already wraps parameters in <>, avoiding list<<T>>
			if view and view:sub(1, 1) == "<" and view:sub(-1) == ">" then
				view = view:sub(2, -2)
			end
			buf[i] = view
		end
		local node = vm.compileNode(source)
		for c in node:eachObject() do
			if guide.isLiteral(c) then
				---@cast c parser.object
				local view = vm.getInfer(c):view(uri)
				if view then
					infer._drop[view] = true
				end
			end
		end
		return ("%s<%s>"):format(source.node[1], table.concat(buf, ", "))
	end)
	:case("doc.type.table")
	--TODO: kind "tuple"
	:call(function(source, infer, uri)
		--[[ if #source.fields == 0 then
			infer._hasTable = true
			return
		end ]]
		-- infer._hasClass = true
		local fields = {}
		for i, field in ipairs(source.fields) do
			local name = source.isTuple and ("[%d]"):format(i)
				or field.name.type == "doc.type" and ("[%s]"):format(vm.getInfer(field.name):view(uri))
				or type(field.name[1] == "string") and field.name[1]
				or ("[%q]"):format(field.name[1])
			local documentation = ""
			local fieldType, fieldInfer = export.documentObjectType(field)
			local optional = fieldInfer.node:isOptional()

			table.insert(fields, { name = name, documentation = documentation, infer = fieldType, optional = optional })
		end

		return {
			kind = "table",
			fields = fields,
		}
		--[[ local buf = {}
		buf[#buf + 1] = source.isTuple and "[" or "{ "
		for i, field in ipairs(source.fields) do
			if i > 1 then
				buf[#buf + 1] = ", "
			end
			if not source.isTuple then
				local key = field.name
				if key.type == "doc.type" then
					buf[#buf + 1] = ("[%s]: "):format(vm.getInfer(key):view(uri))
				elseif type(key[1]) == "string" then
					buf[#buf + 1] = key[1] .. ": "
				else
					buf[#buf + 1] = ("[%q]: "):format(key[1])
				end
			end
			buf[#buf + 1] = vm.getInfer(field.extends):view(uri)
		end
		buf[#buf + 1] = source.isTuple and "]" or " }"
		return table.concat(buf) ]]
	end)
	:case("doc.type.string")
	:call(function(source, _infer)
		return util.viewString(source[1], source[2])
	end)
	:case("doc.type.integer")
	:case("doc.type.boolean")
	:call(function(source, _infer)
		return ("%q"):format(source[1])
	end)
	:case("doc.type.code")
	:call(function(source, _infer)
		return ("`%s`"):format(source[1])
	end)
	:case("doc.type.function")
	:call(function(source, infer, uri)
		local generics = {}

		--Forcing sign calculation
		local sign = vm.getSign(source)

		--TODO: does this work?
		--This was copied as is from :case("function")
		--When that was done, this had no effect on the output
		if sign then
			for _, docGeneric in ipairs(sign.docGeneric) do
				for _, generic in ipairs(docGeneric.generics) do
					local name = generic.generic[1]
					local constraints = {}

					if generic.extends then
						for _, genericTypeConstraint in ipairs(generic.extends.types) do
							local constraintType = export.documentObjectType(genericTypeConstraint)
							table.insert(constraints, constraintType)
						end
					end

					table.insert(generics, { name = name, constraints = constraints })
				end
			end
		end

		local arguments = {}

		for _, arg in ipairs(source.args) do
			local name = arg.name[1]
			local argType, argInfer = export.documentObjectType(arg)
			local optional = argInfer.node:isOptional()

			table.insert(arguments, { name = name, infer = argType, optional = optional })
		end

		local returns = {}

		for _, ret in ipairs(source.returns) do
			local name = ret.name and ret.name[1] or ""
			local retType, retInfer = export.documentObjectType(ret)
			local optional = retInfer.node:isOptional()

			table.insert(returns, { name = name, infer = retType, optional = optional })
		end

		return {
			kind = "function",
			generics = generics,
			arguments = arguments,
			returns = returns,
		}
	end)
	:case("doc.field.name")
	:call(function(source, _infer, uri)
		return vm.viewKey(source, uri)
	end)

--TODO: Finish implementing infer:view logics
function export.documentObjectType(source)
	local infer = vm.getInfer(source)

	local view = infer:view(ws.rootUri)
	local types = {}

	print("view", export.inspect(infer.views))

	for nodeObject in infer.node:eachObject() do
		if nodeObject.hideView then
			goto continue
		end

		-- Avoiding to avoid expand aliases
		-- infer._drop contains the string representations of expanded types that should be suppressed
		local nodeObjectView = vm.viewObject(nodeObject, ws.rootUri)

		if not nodeObjectView or infer._drop[nodeObjectView] or not infer.views[nodeObjectView] then
			goto continue
		end

		print("nView", nodeObjectView)

		local nodeObjectType = viewNodeSwitch(nodeObject.type, nodeObject, infer, ws.rootUri)

		if not nodeObjectType then
			goto continue
		end

		table.insert(types, nodeObjectType)

		::continue::
	end

	return {
		view = view,
		types = types,
	}, infer
end

function export.followTypes(source, sourceType, documentObjectField)
	documentObjectField = documentObjectField or export.documentObject

	for i, t in ipairs(sourceType.types) do
		if (t.kind == "table" and not t.fields) or (t.kind == "builtin" and t.value == "table") then
			local fields = export.documentObjectFields(source, documentObjectField)

			if #fields > 0 then
				sourceType.types[i] = {
					kind = "table",
					fields = fields,
				}
			else
				sourceType.types[i] = {
					kind = "builtin",
					value = "table",
				}
			end
		end
	end

	return sourceType
end

--TODO: merge fields with the same name (look at vim.v)
--TODO: optional fields
function export.documentObjectFields(source, documentObjectField)
	local sources = {}
	local has_seen = {}

	local function addFields(fields)
		for _, fieldObj in ipairs(fields or {}) do
			if not has_seen[fieldObj] then
				has_seen[fieldObj] = true
				sources[#sources + 1] = fieldObj
			end
		end
	end

	addFields(vm.getFields(source))

	local docs = {}

	for _, fieldObj in ipairs(sources) do
		-- if not isSourceBuiltin(fieldObj) then
		local fieldDoc = documentObjectField(fieldObj)

		if fieldDoc then
			docs[#docs + 1] = fieldDoc
		end
		-- end
	end

	table.sort(docs, export.sortDoc)

	return docs
end

--[[ export.makeDocObject["doc.type"] = function(source, init, has_seen)
	local doc = init
	print(export.inspect(init))

	doc.infer = export.documentObjectType(source)

	return init
end ]]

export.makeDocObject["variable"] = function(source, init, has_seen)
	-- local definitions = source:getSets(ws.rootUri)
	local definitions = vm.getDefs(source)

	assert(
		#definitions == 1,
		string.format("Global '%s' with an unexpected number of definitions received", source.name)
	)

	local doc = util.tableMerge(init, export.documentObject(definitions[1], has_seen))

	local definition = definitions[1]

	doc.meta = export.documentObjectMeta(definition)

	local definitionType = export.documentObjectType(definition)
	doc.infer = export.followTypes(definition, definitionType)

	return doc
end

-- TODO: global types
export.makeDocObject["global"] = function(source, init, has_seen)
	assert(source.cate == "variable", "ERROR: global types yet to be implemented")

	return export.makeDocObject[source.cate](source, init, has_seen)

	--[[ for _, set in ipairs(definitions) do
		if
			set.type == "setglobal"
			or set.type == "setfield"
			or set.type == "setmethod"
			or set.type == "setindex"
			or set.type == "doc.alias"
			or set.type == "doc.enum"
			or set.type == "doc.class"
		then
			table.insert(obj.defines, export.documentObject(set, has_seen))
		end
	end
	if #obj.defines == 0 then
		return false
	end
	table.sort(obj.defines, export.sortDoc) ]]
end

export.documentObjectMeta = function(source, meta)
	meta = meta or {}

	table.insert(meta, vm.getVisibleType(source))

	if vm.isAsync(source, true) then
		table.insert(meta, "async")
	end

	if vm.getDeprecated(source) then
		table.insert(meta, "deprecated")
	end

	return meta
end

export.makeDocObject["setfield"] = function(source, init, has_seen)
	assert(source.field ~= nil, "Setfield source without field property")
	assert(source.value ~= nil, "Setfield source without value property")

	-- local extends = source.extends or source.value --doc.class or other
	-- init.file = export.getLocalPath(guide.getUri(source))
	-- init.extends = extends and export.documentObject(extends, has_seen)

	local doc = init

	local field = source.field

	--TODO: handle metamethods
	if field[1] == "__index" then
		return nil
	end

	debug = field[1] == "_enabled_configs"

	if debug then
		-- print(export.inspect(source))
	end

	-- TODO: why does this do
	doc.name = field[1]
	doc.meta = export.documentObjectMeta(source)

	local candidates = { source.field, source.value }

	for _, src in ipairs(candidates) do
		local srcType = export.documentObjectType(src)
		local srcInfer = export.followTypes(src, srcType)

		if #srcInfer.types > 0 then
			doc.infer = srcInfer
			break
		end
	end

	return doc
end

export.makeDocObject["setmethod"] = function(source, init, has_seen)
	assert(source.method ~= nil, "ERROR: Setmethod source without method property")

	local doc = init
	local method = source.method

	doc.name = method[1]
	doc.meta = export.documentObjectMeta(source, { "method" })

	local methodType = export.documentObjectType(method)
	doc.infer = export.followTypes(method, methodType)

	return doc
end

export.makeDocObject["tablefield"] = function(source, init, has_seen)
	local doc = init

	assert(source.field ~= nil, "Tablefield source without field property")
	assert(source.value ~= nil, "Tablefield source without value property")

	local field = source.field

	doc.name = field[1]
	doc.meta = export.documentObjectMeta(source)

	local candidates = { source.field, source.value }

	for _, src in ipairs(candidates) do
		local srcInfer = export.documentObjectType(src)
		local srcType = export.followTypes(src, srcInfer)

		if #srcType.types > 0 then
			doc.infer = srcType
			break
		end
	end

	return doc
end

export.makeDocObject["doc.field"] = function(source, init, has_seen)
	assert(source.field ~= nil, "ERROR: doc.field source missing field")
	assert(source.extends ~= nil, "ERROR: doc.field source missing extends")

	local doc = init

	doc.meta = export.documentObjectMeta(source)

	local field = source.field
	local extends = source.extends

	if field.type == "doc.field.name" then
		doc.name = field[1]
	else
		doc.name = ("[%s]"):format(vm.getInfer(field):view(ws.rootUri))
	end

	local extendsType = export.documentObjectType(extends)
	doc.infer = export.followTypes(extends, extendsType)

	--[[ doc.file = export.getLocalPath(guide.getUri(source))
	doc.extends = source.extends and export.documentObject(source.extends, has_seen) --check if bug?
	doc.async = vm.isAsync(source, true) and true or false --if vm.isAsync(set, true) then result.defines[#result.defines].extends['async'] = true end
	doc.deprecated = vm.getDeprecated(source) and true or false --  if (depr and not depr.versions) the result.defines[#result.defines].extends['deprecated'] = true end
	doc.visible = vm.getVisibleType(source) ]]

	return doc
end

---gathers the globals that are to be exported in documentation
---@async
---@return table globals
function export.gatherGlobals()
	-- Adding a new line after "Loading documents ..." in cli output
	print("")

	local globals = util.valuesOf(vm.getExportableGlobals())
	local exportedGlobals = {}

	for _, source in ipairs(globals) do
		-- TODO: env var
		if source.name == "vim" then
			table.insert(exportedGlobals, source)
		end
	end

	return exportedGlobals
end

---builds a lua table of based on `globals` and their elements
---@async
---@param globals table
---@param callback fun(i, max)
function export.makeDocs(globals, callback)
	local docs = {}
	for i, globalVar in ipairs(globals) do
		table.insert(docs, export.documentObject(globalVar))
		callback(i, #globals)
	end
	table.sort(docs, export.sortDoc)
	return docs
end

function export.getLualsConfig()
	return nil
end

---takes the table from `makeDocs`, serializes it, and exports it
export.serializeAndExport = function(docs, outputDir)
	local old_jsonb_supportSparseArray = jsonb.supportSparseArray
	jsonb.supportSparseArray = true

	local exportPath = outputDir
	local variablesPath = outputDir
	local typesPath = outputDir .. "/type"

	--export all docs  to json

	local jsonPath = exportPath .. "/doc.json"
	local ok, err = util.saveFile(jsonPath, jsonb.beautify(docs))

	if not ok then
		return false, { jsonPath }, { err }
	end

	--export type docs to json

	local typeDocs = {}

	for i = 1, #docs do
		local doc = docs[i]
		if doc.category == "type" then
			table.insert(typeDocs, doc)
		end
	end

	for i = 1, #typeDocs do
		local typeDoc = typeDocs[i]
		local typePath = typesPath .. "/" .. typeDoc.name .. ".json"

		local typeOk, typeErr = util.saveFile(typePath, jsonb.beautify(typeDoc))

		if not typeOk then
			return false, { typePath }, { typeErr }
		end
	end

	--export variable docs to json

	-- local variableDocs = {}

	local function getJsonDocs(parent, prefix)
		local parentName = prefix or parent.name

		local jsons = { { name = parentName, doc = parent } }

		-- TODO: remove
		if not parent.infer then
			return jsons
		end

		if #parent.infer.types ~= 1 or parent.infer.types[1].kind ~= "table" or not parent.infer.types[1].fields then
			return jsons
		end

		for i = #parent.infer.types[1].fields, 1, -1 do
			local field = parent.infer.types[1].fields[i]

			-- TODO: remove
			if not field.infer then
				goto continue
			end

			if #field.infer.types ~= 1 or field.infer.types[1].kind ~= "table" or not field.infer.types[1].fields then
				goto continue
			end

			local childName = string.format("%s.%s", parentName, field.name)
			local child = table.remove(parent.infer.types[1].fields, i)

			for _, childChild in ipairs(getJsonDocs(child, childName)) do
				table.insert(jsons, childChild)
			end

			::continue::
		end

		return jsons
	end

	local variableJsons = {}
	for i = 1, #docs do
		local doc = docs[i]
		if doc.category == "variable" then
			for _, jsonDoc in ipairs(getJsonDocs(doc)) do
				table.insert(variableJsons, jsonDoc)
			end
		end
	end

	for _, variableJson in ipairs(variableJsons) do
		local variableJsonPath = variablesPath .. "/" .. variableJson.name .. ".json"

		local variableOk, variableJsonErr = util.saveFile(variableJsonPath, jsonb.beautify(variableJson.doc))

		if not variableOk then
			return false, { variableJsonPath }, { variableJsonErr }
		end
	end

	--[[ for i = 1, #variableDocs do
		local variableDoc = variableDocs[i]
		local variablePath = variablesPath .. "/" .. variableDoc.name .. ".json"

		local variableOk, variableErr = util.saveFile(variablePath, jsonb.beautify(variableDoc))

		if not variableOk then
			return false, { variablePath }, { variableErr }
		end
	end ]]

	jsonb.supportSparseArray = old_jsonb_supportSparseArray

	return true, { jsonPath }
end

return export
