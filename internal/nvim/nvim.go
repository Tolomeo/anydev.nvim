package nvim

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type Nvim struct {
	options options
	cmd     *exec.Cmd
	rpc     *rpc
}

func (n *Nvim) Options() options {
	return n.options
}

func (n *Nvim) Open() error {
	return n.cmd.Start()
}

func (n *Nvim) Close() error {
	request := requestMessage{
		method: "nvim_command",
		params: []any{"qa!"},
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

/* func (n *Nvim) ApiInfo() (any, error) {
	request := requestMessage{
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
} */

func (n *Nvim) ExecLua(lua string, args []any) (any, error) {
	request := requestMessage{
		method: "nvim_exec_lua",
		params: []any{lua, args},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return nil, fmt.Errorf("Error executing lua: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return nil, fmt.Errorf("Error executing lua: %v\n", err)
	}

	return result, nil
}

func New(opts ...optionProvider) (*Nvim, error) {
	options, err := NewOptions(opts...)

	if err != nil {
		return nil, fmt.Errorf("Error getting nvim options: %v", err)
	}

	arguments := []string{"--embed", "--headless"}

	if options.vimrc != "" {
		arguments = append(arguments, "-u", options.vimrc)
	}

	cmd := exec.Command(options.cmd, arguments...)
	rpc, err := NewRpc(cmd)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim rpc: %v", err)
	}

	return &Nvim{
		options: options,
		cmd:     cmd,
		rpc:     rpc,
	}, nil
}
