package main

import (
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

func main() {
	instance, err := nvim.NewNvim("nvim")

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %v", err))
	}

	defer instance.Kill()

	fmt.Print(instance.Args())
}
