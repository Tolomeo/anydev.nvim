package nvim

import (
	"path/filepath"
)

type Config struct {
	directory string
}

func (c Config) Dir() string {
	return c.directory
}

func (c Config) File(path string) string {
	return filepath.Join(c.directory, path)
}

func (c Config) InitFile() string {
	return c.File("init.lua")
}

func NewConfig(directory string) Config {
	newconfig := Config{
		directory: directory,
	}

	return newconfig
}

type options struct {
	command   string
	arguments []string
}

type optionProvider func(*options)

func WithCmd(path string) optionProvider {
	return func(o *options) {
		o.command = path
	}
}

func WithArguments(argument ...string) optionProvider {
	return func(o *options) {
		o.arguments = append(o.arguments, argument...)
	}
}

func NewOptions(opts ...optionProvider) options {
	newOptions := options{
		command: "nvim",
	}

	for _, opt := range opts {
		opt(&newOptions)
	}

	return newOptions
}
