---@diagnostic disable: undefined-global

--[[
Custom LuaLS documentation exporter.

LuaLS loads this file after cloning the built-in `cli.doc.export`
module. The following globals are injected into this script:

    export
    ws
    vm
    guide
    getDesc
    getLabel
    jsonb
    util
    markdown

The normal Lua global environment is also available through _G.

The built-in exporter additionally uses these modules, but they are
not currently injected into the custom-script environment:
]]

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

local originalSetfield = export.makeDocObject["setfield"]

export.makeDocObject["setfield"] = function(source, obj, has_seen)
	local result = originalSetfield(source, obj, has_seen)

	if not source.value then
		return result
	end

	local fields = vm.getFields(source.value)

	if #fields == 0 then
		return result
	end

	obj.fields = {}

	for _, child in ipairs(fields) do
		local childObj = export.documentObject(child, has_seen)

		if childObj then
			obj.fields[#obj.fields + 1] = childObj
		end
	end

	table.sort(obj.fields, export.sortDoc)

	return result
end

local originalVariable = export.makeDocObject["variable"]

export.makeDocObject["variable"] = function(source, obj, has_seen)
	local result = originalVariable(source, obj, has_seen)

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

	-- `docs` is the final Lua table before serialization.
	--
	-- This is often the easiest place to make broad changes:
	--
	-- for _, doc in ipairs(docs) do
	--     doc.customField = true
	-- end

	return docs
end

-- ---------------------------------------------------------------------------
-- Example: customize the generated LuaLS config object
-- ---------------------------------------------------------------------------

local originalGetLualsConfig = export.getLualsConfig

export.getLualsConfig = function()
	local config = originalGetLualsConfig()

	-- `fs` is useful whenever you need filesystem manipulation.
	--
	-- local cwd = fs.current_path():string()
	--
	-- `furi` converts between filesystem paths and file:// URIs.
	--
	-- local uri  = furi.encode(cwd)
	-- local path = furi.decode(uri)

	return config
end

-- ---------------------------------------------------------------------------
-- Example: customize serialization/output
-- ---------------------------------------------------------------------------

local originalSerializeAndExport = export.serializeAndExport

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

	-- Preserve the normal doc.json + doc.md output.
	return originalSerializeAndExport(docs, outputDir)
end

-- ---------------------------------------------------------------------------
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
-- while developing the real exporter.
