package nvim

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	// msgpack "github.com/vmihailenco/msgpack/v5"
)

const (
	rpcRequest  int8 = 0
	rpcResponse int8 = 1
)

type rpc struct {
	requestId int8
	writer    *io.WriteCloser
	reader    *bufio.Reader
}

func Rpc(cmd *exec.Cmd) (*rpc, error) {
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

	return &rpc{
		writer: writer,
		reader: reader,
	}, nil
}
