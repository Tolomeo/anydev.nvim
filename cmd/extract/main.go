package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Tolomeo/anydev.nvim/internal/extract"
)

var override = extract.Override{
	Definition: map[string][]string{
		"vim":        {"vim = {}"},
		"vim.base64": {"vim.base64 = {}"},
		"vim.cmd":    {"---@type fun(command: string|table)|table<string,fun(...:any)>", "vim.cmd = ..."},
		"vim.env":    {"---@type table<string, string>", "vim.env = ..."},
		// TODO: improve module export query
		"vim.iter": {"---@type IterMod", "vim.iter = ..."},
	},
}

func main() {
	debug := flag.Bool("debug", false, "Debug mode.")
	logLevel := flag.Uint("log-level", 2, "Log level. 0: Error, 1: Warning, 2: Info, 3: Verbose, 4: Silly")
	outDir := flag.String("out-dir", "out", "Output directory")
	tmpDir := flag.String("tmp-dir", "tmp", "Temp directory")

	flag.Parse()

	targets := flag.Args()

	if len(targets) < 1 {
		fmt.Println("Error: extraction targets required.")
		flag.Usage()
		os.Exit(1)
	}

	extractor, err := extract.NewExtractor(extract.Options{
		Debug:    *debug,
		LogLevel: *logLevel,
		OutDir:   *outDir,
		TmpDir:   *tmpDir,
		Override: override,
	})

	if err != nil {
		panic(err)
	}

	// TODO: deal with Destroy returning an error
	defer extractor.Destroy()

	for _, target := range targets {
		err := extractor.Extract("value", target)

		if err != nil {
			panic(err)
		}
	}
}
