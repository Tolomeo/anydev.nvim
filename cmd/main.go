package main

import (
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"os/exec"
	"strings"
)

func nvimConfig() (string, error) {
	cmdOutput, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()

	if err != nil {
		return "", err
	}

	rootDir := strings.TrimSpace(string(cmdOutput))
	configDir := rootDir + "/.config/nvim"

	return configDir, nil
}

func main() {
	vimrc, err := nvimConfig()

	if err != nil {
		panic(fmt.Errorf("Error retrieving lua init: %v", err))
	}

	nvimClient, err := nvim.New(nvim.WithVimrc(vimrc))

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %v", err))
	}

	err = nvimClient.Open()

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %v", err))
	}

	luacode := "return vim.fn.json_encode(vim.lsp.config.lua_ls)"

	result, err := nvimClient.ExecLua(luacode, []any{})

	if (err != nil) {
		panic(fmt.Errorf("Error executing lua: %v", err))
	}

	fmt.Println(result)

	if err := nvimClient.Close(); err != nil {
		fmt.Printf("Error closing nvim gracefully: %v", err)
	}
}
