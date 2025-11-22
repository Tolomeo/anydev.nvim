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

	outputDir := filepath.Join(rootDir, "output/")

	err = os.MkdirAll(outputDir, os.ModePerm)

	if err != nil {
		return "", err
	}

	return outputDir, nil
}
