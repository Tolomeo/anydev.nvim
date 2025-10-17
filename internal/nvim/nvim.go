package nvim

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type nvim struct {
	cmd    *exec.Cmd
	writer *io.WriteCloser
	reader *bufio.Reader
}

func (n *nvim) Args() []string {
	return n.cmd.Args
}

func (n *nvim) Kill() error {
	return n.cmd.Process.Kill()
}

func NewNvim(path string) (*nvim, error) {
	cmd := exec.Command(path, "--clean", "--embed")
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
		writer: writer,
		reader: reader,
	}, nil
}
