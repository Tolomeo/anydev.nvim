package nvim

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/msgpackrpc"
)

type CursorPosition struct {
	Line      uint
	Character uint
}

type Nvim struct {
	options options
	cmd     *exec.Cmd
	rpc     *msgpackrpc.MsgpackRpc
}

func (n *Nvim) Options() options {
	return n.options
}

func (n *Nvim) Start() error {
	return n.cmd.Start()
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

func (n *Nvim) Redir(file string) error {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_command",
		Params: []any{fmt.Sprintf("redir! %s", file)},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return err
	}

	return nil
}

var counter = 0

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

func (n *Nvim) CallFunction(function string, functionArgs []any) (any, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_call_function",
		Params: []any{function, functionArgs},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return nil, fmt.Errorf("Error executing function: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return nil, fmt.Errorf("Error executing function: %v\n", err)
	}

	return result, nil
}

func (n *Nvim) ExecLua(lua string, args []any) (any, error) {
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

func New(config Config, opts ...optionProvider) (*Nvim, error) {
	options, err := NewOptions(config, opts...)

	if err != nil {
		return nil, fmt.Errorf("Error getting nvim options: %v", err)
	}

	arguments := []string{"--embed", "--headless", "-i", "NONE", "-u", options.config.InitFile()}
	arguments = append(arguments, options.arguments...)
	cmd := exec.Command(options.command, arguments...)

	rpc, err := msgpackrpc.New(cmd)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim rpc: %v", err)
	}

	return &Nvim{
		options: options,
		cmd:     cmd,
		rpc:     rpc,
	}, nil
}
