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

	fmt.Println(nvimClient.Options())

	nvimClient.ApiInfo()

	if err := nvimClient.Close(); err != nil {
		fmt.Printf("Error closing nvim gracefully: %v", err)
	}
}
