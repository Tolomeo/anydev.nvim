package msgpackrpc

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

var mu sync.Mutex

type RequestMessage struct {
	Method string
	Params []any
}

type ResponseMessage struct {
	error  any
	result any
}

func (r ResponseMessage) Result() (any, error) {
	if r.error != nil {
		return nil, fmt.Errorf("Error response: %v", r.error)
	}

	return r.result, nil
}

type MsgpackRpc struct {
	requestId int8
	writer    io.WriteCloser
	reader    *bufio.Reader
}

func (r *MsgpackRpc) getMessageId() int8 {
	mu.Lock()
	requestId := r.requestId
	// the requestId folds and repeats from the start when we reach max count permitted by int8
	r.requestId = (r.requestId + 1) % 127
	mu.Unlock()

	return requestId
}

func (r *MsgpackRpc) Send(request RequestMessage) (*ResponseMessage, error) {
	messageId := r.getMessageId()

	// fmt.Printf("RPC request %s: [%d, %s]\n", time.Now(), messageId, request.method)
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

	// fmt.Printf("RPC response %s: [%d, %v, %v]\n", time.Now(), messageId, response.error, response.result)

	if err != nil {
		return nil, fmt.Errorf("Invalid response received: %w", err)
	}

	return response, nil
}

// https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md#request-message
// [type, msgid, method, params]
func (r *MsgpackRpc) requestToMessagePackRequest(request RequestMessage, messageId int8) ([]byte, error) {
	message := []any{
		int8(0),
		messageId,
		request.Method,
		request.Params,
	}

	return msgpack.Marshal(message)
}

// https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md#response-message
// [type, msgid, error, result]
func (r *MsgpackRpc) messagePackResponseToResponse(messagePackResponse []any, messageId int8) (*ResponseMessage, error) {
	if len(messagePackResponse) < 4 {
		return nil, fmt.Errorf("Invalid response length received: %v", r)
	}

	if messagePackResponse[0].(int8) != 1 {
		return nil, fmt.Errorf("Invalid response type received: %v", r)
	}

	if messagePackResponse[1].(int8) != messageId {
		return nil, fmt.Errorf("Out of sync response received: %v", r)
	}

	return &ResponseMessage{
		error:  messagePackResponse[2],
		result: messagePackResponse[3],
	}, nil
}

func New(cmd *exec.Cmd) (*MsgpackRpc, error) {
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

	return &MsgpackRpc{
		writer: writer,
		reader: reader,
	}, nil
}
