package nvim

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/Tolomeo/anydev.nvim/internal/lsp"
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

func (n *Nvim) Edit(file string) error {
	request := requestMessage{
		method: "nvim_command",
		params: []any{"edit" + file},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to edit %s: %v\n", file, err)
	}

	return nil
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
		params: []any{0, 0, -1, false, lines},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying set buffer lines: %v\n", err)
	}

	return nil

}

func (n *Nvim) GetBufferLines() ([]string, error) {
	request := requestMessage{
		method: "nvim_buf_get_lines",
		params: []any{0, 0, -1, false},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer name: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return []string{}, fmt.Errorf("Error executing lua: %v\n", err)
	}

	sliceOfAny, ok := result.([]any)

	if !ok {
		return []string{}, fmt.Errorf("Error reading buffer lines return value: $v")
	}

	sliceOfStrings := make([]string, len(sliceOfAny))

	for index, value := range sliceOfAny {
		str, ok := value.(string)

		if !ok {
			return []string{}, fmt.Errorf("Error reading buffer lines return value item %d: %s", index, value)
		}

		sliceOfStrings[index] = str
	}

	return sliceOfStrings, nil
}

func (n *Nvim) DeleteBuffer() error {
	request := requestMessage{
		method: "nvim_buf_delete",
		params: []any{0, struct{ force bool }{force: true}}}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to write buffer: %v\n", err)
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

func (n *Nvim) WaitForLSP() error {
	luaCode := `
		vim.wait(2000, function()
			return next(vim.lsp.get_active_clients()) ~= nil
		end)
	`

	_, err := n.ExecLua(luaCode, []any{})

	if err != nil {
		return fmt.Errorf("Error waiting for lsp to attach: %v", err)
	}

	return nil
}

func (n *Nvim) GetDocumentSymbols() (lsp.TextDocumentDocumentSymbolResponse, error) {
	documentSymbols := lsp.TextDocumentDocumentSymbolResponse{}

	err := n.WaitForLSP()

	if err != nil {
		return documentSymbols, err
	}

	luaCode := `
		local textDocumentParams = vim.lsp.util.make_text_document_params(0)
		local result = vim.lsp.buf_request_sync(0, 'textDocument/documentSymbol', { textDocument = textDocumentParams }, 2000)
		return vim.fn.json_encode(result[1])
	`

	result, err := n.ExecLua(luaCode, []any{})

	if err != nil {
		return documentSymbols, fmt.Errorf("Error getting document symbols: %v", err)
	}

	stringResult, ok := result.(string)

	if !ok {
		return documentSymbols, fmt.Errorf("Error reading document symbols response: %v", result)
	}

	err = documentSymbols.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return documentSymbols, fmt.Errorf("Error unmarshalling document symbols response: %v", err)
	}

	return documentSymbols, nil
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

func New(opts ...optionProvider) (*Nvim, error) {
	options, err := NewOptions(opts...)

	if err != nil {
		return nil, fmt.Errorf("Error getting nvim options: %v", err)
	}

	arguments := []string{"--embed", "--headless", "-u", options.config.InitFile()}
	arguments = append(arguments, options.arguments...)
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
