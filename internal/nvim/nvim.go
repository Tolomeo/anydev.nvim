package nvim

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
	"github.com/Tolomeo/anydev.nvim/internal/utils/anyx"
)

type CursorPosition struct {
	Line      uint
	Character uint
}

type Location struct {
	lsp.DefinitionLocation
	Url string
}

type Nvim struct {
	options options
	cmd     *exec.Cmd
	rpc     *rpc
}

func (n *Nvim) Options() options {
	return n.options
}

func (n *Nvim) Start() error {
	return n.cmd.Start()
}

func (n *Nvim) Quit() error {
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

func (n *Nvim) Open(file string) (string, error) {
	request := requestMessage{
		method: "nvim_command",
		params: []any{"edit" + file},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return file, fmt.Errorf("Error opening %s: %v\n", file, err)
	}

	return file, nil
}

func (n *Nvim) Write() error {
	request := requestMessage{
		method: "nvim_command",
		params: []any{"write"},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to write buffer: %v\n", err)
	}

	return nil
}

func (n *Nvim) GetBufferName() (string, error) {
	request := requestMessage{
		method: "nvim_buf_get_name",
		params: []any{0},
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
}

func (n *Nvim) SetBufferLines(lines []string) error {
	request := requestMessage{
		method: "nvim_buf_set_lines",
		params: []any{0, 0, -1, true, lines},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return err
	}

	_, err = response.Result()

	if err != nil {
		return err
	}

	return nil
}

func (n *Nvim) GetBufferText(startRow int, startCol int, endRow int, endCol int) ([]string, error) {
	request := requestMessage{
		method: "nvim_buf_get_text",
		params: []any{0, startRow, startCol, endRow, endCol, struct{}{}},
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
}

func (n *Nvim) GetBufferLines(start int, end int) ([]string, error) {
	request := requestMessage{
		method: "nvim_buf_get_lines",
		params: []any{0, start, end, false},
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

func (n *Nvim) DeleteBuffer() error {
	request := requestMessage{
		method: "nvim_buf_delete",
		params: []any{0, struct{ force bool }{force: true}}}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to delete buffer: %v\n", err)
	}

	return nil
}

var counter = 0

type buffer struct {
	name     string
	previous string
	Edit     func([]string) error
	Delete   func() error
}

func (n *Nvim) Buffer() (*buffer, error) {
	counter += 1
	name := path.Join(n.Options().Config().Dir(), fmt.Sprintf("anydev.%d.lua", counter))

	previous, err := n.GetBufferName()

	if err != nil {
		return nil, err
	}

	return &buffer{
		name:     name,
		previous: previous,
		Edit: func(lines []string) error {
			_, err := n.Open(name)

			if err != nil {
				return fmt.Errorf("Error writing to buffer '%s': %w", name, err)
			}

			return n.SetBufferLines(lines)
		},
		Delete: func() error {
			_, err := n.Open(name)

			if err != nil {
				return err
			}

			err = n.DeleteBuffer()

			if err != nil {
				return err
			}

			_, err = n.Open(previous)

			if err != nil {
				return err
			}

			return nil
		},
	}, nil
}

/*
	 func (n *Nvim) ApiInfo() (any, error) {
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
	}
*/

func (n *Nvim) CallFunction(function string, functionArgs []any) (any, error) {
	request := requestMessage{
		method: "nvim_call_function",
		params: []any{function, functionArgs},
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

func (n *Nvim) GetRuntimeType(variable string) (string, error) {
	runtimePath := variable
	parts := strings.Split(runtimePath, ".")

	switch len(parts) {
	case 1:
	default:
		tail := parts[len(parts)-1]
		// https://www.lua.org/manual/5.1/manual.html#2.1
		switch tail {
		case "and", "break", "do", "else", "elseif", "end", "false", "for", "function", "if", "in", "local", "nil", "not", "or", "repeat", "return", "then", "true", "until", "while":
			head := parts[:len(parts)-1]
			runtimePath = strings.Join(head, ".") + "['" + tail + "']"
		}
	}

	luaCode := fmt.Sprintf("return type(%s)", runtimePath)

	result, err := n.ExecLua(luaCode, []any{})

	if err != nil {
		return "", fmt.Errorf("Error getting the type of %s: %w", variable, err)
	}

	typeName, ok := result.(string)

	if !ok {
		return "", fmt.Errorf("Error getting the type of %s: Error converting the result to a string", runtimePath)
	}

	return typeName, nil
}

func New(config Config, opts ...optionProvider) (*Nvim, error) {
	options, err := NewOptions(config, opts...)

	if err != nil {
		return nil, fmt.Errorf("Error getting nvim options: %v", err)
	}

	arguments := []string{"--embed", "--headless", "-u", options.config.InitFile()}
	arguments = append(arguments, options.arguments...)
	cmd := exec.Command(options.command, arguments...)
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
