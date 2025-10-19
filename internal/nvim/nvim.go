package nvim

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type nvim struct {
	options nvimOptions
	cmd     *exec.Cmd
	rpc     *rpc
}

func (n *nvim) Options() nvimOptions {
	return n.options
}

func (n *nvim) Open() error {
	return n.cmd.Start()
}

func (n *nvim) Close() error {
	request := rpcRequest{
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

func (n *nvim) ApiInfo() {
	request := rpcRequest{
		method: "nvim_get_api_info",
		params: []any{},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		fmt.Printf("Error getting API info: %v\n", err)
		return
	}

	result, err := response.Result()

	if err != nil {
		fmt.Printf("Error getting API info: %v\n", err)
		return
	}

	fmt.Printf("Api info: %v\n", result)
}

func New(optionOverrides ...nvimOptionProvider) (*nvim, error) {
	options := nvimOptions{
		path: "nvim",
	}
	options.Set(optionOverrides...)

	cmd := exec.Command(options.path, "--clean", "--embed")
	rpc, err := NewRpc(cmd)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim rpc: %v", err)
	}

	return &nvim{
		options: options,
		cmd:     cmd,
		rpc:     rpc,
	}, nil
}
