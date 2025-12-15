package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var cache = make(map[string]string)

func GetRoot() (string, error) {
	if rootDir, cached := cache["root"]; cached {
		return rootDir, nil
	}

	root, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()

	if err != nil {
		return "", fmt.Errorf("Error retrieving project folder root path: %w", err)
	}

	rootDir := strings.TrimSpace(string(root))

	cache["root"] = rootDir

	return rootDir, nil
}

func GetOutputDir() (string, error) {
	if outputDir, cached := cache["output"]; cached {
		return outputDir, nil
	}

	rootDir, err := GetRoot()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(rootDir, "out/")

	err = os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return "", err
	}

	cache["output"] = dir

	return dir, nil
}

func GetConfigDir() (string, error) {
	if configDir, cached := cache["config"]; cached {
		return configDir, nil
	}

	rootDir, err := GetRoot()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(rootDir, ".config/nvim/")

	err = os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return "", err
	}

	cache["config"] = dir

	return dir, nil
}

func GetTmpDir() (string, error) {
	if configDir, cached := cache["tmp"]; cached {
		return configDir, nil
	}

	rootDir, err := GetRoot()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(rootDir, "tmp")

	err = os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return "", err
	}

	cache["config"] = dir

	return dir, nil
}
