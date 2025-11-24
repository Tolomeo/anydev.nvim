package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetRoot() (string, error) {
	root, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()

	if err != nil {
		return "", fmt.Errorf("Error retrieving project folder root path: %w", err)
	}

	return strings.TrimSpace(string(root)), nil
}

func GetOutputDir() (string, error) {
	rootDir, err := GetRoot()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(rootDir, "out/")

	err = os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return "", err
	}

	return dir, nil
}

func GetConfigDir() (string, error) {
	rootDir, err := GetRoot()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(rootDir, ".config/nvim/")

	err = os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return "", err
	}

	return dir, nil
}
