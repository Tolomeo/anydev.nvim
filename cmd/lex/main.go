package main

import (
	"fmt"
	// "path"

	// "github.com/Tolomeo/anydev.nvim/internal/context"
	"github.com/Tolomeo/anydev.nvim/internal/extract"
	// "github.com/Tolomeo/anydev.nvim/internal/lex"
	// "github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	// "github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/project"
	// "github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

const debug = true

var values []string = []string{"vim.F"}

// var values []string = []string{"vim.validate"}

// var values []string = []string{"vim.validate", "vim.F"}
// var values []string = []string{"vim.loop"}
// var values []string = []string{"vim.F", "vim.validate", "vim.loop"}

// var values = []string{}

// var types = []string{"uv.interface_addresses.addr"}
var types = []string{}

// var paths = []string{"vip"}
// var paths = []string{"uv"}

/* func getClient() (*nvim.Nvim, error) {
	configDir, err := project.GetConfigDir()
	tmpDir, err := project.GetTmpDir()

	if err != nil {
		return nil, fmt.Errorf("Error reading project directories: %w", err)
	}

	nvimConfig := nvim.NewConfig(configDir)

	var client *nvim.Nvim

	if !debug {
		client, err = nvim.New(nvimConfig)
	} else {
		client, err = nvim.New(
			nvimConfig,
			nvim.WithArguments(
				fmt.Sprintf("-V%d%s", 10, path.Join(tmpDir, "nvim.verbosefile")),
				"--listen", path.Join(tmpDir, "nvim.server.pipe"),
			),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("Error initialising nvim client: %v", err)
	}

	err = client.Start()

	if err != nil {
		return nil, fmt.Errorf("Error starting nvim client: %w", err)
	}

	return client, nil
} */

func getOutput() (*output.Output, error) {
	outputDir, err := project.GetOutputDir()

	if err != nil {
		return nil, fmt.Errorf("Error getting output location: %w", err)
	}

	out := output.NewOutput(outputDir)

	return out, nil
}

func main() {
	// var paths []string = []string{"vim._defer_require", "vim.deepcopy", "vim.validate"}

	/* client, err := getClient()

	if err != nil {
		panic(err)
	} */

	// lexer := lex.NewLexer()
	out, err := getOutput()

	if err != nil {
		panic(err)
	}

	extractor, err := extract.NewExtractor(extract.Options{Debug: true})

	if err != nil {
		panic(err)
	}

	for _, value := range values {
		err := extractor.Extract("value", value)

		/* logger := log.NewLogger("")
		context := context.New(logger, client)
		err = lexer.LexValue(value, context) */
		// err = lexer.LexType(path, context)

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

	/* for _, type_ := range types {
		logger := log.NewLogger("")
		context := context.New(logger, client)
		err = lexer.LexType(type_, context)

		if err != nil {
			panic(err)
		}

		if err := out.WriteFile(fmt.Sprintf("%s.result.json", type_), context.Result()); err != nil {
			panic(fmt.Errorf("Error writing result.json: %w", err))
		}

		if err := out.WriteFile(fmt.Sprintf("%s.logs.json", type_), logger.Logs()); err != nil {
			panic(fmt.Errorf("Error writing logs.json: %w", err))
		}
	} */

	/* fmt.Println("hey")
	fmt.Scanln() */

	/* err = client.Quit()

	if err != nil {
		return panic(err)
	} */
}
