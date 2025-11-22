package nvim

import (
	"fmt"
	"path/filepath"
	"github.com/Tolomeo/anydev.nvim/internal/utils/project"
)

type config struct {
	directory  string
}

func (c config) Dir() string {
	return c.directory
}

func (c config) File(path string) string {
	return filepath.Join(c.directory, path)
}

func (c config) InitFile() string {
	return c.File("init.lua")
}

func NewConfig() (config, error) {
	rootDir, err := project.GetRoot()

	if err != nil {
		return config{}, fmt.Errorf("Error retrieving project folder path: %v", err)
	}

	configDir := filepath.Join(rootDir, ".config/nvim/")

	newconfig := config{
		directory:  configDir,
	}

	return newconfig, nil
}

type options struct {
	cmd       string
	arguments []string
	config    config
}

func (o options) Config() config {
	return o.config
}

type optionProvider func(*options)

func WithCmd(path string) optionProvider {
	return func(o *options) {
		o.cmd = path
	}
}

func WithArguments(arguments []string) optionProvider {
	return func(o *options) {
		o.arguments = arguments
	}
}

func NewOptions(opts ...optionProvider) (options, error) {
	config, err := NewConfig()

	if err != nil {
		return options{}, fmt.Errorf("Error retrieving default nvim config details: %v", err)
	}

	newoptions := options{
		cmd:    "nvim",
		config: config,
	}

	for _, opt := range opts {
		opt(&newoptions)
	}

	return newoptions, nil
}
