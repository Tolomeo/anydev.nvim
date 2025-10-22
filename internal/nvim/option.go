package nvim

import (
	"fmt"
	"os/exec"
	"strings"
)

type config struct {
	dir  string
}

func (c config) Dir() string {
	return c.dir
}

func (c config) File(path string) string {
	return c.dir + path
}

func (c config) InitFile() string {
	return c.File("init.lua")
}

func NewConfig() (config, error) {
	moduleDir, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()

	if err != nil {
		return config{}, fmt.Errorf("Error retrieving project folder path: %v", err)
	}

	configDir := strings.TrimSpace(string(moduleDir)) + "/.config/nvim/"

	newconfig := config{
		dir:  configDir,
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
