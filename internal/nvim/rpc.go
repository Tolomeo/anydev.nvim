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
	responseMessageType int8 = 1
)

type rpcRequest struct {
	method string
	params []any
}

func (r rpcRequest) Marshal(id int8) []any {
	return []any{
		int8(0), // Message type: Request (0)
		id,
		r.method,
		r.params,
	}

}

type rpcResponse []any

func (r rpcResponse) Validate(id int8) error {
	if len(r) < 4 {
		return fmt.Errorf("Invalid response length received: %v", r)
	}

	if r[0].(int8) != 1 {
		return fmt.Errorf("Invalid response type received: %v", r)
	}

	if r[1].(int8) != id {
		return fmt.Errorf("Out of sync response received: %v", r)
	}

	return nil
}

func (r rpcResponse) Result() (any, error) {
	if r[2] != nil {
		return nil, fmt.Errorf("Error response: %v", r[2])
	}

	return r[3], nil
}

type rpc struct {
	requestId int8
	writer    io.WriteCloser
	reader    *bufio.Reader
}

func (r *rpc) Send(request rpcRequest) (rpcResponse, error) {
	mu.Lock()
	messageId := r.requestId
	r.requestId++
	mu.Unlock()

	messageData, err := msgpack.Marshal(request.Marshal(messageId))

	if err != nil {
		return nil, fmt.Errorf("Error marshalling request message: %v", err)
	}

	if _, err := r.writer.Write(messageData); err != nil {
		return nil, fmt.Errorf("Error sending request data: %w", err)
	}

	decoder := msgpack.NewDecoder(r.reader)

	var response rpcResponse

	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("Error decoding response: %w", err)
	}

	if err := response.Validate(messageId); err != nil {
		return nil, fmt.Errorf("Invalid response received: %v", response)
	}

	return response, nil
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
