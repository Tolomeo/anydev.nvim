package nvim

import (
	"bufio"
	"fmt"
	msgpack "github.com/vmihailenco/msgpack/v5"
	"io"
	"os/exec"
	"sync"
)

var mu sync.Mutex

const (
	requestMessageType  int8 = 0
	responseMessageType int8 = 1
)

type rpcMethod string

type rpc struct {
	requestId int8
	writer    io.WriteCloser
	reader    *bufio.Reader
}

func (r *rpc) Request(method rpcMethod, parameters []any) (any, error) {
	mu.Lock()
	messageId := r.requestId
	r.requestId++
	mu.Unlock()

	message := []any{
		requestMessageType, // Message type: Request (0)
		messageId,
		method,
		parameters,
	}

	messageData, err := msgpack.Marshal(message)

	if err != nil {
		return nil, fmt.Errorf("Error marshalling request message: %v", err)
	}

	if _, err := r.writer.Write(messageData); err != nil {
		return nil, fmt.Errorf("Error sending request data: %w", err)
	}

	decoder := msgpack.NewDecoder(r.reader)

	var response []any

	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("Error decoding response: %w", err)
	}

	if len(response) < 4 || response[0].(int8) != responseMessageType || response[1].(int8) != messageId {
		return nil, fmt.Errorf("Invalid response received: %v", response)
	}

	if response[2] != nil { // Check for error field
		return nil, fmt.Errorf("RPC response error: %v", response[2])
	}

	return response[3], nil
}

func NewRpc(cmd *exec.Cmd) (*rpc, error) {
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim stdin: %v", err)
	}

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim stdout: %v", err)
	}

	writer := stdin
	reader := bufio.NewReader(stdout)

	return &rpc{
		writer: writer,
		reader: reader,
	}, nil
}
