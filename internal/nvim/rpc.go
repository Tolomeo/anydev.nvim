package nvim

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

var mu sync.Mutex

type requestMessage struct {
	method string
	params []any
}

type responseMessage struct {
	error  any
	result any
}

func (r responseMessage) Result() (any, error) {
	if r.error != nil {
		return nil, fmt.Errorf("Error response: %v", r.error)
	}

	return r.result, nil
}

type rpc struct {
	requestId int8
	writer    io.WriteCloser
	reader    *bufio.Reader
}

func (r *rpc) Send(request requestMessage) (*responseMessage, error) {
	mu.Lock()
	messageId := r.requestId
	r.requestId++
	mu.Unlock()

	messagePackRequest, err := r.requestToMessagePackRequest(request, messageId)

	if err != nil {
		return nil, fmt.Errorf("Error marshalling request message: %w", err)
	}

	if _, err := r.writer.Write(messagePackRequest); err != nil {
		return nil, fmt.Errorf("Error sending request data: %w", err)
	}

	decoder := msgpack.NewDecoder(r.reader)

	var messagePackResponse []any

	if err := decoder.Decode(&messagePackResponse); err != nil {
		return nil, fmt.Errorf("Error decoding response: %w", err)
	}

	response, err := r.messagePackResponseToResponse(messagePackResponse, messageId)

	if err != nil {
		return nil, fmt.Errorf("Invalid response received: %w", err)
	}

	return response, nil
}

// https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md#request-message
// [type, msgid, method, params]
func (r *rpc) requestToMessagePackRequest(request requestMessage, messageId int8) ([]byte, error) {
	message := []any{
		int8(0),
		messageId,
		request.method,
		request.params,
	}

	return msgpack.Marshal(message)
}

// https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md#response-message
// [type, msgid, error, result]
func (r *rpc) messagePackResponseToResponse(messagePackResponse []any, messageId int8) (*responseMessage, error) {
	if len(messagePackResponse) < 4 {
		return nil, fmt.Errorf("Invalid response length received: %v", r)
	}

	if messagePackResponse[0].(int8) != 1 {
		return nil, fmt.Errorf("Invalid response type received: %v", r)
	}

	if messagePackResponse[1].(int8) != messageId {
		return nil, fmt.Errorf("Out of sync response received: %v", r)
	}

	return &responseMessage{
		error:  messagePackResponse[2],
		result: messagePackResponse[3],
	}, nil
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
