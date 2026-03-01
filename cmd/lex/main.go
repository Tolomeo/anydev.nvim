package main

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/extract"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/project"
)

const debug = true

// var values []string = []string{"vim.F"}
// var values []string = []string{"vim.validate"}
// var values []string = []string{"vim.validate", "vim.F"}
// var values []string = []string{"vim.loop"}
// var values []string = []string{"vim.log"}
// var values []string = []string{"vim.lsp"}
// var values []string = []string{"vim.treesitter"}
// var values []string = []string{"vim.api"}
// var values []string = []string{"vim.b"}
// var values []string = []string{"vim.opt"}
// "vim.bo.ai",

var values = []string{
	// "vim.F", "vim.validate", "vim.loop",
	// "vim.lsp",
	// "vim.treesitter",
	// "vim.fn.function", // named as a keyword
	// "vim.fn.NetUserPass", // symbol unknown to the lsp
	// "vim.func",

	// "vim.validate", "vim.loop", "vim.log", "vim.lsp",

	// "vim.F",
	// "vim.NIL",
	// "vim.api",
	// "vim.b",
	// "vim.base64",
	// "vim.bo",
	// "vim.call",
	// "vim.cmd",
	// "vim.deep_equal",
	// "vim.deepcopy",
	// "vim.defaulttable",
	// "vim.defer_fn",
	// "vim.deprecate",
	// "vim.diagnostic",
	// "vim.diff",
	// "vim.empty_dict",
	// "vim.endswith",
	// "vim.env",
	// "vim.filetype",
	// "vim.fn",
	// "vim.fs",
	// "vim.func",
	// "vim.funcref",
	// "vim.g",
	// "vim.glob",
	// "vim.go",
	// "vim.gsplit",
	// "vim.health",
	// "vim.highlight",
	// "vim.hl",
	// "vim.iconv",
	// "vim.in_fast_event",
	// "vim.inspect",
	// "vim.inspect_pos",
	// "vim.is_callable",
	// "vim.is_thread",
	// "vim.isarray",
	// "vim.islist",
	// "vim.iter",
	// "vim.json",
	// "vim.keycode",
	// "vim.keymap",
	// "vim.list_contains",
	// "vim.list_extend",
	// "vim.list_slice",
	// "vim.loader",
	// "vim.log",
	// "vim.loop",
	// "vim.lpeg",
	// "vim.lsp",
	// "vim.lua_omnifunc",
	// "vim.mpack",
	// "vim.notify",
	// "vim.notify_once",
	// "vim.o",
	// "vim.on_key",
	// "vim.opt",
	// "vim.opt_global",
	// "vim.opt_local",
	// "vim.paste",
	// "vim.pesc",
	// "vim.print",
	// "vim.provider",
	// "vim.re",
	// "vim.regex",
	// "vim.region",
	// "vim.ringbuf",
	// "vim.rpcnotify",
	// "vim.rpcrequest",
	// "vim.schedule",
	// "vim.schedule_wrap",
	// "vim.secure",
	// "vim.show_pos",
	// "vim.snippet",
	// "vim.spairs",
	// "vim.spell",
	// "vim.split",
	// "vim.startswith",
	// "vim.str_byteindex",
	// "vim.str_utf_end",
	// "vim.str_utf_pos",
	// "vim.str_utf_start",
	// "vim.str_utfindex",
	// "vim.stricmp",
	// "vim.system",
	// "vim.t",
	// "vim.tbl_add_reverse_lookup",
	// "vim.tbl_contains",
	// "vim.tbl_count",
	// "vim.tbl_deep_extend",
	// "vim.tbl_extend",
	// "vim.tbl_filter",
	// "vim.tbl_flatten",
	// "vim.tbl_get",
	// "vim.tbl_isempty",
	// "vim.tbl_islist",
	// "vim.tbl_keys",
	// "vim.tbl_map",
	// "vim.tbl_values",
	// "vim.text",
	// "vim.treesitter",
	// "vim.trim",
	// "vim.type_idx",
	// "vim.types",
	// "vim.ui",
	// "vim.ui_attach",
	// "vim.ui_detach",
	// "vim.uri_encode",
	// "vim.uri_from_bufnr",
	// "vim.uri_from_fname",
	// "vim.uri_to_bufnr",
	// "vim.uri_to_fname",
	// "vim.uv",
	// "vim.v",
	// "vim.val_idx",
	// "vim.validate",
	// "vim.version",
	// "vim.w",
	// "vim.wait",
	// "vim.wo",

	"vim",
}

// var values []string = []string{"vim.F", "vim.validate", "vim.loop", "vim.log", "vim.lsp"}
// var values = []string{"vim.lsp.config"}
// var values = []string{"vim.lsp.log_levels"}
// var values = []string{"vim.lsp.client_errors"}
// var values = []string{}

// var types = []string{"uv.fs_copyfile.flags"}
// var types = []string{"vim.lsp.protocol.Methods"}
// var types = []string{"vim.log.levels"}
// var types = []string{"vim.lsp.Client"}
// var types = []string{"vim.Ringbuf"}
// var types = []string{"elem_or_list"}
// var types = []string{"lsp.DynamicCapabilities"}
// var types = []string{"vim.lsp.Client"}
// var types = []string{"lsp.CallHierarchyClientCapabilities"}
// var types = []string{"Range2"}
var types = []string{}

func getOutput() (*output.Output, error) {
	outputDir, err := project.GetOutputDir()

	if err != nil {
		return nil, fmt.Errorf("Error getting output location: %w", err)
	}

	out := output.NewOutput(outputDir)

	return out, nil
}

func main() {
	out, err := getOutput()

	if err != nil {
		panic(err)
	}

	override := extract.Override{
		Definition: map[string][]string{
			"vim":        {"vim = {}"},
			"vim.base64": {"vim.base64 = {}"},
			"vim.cmd":    {"---@type fun(command: string|table)|table<string,fun(...:any)>", "vim.cmd = ..."},
			"vim.env":    {"---@type table<string, string>", "vim.env = ..."},
			// TODO: improve module export query
			"vim.iter": {"---@type IterMod", "vim.iter = ..."},
		},
	}

	extractor, err := extract.NewExtractor(extract.Options{Debug: debug, Override: override})

	if err != nil {
		panic(err)
	}

	for _, value := range values {
		err := extractor.Extract("value", value)

		if err != nil {
			panic(err)
		}

		result := extractor.Result()

		if err := out.WriteFile(fmt.Sprintf("%s.result.json", value), result); err != nil {
			panic(fmt.Errorf("Error writing result.json: %w", err))
		}

		/* if err := out.WriteFile(fmt.Sprintf("%s.logs.json", value), logger.Logs()); err != nil {
			panic(fmt.Errorf("Error writing logs.json: %w", err))
		} */
		extractor.Flush()
	}

	for _, typ := range types {
		err := extractor.Extract("type", typ)

		if err != nil {
			panic(err)
		}

		result := extractor.Result()

		if err := out.WriteFile(fmt.Sprintf("%s.result.json", typ), result); err != nil {
			panic(fmt.Errorf("Error writing result.json: %w", err))
		}

		/* if err := out.WriteFile(fmt.Sprintf("%s.logs.json", value), logger.Logs()); err != nil {
			panic(fmt.Errorf("Error writing logs.json: %w", err))
		} */
		extractor.Flush()
	}
}
