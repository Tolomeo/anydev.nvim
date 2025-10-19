package nvim

import (
	"fmt"
	"os/exec"
)

type nvim struct {
	options nvimOptions
	cmd     *exec.Cmd
	rpc     *rpc
}

func (n *nvim) Open() error {
	err := n.cmd.Start()

	return err
}

func (n *nvim) Args() []string {
	return n.cmd.Args
}

func (n *nvim) Options() nvimOptions {
	return n.options
}

func (n *nvim) Kill() error {
	return n.cmd.Process.Kill()
}

func (n *nvim) ApiInfo() {
	apiInfo, err := n.rpc.Request("nvim_get_api_info", []any{})

	if err != nil {
		fmt.Printf("Error getting API info: %v\n", err)
		return
	}

	fmt.Printf("Api info: %v\n", apiInfo)
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
