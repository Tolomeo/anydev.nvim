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

	fmt.Println(vimrc)

	nvimClient, err := nvim.New(nvim.WithVimrc(vimrc))

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %v", err))
	}

	err = nvimClient.Open()

	if err != nil {
		panic(fmt.Errorf("Error opening nvim: %v", err))

	}

	fmt.Println(nvimClient.Options())

	apiInfo, err := nvimClient.ApiInfo()

	if err != nil {
		panic(fmt.Errorf("Error obtaining nvim api info: %v", err))
	}

	fmt.Println(apiInfo)

	if err := nvimClient.Close(); err != nil {
		fmt.Printf("Error closing nvim gracefully: %v", err)
	}
}
