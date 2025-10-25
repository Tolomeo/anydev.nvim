package main

import (
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

func main() {
	nvimClient, err := nvim.New()

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %v", err))
	}

	err = nvimClient.Open()

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %v", err))
	}

	tempFile := nvimClient.Options().Config().Dir() + "anydev.lua"
	err = nvimClient.Edit(tempFile)

	if err != nil {
		panic(err)
	}

	bufferName, err := nvimClient.GetBufferName()

	if err != nil {
		panic(err)
	}

	fmt.Println(bufferName)

	err = nvimClient.SetBufferLines([]string{
		"local vim_api = vim",
	})

	if err != nil {
		panic(err)
	}

	lines, err := nvimClient.GetBufferLines()

	if err != nil {
		panic(err)
	}

	fmt.Println(lines)

	documentSymbols, err := nvimClient.GetDocumentSymbols()

	if err != nil {
		panic(err)
	}

	fmt.Println(documentSymbols)

	if err := nvimClient.Close(); err != nil {
		fmt.Printf("Error closing nvim gracefully: %v", err)
	}
}
