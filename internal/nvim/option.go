package nvim

import (
	"path/filepath"
)

type Config struct {
	directory  string
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

func NewConfig(directory string) (Config) {
	newconfig := Config{
		directory:  directory,
	}

	return newconfig
}

type options struct {
	cmd       string
	arguments []string
	config    Config
}

func (o options) Config() Config {
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

func NewOptions(config Config, opts ...optionProvider) (options, error) {
	newoptions := options{
		cmd:    "nvim",
		config: config,
	}

	for _, opt := range opts {
		opt(&newoptions)
	}

	return newoptions, nil
}
