package nvim

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type nvim struct {
	options nvimOptions
	cmd     *exec.Cmd
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

func New(opts ...nvimOptionProvider) (*nvim, error) {
	options := nvimOptions{
		path: "nvim",
	}
	options.Set(opts...)

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
		options: options,
		cmd:     cmd,
		writer:  writer,
		reader:  reader,
	}, nil
}
