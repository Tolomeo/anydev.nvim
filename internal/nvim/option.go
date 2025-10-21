package nvim

import (
	"fmt"
	"os/exec"
	"strings"
)

func vimrc() (string, error) {
	cmdOutput, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()

	if err != nil {
		return "", err
	}

	rootDir := strings.TrimSpace(string(cmdOutput))
	configDir := rootDir + "/.config/nvim"
	initLua := configDir + "/init.lua"

	return initLua, nil
}

type options struct {
	cmd   string
	vimrc string
}

func (o *options) Set(opts ...optionProvider) {
	for _, opt := range opts {
		opt(o)
	}
}

type optionProvider func(*options)

func WithCmd(path string) optionProvider {
	return func(o *options) {
		o.cmd = path
	}
}

func WithVimrc(vimrc string) optionProvider {
	return func(o *options) {
		o.vimrc = vimrc
	}
}

func NewOptions(opts ...optionProvider) (options, error) {
	rc, err := vimrc()

	if err != nil {
		return options{}, fmt.Errorf("Error retrieving default nvim config location: %v", err)
	}

	newoptions := options{
		cmd:   "nvim",
		vimrc: rc,
	}

	for _, opt := range opts {
		opt(&newoptions)
	}

	return newoptions, nil
}
