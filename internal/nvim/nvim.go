package nvim

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/internal/scripts"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/msgpackrpc"
)

type CursorPosition struct {
	Line      uint
	Character uint
}

type Nvim struct {
	options Options
	logger  *log.Logger
	cmd     *exec.Cmd
	rpc     *msgpackrpc.MsgpackRpc
	config  string
	fileBuffers   []*FileBuffer
}

func (n *Nvim) Options() Options {
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

	script, err := scripts.Read("start")

	if err != nil {
		return err
	}

	_, err = n.execLua(script, []any{30000})

	if err != nil {
		return fmt.Errorf("Error initializing api: %w", err)
	}

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

type Options struct {
	LogLevel  uint
	Arguments []string
}

func New(options Options) (*Nvim, error) {
	cmd := exec.Command("nvim", options.Arguments...)

	rpc, err := msgpackrpc.New(cmd)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to nvim rpc: %v", err)
	}

	return &Nvim{
		options: options,
		logger:  log.NewLogger("nvim", options.LogLevel),
		cmd:     cmd,
		rpc:     rpc,
		fileBuffers:   make([]*FileBuffer, 0, 30),
	}, nil
}
