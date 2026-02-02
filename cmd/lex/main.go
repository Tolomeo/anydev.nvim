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
var values []string = []string{"vim.lsp.protocol.Methods"}

// var values []string = []string{"vim.F", "vim.validate", "vim.loop"}

// var values = []string{}

// var types = []string{"uv.fs_copyfile.flags"}
// var types = []string{"vim.lsp.protocol.Methods"}
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

	extractor, err := extract.NewExtractor(extract.Options{Debug: debug})

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
