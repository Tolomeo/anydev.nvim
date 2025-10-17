package main

import (
	"fmt";
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

func main() {
	nvim := nvim.NewNvim("nvim")

	fmt.Print(nvim.Path())
}
