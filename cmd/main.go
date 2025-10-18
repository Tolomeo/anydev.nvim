package main

import (
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

func main() {
	nvimInstance, err := nvim.NewNvim()

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %v", err))
	}

	fmt.Print(nvimInstance.Args())
	fmt.Print(nvimInstance.Options())
}
