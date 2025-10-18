package nvim

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type nvim struct {
	cmd     *exec.Cmd
	options nvimOptions
	writer  *io.WriteCloser
	reader  *bufio.Reader
}

func (n *nvim) Args() []string {
	return n.cmd.Args
}

func (n *nvim) Options() nvimOptions {
	return n.options
}

func (n *nvim) Kill() error {
	return n.cmd.Process.Kill()
}

type nvimOptions struct {
	path string
}

type nvimOptionProvider func(*nvimOptions)

func WithPath(path string) nvimOptionProvider {
	return func(o *nvimOptions) {
		o.path = path
	}
}

func NewNvim(opts ...nvimOptionProvider) (*nvim, error) {
	options := nvimOptions{
		path: "nvim",
	}

	for _, opt := range opts {
		opt(&options)
	}

	cmd := exec.Command(options.path, "--clean", "--embed")
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim stdin: %v", err)
	}

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim stdout: %v", err)
	}

	writer := &stdin
	reader := bufio.NewReader(stdout)

	return &nvim{
		cmd:    cmd,
		options: options,
		writer: writer,
		reader: reader,
	}, nil
}
