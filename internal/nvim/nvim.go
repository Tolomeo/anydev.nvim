package nvim

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/msgpackrpc"
	"github.com/Tolomeo/anydev.nvim/internal/utils/anyx"
)

type CursorPosition struct {
	Line      uint
	Character uint
}

type Nvim struct {
	options options
	logger  *log.Logger
	cmd     *exec.Cmd
	rpc     *msgpackrpc.MsgpackRpc
	config  string
}

func (n *Nvim) Options() options {
	return n.options
}

func (n *Nvim) Start() error {
	err := n.cmd.Start()

	if err != nil {
		return err
	}

	stdpath, err := n.callFunction("stdpath", []any{"config"})

	if err != nil {
		return fmt.Errorf("Error retrieving config path: %w", err)
	}

	config, ok := stdpath.(string)

	if !ok {
		return fmt.Errorf("Error retrieving config path from value <%+v>", config)
	}

	n.config = config

	// n.config = "/root/.config/nvim"
	// n.config = "/Users/diegofrattini/Projects/anydev.nvim/resources/config"

	return nil
}

func (n *Nvim) Quit() error {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_command",
		Params: []any{"qa!"},
	}
	_, err := n.rpc.Send(request)

	// EOF error expected
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	err = n.cmd.Wait()

	if err != nil {
		return err
	}

	return nil
}

/* func (n *Nvim) redir(file string) error {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_command",
		Params: []any{fmt.Sprintf("redir! %s", file)},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return err
	}

	return nil
} */

/*
	 func (n *Nvim) ApiInfo() (any, error) {
		request := msgpackrpc.RequestMessage{
			method: "nvim_get_api_info",
			params: []any{},
		}
		response, err := n.rpc.Send(request)

		if err != nil {
			return nil, fmt.Errorf("Error getting API info: %v\n", err)
		}

		result, err := response.Result()

		if err != nil {
			return nil, fmt.Errorf("Error getting API info: %v\n", err)
		}

		return result, nil
	}
*/

func (n *Nvim) callFunction(function string, functionArgs []any) (any, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_call_function",
		Params: []any{function, functionArgs},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return nil, fmt.Errorf("Error executing function '%s' with arguments <%+v>: %w\n", function, functionArgs, err)
	}

	result, err := response.Result()

	if err != nil {
		return nil, fmt.Errorf("Error executing function '%s' with arguments <%+v>: %w\n", function, functionArgs, err)
	}

	return result, nil
}

func (n *Nvim) execLua(lua string, args []any) (any, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_exec_lua",
		Params: []any{lua, args},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return nil, fmt.Errorf("Error sending nvim_exec_lua rpc message: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return nil, fmt.Errorf("Error executing lua: %v\n", err)
	}

	return result, nil
}

func (n *Nvim) open(file string) (string, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_command",
		Params: []any{fmt.Sprintf("edit %s", file)},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return file, fmt.Errorf("Error opening %s: %v\n", file, err)
	}

	return file, nil
}

/* func (n *Nvim) write() error {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_command",
		Params: []any{"write"},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to write buffer: %v\n", err)
	}

	return nil
} */

/* func (n *Nvim) getBufferName() (string, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_buf_get_name",
		Params: []any{0},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return "", fmt.Errorf("Error reading buffer name: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return "", fmt.Errorf("Error executing lua: %v\n", err)
	}

	return result.(string), nil
} */

func (n *Nvim) setBufferLines(lines []string) error {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_buf_set_lines",
		Params: []any{0, 0, -1, true, lines},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error sending nvim_buf_set_lines rpc message: %v\n", err)
	}

	_, err = response.Result()

	if err != nil {
		return fmt.Errorf("Error setting buffer lines: %v\n", err)
	}

	return nil
}

/* func (n *Nvim) getBufferText(startRow int, startCol int, endRow int, endCol int) ([]string, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_buf_get_text",
		Params: []any{0, startRow, startCol, endRow, endCol, struct{}{}},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer text: %w\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer text: %w\n", err)
	}

	bufferText, err := anyx.ToSliceOf[string](result)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer text return value: %w", err)
	}

	return bufferText, nil
} */

func (n *Nvim) getBufferLines(start int, end int) ([]string, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_buf_get_lines",
		Params: []any{0, start, end, false},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer name: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return []string{}, fmt.Errorf("Error executing lua: %v\n", err)
	}

	bufferLines, err := anyx.ToSliceOf[string](result)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer lines return value: %w", err)
	}

	return bufferLines, nil
}

func (n *Nvim) deleteBuffer() error {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_buf_delete",
		Params: []any{0, struct{ force bool }{force: true}}}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to delete buffer: %v\n", err)
	}

	return nil
}

func New(opts ...optionProvider) (*Nvim, error) {
	options := NewOptions(opts...)

	arguments := []string{"--embed", "--headless", "-i", "NONE"}
	arguments = append(arguments, options.arguments...)
	cmd := exec.Command(options.command, arguments...)

	rpc, err := msgpackrpc.New(cmd)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim rpc: %v", err)
	}

	return &Nvim{
		options: options,
		logger:  log.NewLogger("nvim"),
		cmd:     cmd,
		rpc:     rpc,
	}, nil
}
