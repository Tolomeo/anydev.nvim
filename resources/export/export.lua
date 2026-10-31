---@diagnostic disable: undefined-global

local fs = require("bee.filesystem")
local furi = require("file-uri")

-- ---------------------------------------------------------------------------
-- Available APIs
-- ---------------------------------------------------------------------------
--
-- Keeping these aliases here makes it easy to remember what LuaLS exposes to
-- custom exporters and gives LuaLS/type tooling something explicit to inspect.

local api = {
	-- Built-in exporter clone.
	export = export,

	-- Workspace.
	ws = ws,

	-- LuaLS VM/type-analysis API.
	vm = vm,

	-- Parser/navigation helpers.
	guide = guide,

	-- Hover/documentation helpers.
	getDesc = getDesc,
	getLabel = getLabel,

	-- JSON serializer used by the default exporter.
	jsonb = jsonb,

	-- General LuaLS utilities.
	util = util,

	-- Markdown builder.
	markdown = markdown,

	-- Internal modules available through require().
	fs = fs,
	furi = furi,
}

local anydev = {}

local function inspect(value, maxDepth, depth, seen)
	maxDepth = maxDepth or 3
	depth = depth or 0
	seen = seen or {}

	local kind = type(value)

	if kind ~= "table" then
		if kind == "string" then
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
			.. inspect(entry.value, maxDepth, depth + 1, seen)
			.. ","
	end

	lines[#lines + 1] = string.rep("  ", depth) .. "}"

	seen[value] = nil

	return table.concat(lines, "\n")
end

--[[ local function getDeclaredType(source)
	local docs = source.bindDocs

	if not docs then
		return nil
	end

	-- LuaLS generally treats later bound docs as higher priority,
	-- so scan backwards.
	for i = #docs, 1, -1 do
		local doc = docs[i]

		if doc.type == "doc.type" then
			return vm.getInfer(doc):view(ws.rootUri)
		end
	end

	return nil
end ]]

local function split(str, separator)
	local result = {}
	local start = 1

	while true do
		local pos = str:find(separator, start, true)
		if not pos then
			result[#result + 1] = str:sub(start)
			break
		end

		result[#result + 1] = str:sub(start, pos - 1)
		start = pos + #separator
	end

	return result
end

local function builtinTypeReference(source)
	return {
		name = source:getCodeName(),
		type = source.cate,
		view = vm.getInfer(source):view(ws.rootUri),
	}
end

local function isSourceBuiltin(source)
	return type(source) == "table"
		and source.type == "global"
		and source.cate == "type"
		and source.uri ~= nil
		and source.uri:find(METAPATH, 1, true) ~= nil
end

local function getSourceFields(source, has_seen)
	local fieldObjs = {}
	local collected = {}

	local function addFields(fields)
		for _, fieldObj in ipairs(fields or {}) do
			if not collected[fieldObj] then
				collected[fieldObj] = true
				fieldObjs[#fieldObjs + 1] = fieldObj
			end
		end
	end

	if source.type == "doc.class" then
		-- Explicit @field annotations.
		-- addFields(source.fields)

		-- All fields from the class definition, including inherited ones.
		addFields(vm.getFields(source))

		-- Fields from the declaration bound to the class.
		local bindSource = source.bindSource

		if bindSource then
			if bindSource.value then
				addFields(vm.getFields(bindSource.value))
			end

			addFields(vm.getFields(bindSource))
		end
	elseif source.value then
		-- Existing behaviour for other definitions.
		addFields(vm.getFields(source.value))
	end

	local fieldDocs = {}

	for _, fieldObj in ipairs(fieldObjs) do
		if not isSourceBuiltin(fieldObj) then
			local fieldDoc = export.documentObject(fieldObj, has_seen)

			if fieldDoc then
				fieldDocs[#fieldDocs + 1] = fieldDoc
			end
		end
	end

	if #fieldDocs == 0 then
		return false, {}
	end

	table.sort(fieldDocs, export.sortDoc)

	return true, fieldDocs
end

local function getSourceDefines(source, has_seen)
	local defines = {}

	for _, set in ipairs(source:getSets(ws.rootUri)) do
		if
			set.type == "setglobal"
			or set.type == "setfield"
			or set.type == "setmethod"
			or set.type == "setindex"
			or set.type == "doc.alias"
			or set.type == "doc.enum"
			or set.type == "doc.class"
			or set.type == "doc.field"
		then
			local define = export.documentObject(set, has_seen)

			if define then
				local hasFields, fields = getSourceFields(set, has_seen)

				if hasFields then
					define.fields = fields
				end

				table.insert(defines, define)
			end
		end
	end

	if #defines == 0 then
		return false, defines
	end

	table.sort(defines, export.sortDoc)

	return true, defines
end

local function getInferredType(source, has_seen)
	local infer = vm.getInfer(source)
	local node = infer.node

	local types = {}

	for typeSource in node:eachObject() do
		local typeDoc

		if isSourceBuiltin(typeSource) then
			typeDoc = builtinTypeReference(typeSource)
		else
			typeDoc = export.documentObject(typeSource, has_seen)
		end

		if typeDoc then
			types[#types + 1] = typeDoc
		end
	end

	return {
		optional = node:isOptional(),
		types = types,
	}
end

export.makeDocObject["funcargs"] = function(source, obj, has_seen)
	local objs = {} --make a pure numerical array
	for i, child in ipairs(source) do
		objs[i] = export.documentObject(child, has_seen)
		if objs[i] then
			objs[i].t = getInferredType(child, has_seen)
		end
	end

	return objs
end

export.makeDocObject["doc.alias"] = function(source, obj, has_seen)
	obj.file = export.getLocalPath(guide.getUri(source))
	obj.extends = source.extends and export.documentObject(source.extends, has_seen)
end

local originalSetfield = export.makeDocObject["setfield"]

export.makeDocObject["setfield"] = function(source, obj, has_seen)
	local result = originalSetfield(source, obj, has_seen)
	--[[ local hasFields, fields = getSourceFields(source, has_seen)

	if hasFields then
		obj.fields = fields
	end ]]

	return result
end

-- local originalVariable = export.makeDocObject["variable"]
--

export.makeDocObject["variable"] = function(source, obj, has_seen)
	local hasDefines, defines = getSourceDefines(source, has_seen)

	if not hasDefines then
		return false
	end

	obj.defines = defines

	--[[ local hasFields, fields = getSourceFields(source, has_seen)

	if hasFields then
		obj.fields = fields
	end ]]

	obj.anydev = {
		category = "variable",
	}

	return result
end

-- local originalType = export.makeDocObject["type"]

export.makeDocObject["type"] = function(source, obj, has_seen)
	local hasDefines, defines = getSourceDefines(source, has_seen)

	if not hasDefines then
		return false
	end

	obj.defines = defines

	--[[ local fields = {}

	vm.getClassFields(ws.rootUri, source, vm.ANY, function(next_source, mark)
		if
			next_source.type == "doc.field"
			or next_source.type == "setfield"
			or next_source.type == "setmethod"
			or next_source.type == "tableindex"
		then
			table.insert(fields, export.documentObject(next_source, has_seen))
		end
	end)

	if #fields > 0 then
		table.sort(fields, export.sortDoc)
		obj.fields = fields
	end ]]

	obj.anydev = {
		category = "type",
	}

	return result
end

export.makeDocObject["doc.type.field"] = function(source, obj, has_seen)
	if source.name then
		if source.name.type == "doc.field.name" then
			-- Named field:
			-- { foo: string }
			obj.name = source.name[1]
		else
			-- Typed/index field:
			-- { [string]: number }
			obj.name = ("[%s]"):format(vm.getInfer(source.name):view(ws.rootUri))
		end
	end

	if source.extends then
		obj.extends = export.documentObject(source.extends, has_seen)
	end

	if source.optional then
		obj.optional = true
	end
end

-- ---------------------------------------------------------------------------
-- Built-in exporter functions
-- ---------------------------------------------------------------------------
--
-- The current built-in exporter provides:
--
--     export.getLocalPath(uri)
--     export.positionOf(rowcol)
--     export.sortDoc(a, b)
--     export.documentObject(source, has_seen)
--     export.makeDocObject
--     export.gatherGlobals()
--     export.makeDocs(globals, callback)
--     export.getLualsConfig()
--     export.serializeAndExport(docs, outputDir)
--
-- Usually you only need to wrap/replace one or two of these.

-- ---------------------------------------------------------------------------
-- Example: customize individual AST/document object types
-- ---------------------------------------------------------------------------

local originalInit = export.makeDocObject.INIT

export.makeDocObject.INIT = function(source, has_seen)
	-- `source` is the original LuaLS parser/VM object.
	--
	-- Examples of APIs that are useful at this level:
	--
	-- local uri = guide.getUri(source)
	-- local description = getDesc(source)
	-- local inferred = vm.getInfer(source)
	-- local typeView = inferred:view(ws.rootUri)

	local obj = originalInit(source, has_seen)

	-- Add common custom properties here.
	--
	-- obj.myProperty = ...

	return obj
end

-- ---------------------------------------------------------------------------
-- Example: customize a particular document type
-- ---------------------------------------------------------------------------

local originalClass = export.makeDocObject["doc.class"]

export.makeDocObject["doc.class"] = function(source, obj, has_seen)
	-- Preserve LuaLS's normal handling.
	local result = originalClass(source, obj, has_seen)

	-- Examples:
	--
	-- obj.uri = guide.getUri(source)
	--
	-- local uri = guide.getUri(source)
	-- obj.absoluteFile = furi.decode(uri)
	--
	-- obj.description = getDesc(source)
	--
	-- obj.visibility = vm.getVisibleType(source)

	return result
end

local originalDocField = export.makeDocObject["doc.field"]

export.makeDocObject["doc.field"] = function(source, obj, has_seen)
	originalDocField(source, obj, has_seen)

	--[[ local hasDefines, defines = getSourceDefines(source, has_seen)

	if not hasDefines then
		return result
	end

	obj.defines = defines ]]
end

export.makeDocObject["doc.type.array"] = function(source, obj, has_seen)
	local node = source.node
	local itemType = export.documentObject(node, has_seen)

	obj.itemType = {
		type = itemType and itemType.type or node.type,
		view = itemType and itemType.view or node.view,
	}
end

-- ---------------------------------------------------------------------------
-- Example: customize global collection
-- ---------------------------------------------------------------------------

local originalGatherGlobals = export.gatherGlobals

export.gatherGlobals = function()
	local globals = originalGatherGlobals()

	-- The default implementation is essentially:
	--
	--     util.valuesOf(vm.getExportableGlobals())
	--
	-- You could filter or reorder `globals` here.

	return globals
end

-- ---------------------------------------------------------------------------
-- Example: modify the completed documentation tree
-- ---------------------------------------------------------------------------

local originalMakeDocs = export.makeDocs

export.makeDocs = function(globals, callback)
	local docs = originalMakeDocs(globals, callback)

	-- Removing the "luals.config" doc used by luals --doc_update
	for i = #docs, 1, -1 do
		if docs[i].type == "luals.config" then
			table.remove(docs, i)
		end
	end

	-- Existing real globals, plus synthetic globals as we create them.
	local byName = {}

	for _, doc in ipairs(docs) do
		if type(doc.name) == "string" then
			byName[doc.name] = doc
		end
	end

	-- Only dotted identifiers can naturally become `foo.bar.baz`.
	local function isNamedField(name)
		return type(name) == "string" and name:match("^[%a_][%w_]*$") ~= nil
	end

	local function addSynthetic(parentName, field)
		if not isNamedField(field.name) then
			return nil
		end

		local fullName = parentName .. "." .. field.name

		local doc = byName[fullName]

		if not doc then
			doc = {
				name = fullName,
				type = "variable",
				view = "unknown",
				defines = {
					field,
				},
			}

			byName[fullName] = doc
			docs[#docs + 1] = doc
		end

		return doc
	end

	local function visitDefinition(parentName, definition)
		-- These are the runtime fields discovered by our custom
		-- `vm.getFields(source.value)` handling.
		for _, field in ipairs(definition.fields or {}) do
			if isNamedField(field.name) then
				local fullName = parentName .. "." .. field.name

				addSynthetic(parentName, field)

				-- If this field itself resolved to another table/object and our
				-- setfield handler attached `fields` to it, recurse.
				visitDefinition(fullName, field)
			end
		end
	end

	-- Important: iterate only over the original globals. `docs` grows while
	-- synthetic globals are appended.
	local originalCount = #docs

	for i = 1, originalCount do
		local doc = docs[i]

		if type(doc.name) == "string" and doc.type == "variable" then
			for _, definition in ipairs(doc.defines or {}) do
				visitDefinition(doc.name, definition)
			end
		end
	end

	table.sort(docs, export.sortDoc)

	return docs
end

-- ---------------------------------------------------------------------------
-- Example: customize serialization/output
-- ---------------------------------------------------------------------------

-- local originalSerializeAndExport = export.serializeAndExport

export.serializeAndExport = function(docs, outputDir)
	-- `jsonb`
	--
	-- local json = jsonb.beautify(docs)
	--
	-- `util`
	--
	-- util.saveFile(outputDir .. '/custom.json', json)
	--
	-- `markdown`
	--
	-- local md = markdown()
	-- md:add('md', '# Custom documentation')
	-- md:emptyLine()
	-- md:add('lua', 'print("hello")')
	-- util.saveFile(outputDir .. '/custom.md', md:string())

	local old_jsonb_supportSparseArray = api.jsonb.supportSparseArray
	api.jsonb.supportSparseArray = true

	local exportPath = outputDir .. "/export"
	local variablesPath = outputDir .. "/export"
	local typesPath = outputDir .. "/export/type"

	--export all docs  to json

	local jsonPath = exportPath .. "/doc.json"
	local ok, err = api.util.saveFile(jsonPath, api.jsonb.beautify(docs))

	if not ok then
		return false, { jsonPath }, { err }
	end

	--export type docs to json

	local typeDocs = {}

	for i = 1, #docs do
		local doc = docs[i]
		if doc.type == "type" then
			table.insert(typeDocs, doc)
		end
	end

	for i = 1, #typeDocs do
		local typeDoc = typeDocs[i]
		local typePath = typesPath .. "/" .. typeDoc.name .. ".json"

		local typeOk, typeErr = api.util.saveFile(typePath, api.jsonb.beautify(typeDoc))

		if not typeOk then
			return false, { typePath }, { typeErr }
		end
	end

	--export variable docs to json

	local variableDocs = {}

	for i = 1, #docs do
		local doc = docs[i]
		if doc.type == "variable" then
			table.insert(variableDocs, doc)
		end
	end

	for i = 1, #variableDocs do
		local variableDoc = variableDocs[i]
		local variablePath = variablesPath .. "/" .. variableDoc.name .. ".json"

		local variableOk, variableErr = api.util.saveFile(variablePath, api.jsonb.beautify(variableDoc))

		if not variableOk then
			return false, { variablePath }, { variableErr }
		end
	end

	api.jsonb.supportSparseArray = old_jsonb_supportSparseArray

	return true, { jsonPath }
end

function anydev.getDocDescription(_, desc)
	if not desc then
		return {}
	end

	return split(desc, "\n")
end

function anydev.getTypeDefinition(docType)
	assert(
		docType.type == "doc.type.name",
		string.format("getDocTypeNameType received a doc with type '%s' instead", docType.type)
	)

	if docType.view == "number" then
		return {
			kind = "builtin",
			value = "number",
		}
	elseif docType.view == "string" then
		return {
			kind = "builtin",
			value = "string",
		}
	else
		return {
			kind = "typereference",
			value = docType.view,
		}
	end
end

function anydev.getAliasDocType(_, define)
	assert(
		define.type == "doc.alias",
		string.format("getAliasDocType received a doc with type '%s' instead", define.type)
	)

	local extends = define.extends

	if extends.type == "doc.type" then
		assert(#extends.types == 1, "TODO: Support multiple alias types definitions")

		local docType = extends.types[1]

		if docType.type == "doc.type.array" then
			return {
				kind = "array",
				items = anydev.getTypeDefinition(docType.itemType),
			}
		elseif docType.type == "doc.type.name" then
			return anydev.getTypeDefinition(docType)
		end
	end

	assert(false, string.format("TODO: support for '%s' type alias", extends.type))
end

-- TODO: parent
-- TODO: indexes
function anydev.getTableDocType(_, define)
	assert(
		define.type == "doc.class" or define.view == "table",
		string.format("getClassDocType received a doc with type '%s' instead", define.type)
	)

	local tableFields = {}
	local tableIndexes = {}

	for _, field in ipairs(define.fields) do
		local fieldDoc = {
			name = field.name,
			type = "field",
			defines = { field },
			desc = field.desc,
			rawDesc = field.rawDesc,
		}
		table.insert(tableFields, anydev.getDoc(fieldDoc))
	end

	return {
		kind = "table",
		fields = tableFields,
		indexes = tableIndexes,
	}
end

function anydev.getFieldDocType(_, define)
	assert(
		define.type == "doc.field",
		string.format("getFieldDocType received a doc with type '%s' instead", define.type)
	)

	local extends = define.extends

	if extends.type == "doc.type" then
		assert(#extends.types == 1, "TODO: Support multiple field type definitions")

		local docType = extends.types[1]

		if docType.type == "doc.type.array" then
			return {
				kind = "array",
				items = anydev.getTypeDefinition(docType.itemType),
			}
		elseif docType.type == "doc.type.name" then
			return anydev.getTypeDefinition(docType)
		end
	end

	assert(false, string.format("TODO: support for '%s' field type", extends.type))
end

function anydev.getSetFieldDocType(doc, define)
	assert(
		define.type == "setfield",
		string.format("getSetFieldDocType received a doc with type '%s' instead", define.type)
	)

	local extends = define.extends

	if extends.type == "select" and extends.view == "table" then
		return anydev.getTableDocType(doc, define)
	end

	-- assert(false, string.format("TODO: support for type '%s' and view '%s' setfield type", extends.type, extends.view))
end

function anydev.getDocType(doc, define)
	if define.type == "doc.alias" then
		return anydev.getAliasDocType(doc, define)
	elseif define.type == "doc.class" then
		return anydev.getTableDocType(doc, define)
	elseif define.type == "doc.field" then
		return anydev.getFieldDocType(doc, define)
	elseif define.type == "setfield" then
		return anydev.getSetFieldDocType(doc, define)
	end

	assert(false, string.format("TODO: support definition type '%s' for '%s'", define.type, doc.name))
end

local testDocs = {
	decimal = true,
	CmdContent = true,
	["lsp.AnnotatedTextEdit"] = true,
	["vim.F"] = true,
}

function anydev.getDoc(doc)
	assert(#doc.defines == 1, "Multiple doc defines received")

	local define = doc.defines[1]

	return {
		category = doc.type,
		name = doc.name,
		kind = define.type,
		documentation = anydev.getDocDescription(doc, doc.desc),
		type = anydev.getDocType(doc, define),
	}
end

local originalDocumentObject = export.documentObject

export.documentObject = function(source, has_seen)
	local doc = originalDocumentObject(source, has_seen)

	if not doc then
		return nil
	end

	if testDocs[doc.name] then
		print(doc.name)
		doc.anydev = anydev.getDoc(doc)
	end

	return doc
end

--[[ -- ---------------------------------------------------------------------------
-- Miscellaneous API examples
-- ---------------------------------------------------------------------------

local function examples(source)
	-- Workspace root URI.
	local rootUri = ws.rootUri

	-- Parser source -> URI.
	local uri = guide.getUri(source)

	-- URI -> filesystem path.
	local path = furi.decode(uri)

	-- Parser byte position -> row/column.
	local row, column
	if source.start then
		row, column = guide.rowColOf(source.start)
	end

	-- Type inference.
	local inferred = vm.getInfer(source)
	local view = inferred:view(rootUri)

	-- Documentation text.
	local description = getDesc(source)

	-- Raw documentation text.
	local rawDescription = getDesc(source, true)

	-- Label formatting. This is normally useful for function sources.
	-- local label = getLabel(source, false, 1)

	-- JSON serialization.
	local json = jsonb.beautify({
		uri = uri,
		path = path,
		row = row,
		column = column,
		type = view,
		description = description,
		rawDescription = rawDescription,
	})

	-- Filesystem.
	local cwd = fs.current_path():string()

	-- Markdown generation.
	local md = markdown()
	md:add("md", "# Example")
	md:emptyLine()
	md:add("lua", view)

	-- File writing.
	-- util.saveFile(cwd .. '/example.json', json)

	return {
		cwd = cwd,
		json = json,
		markdown = md:string(),
	}
end

-- `examples` deliberately isn't called. It exists as a small API reference
-- while developing the real exporter. ]]
