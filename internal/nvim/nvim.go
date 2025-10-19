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

func (n *nvim) Args() []string {
	return n.cmd.Args
}

func (n *nvim) Options() nvimOptions {
	return n.options
}

func (n *nvim) Kill() error {
	return n.cmd.Process.Kill()
}

func New(opts ...nvimOptionProvider) (*nvim, error) {
	options := nvimOptions{
		path: "nvim",
	}
	options.Set(opts...)

	cmd := exec.Command(options.path, "--clean", "--embed")
	rpc, err := Rpc(cmd)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim rpc: %v", err)
	}

	return &nvim{
		options: options,
		cmd:     cmd,
		rpc:     rpc,
	}, nil
}
