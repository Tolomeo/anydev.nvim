package main

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/context"
	"github.com/Tolomeo/anydev.nvim/internal/lex"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/output"
	"github.com/Tolomeo/anydev.nvim/internal/project"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

const debug = true

// var paths []string = []string{"vim.loop"}
var paths []string = []string{"vim.F", "vim.validate"}

func getClient() (*nvim.Nvim, error) {
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
}

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

	client, err := getClient()

	if err != nil {
		panic(err)
	}

	logger := log.NewLogger("")
	lexer := lex.NewLexer()
	out, err := getOutput()

	for _, path := range paths {
		context := context.New(logger, client)
		err = lexer.LexValue(path, context)

		if err != nil {
			/* fmt.Println("Errorrrrr")
			var input string
			_, _ = fmt.Scanln(&input) */
			panic(err)
		}

		if err := out.WriteFile(fmt.Sprintf("%s.result.json", path), context.Result()); err != nil {
			panic(fmt.Errorf("Error writing result.json: %w", err))
		}

		if err := out.WriteFile(fmt.Sprintf("%s.logs.json", path), logger.Logs()); err != nil {
			panic(fmt.Errorf("Error writing logs.json: %w", err))
		}
	}

	/* fmt.Println("hey")
	fmt.Scanln() */

	/* err = client.Quit()

	if err != nil {
		return panic(err)
	} */
}
