package nvim

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"reflect"
	"sync"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

var mu sync.Mutex

type requestMessage []any

func (r requestMessage) Validate() error {
	if len(r) < 2 {
		return fmt.Errorf("Invalid request message length received: %v", r)
	}

	if reflect.TypeOf(r[0]).String() != "string" {
		return fmt.Errorf("Invalid request message method type: %v", r)
	}

	if reflect.TypeOf(r[1]).String() != "[]interface {}" {
		return fmt.Errorf("Invalid request message parameters type: %v", r)
	}

	return nil
}

// https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md#request-message
// [type, msgid, method, params]
func (r requestMessage) Marshal(id int8) []any {
	return []any{
		int8(0),
		id,
		r[0],
		r[1],
	}
}

type responseMessage []any

// https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md#response-message
// [type, msgid, error, result]
func (r responseMessage) Validate(id int8) error {
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

func (r responseMessage) Result() (any, error) {
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

func (r *rpc) Send(request requestMessage) (responseMessage, error) {
	mu.Lock()
	messageId := r.requestId
	r.requestId++
	mu.Unlock()

	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid request received: %v", err)
	}

	message, err := msgpack.Marshal(request.Marshal(messageId))

	if err != nil {
		return nil, fmt.Errorf("Error marshalling request message: %v", err)
	}

	if _, err := r.writer.Write(message); err != nil {
		return nil, fmt.Errorf("Error sending request data: %w", err)
	}

	decoder := msgpack.NewDecoder(r.reader)

	var response responseMessage

	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("Error decoding response: %w", err)
	}

	if err := response.Validate(messageId); err != nil {
		return nil, fmt.Errorf("Invalid response received: %v", err)
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
