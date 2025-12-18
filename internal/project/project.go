package project

import (
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/utils/cache"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var projectCache = cache.NewCache[string]()

func GetRoot() (string, error) {
	if rootDir, cached := projectCache.Get("root"); cached {
		return rootDir, nil
	}

	root, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()

	if err != nil {
		return "", fmt.Errorf("Error retrieving project folder root path: %w", err)
	}

	rootDir := strings.TrimSpace(string(root))

	projectCache.Set(rootDir, "root")

	return rootDir, nil
}

func GetResourcesDir() (string, error) {
	if outputDir, cached := projectCache.Get("resources"); cached {
		return outputDir, nil
	}

	rootDir, err := GetRoot()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(rootDir, "resources/")

	projectCache.Set(dir, "output")

	return dir, nil
}

func GetOutputDir() (string, error) {
	if outputDir, cached := projectCache.Get("out"); cached {
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

	projectCache.Set(dir, "out")

	return dir, nil
}

func GetConfigDir() (string, error) {
	if configDir, cached := projectCache.Get("config"); cached {
		return configDir, nil
	}

	resourcesDir, err := GetResourcesDir()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(resourcesDir, "config/")

	projectCache.Set(dir, "config")

	return dir, nil
}

func GetTmpDir() (string, error) {
	if configDir, cached := projectCache.Get("tmp"); cached {
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

	projectCache.Set(dir, "config")

	return dir, nil
}
