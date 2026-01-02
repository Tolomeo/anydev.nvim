package nvim

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/msgpackrpc"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/anyx"
)

type Buffer struct {
	name                   string
	previous               string
	ReadLines              func() ([]string, error)
	SetLines               func([]string) error
	Close                  func() error
	GetTSNodeAt            func([]string, uint, uint) (*treesitter.TsNode, error)
	GetTsCommentBlockAt    func(uint, uint) (*treesitter.TsNode, error)
	TsQueryOne             func(treesitter.Query) (*TsQueryMatch, error)
	TsQueryAll             func(treesitter.Query) (*[]TsQueryMatch, error)
	SafeTsQueryAll         func(treesitter.Query) (*[]SafeTsQueryResult, error)
	GetDefinitionLocations func(uint, uint) (*[]Location, error)
}

func (n *Nvim) OpenBuffer(name string) (*Buffer, error) {
	fmt.Println("Opening buffer", name)

	previous, err := n.getBufferName()

	if err != nil {
		return nil, err
	}

	buf := &Buffer{
		name:     name,
		previous: previous,
		ReadLines: func() ([]string, error) {
			_, err := n.open(name)

			if err != nil {
				return []string{}, fmt.Errorf("Error reading buffer '%s': %w", name, err)
			}

			lines, err := n.getBufferLines(0, -1)

			if err != nil {
				return []string{}, fmt.Errorf("Error reading buffer '%s': %w", name, err)
			}

			return lines, nil
		},
		SetLines: func(lines []string) error {
			_, err := n.open(name)

			if err != nil {
				return fmt.Errorf("Error writing to buffer '%s': %w", name, err)
			}

			if len(lines) < 1 {
				return n.setBufferLines([]string{""})
			}

			return n.setBufferLines(lines)
		},
		GetTSNodeAt: func(nodeTypes []string, line uint, character uint) (*treesitter.TsNode, error) {
			_, err := n.open(name)

			if err != nil {
				return nil, err
			}

			return n.getTSNodeAt(nodeTypes, line, character)
		},
		GetTsCommentBlockAt: func(line uint, character uint) (*treesitter.TsNode, error) {
			_, err := n.open(name)

			if err != nil {
				return nil, err
			}

			return n.getTsCommentBlockAt(line, character)
		},
		TsQueryOne: func(query treesitter.Query) (*TsQueryMatch, error) {
			_, err := n.open(name)

			if err != nil {
				return nil, err
			}

			return n.tsQueryOne(query)
		},
		TsQueryAll: func(query treesitter.Query) (*[]TsQueryMatch, error) {
			_, err := n.open(name)

			if err != nil {
				return nil, err
			}

			return n.tsQueryAll(query)
		},
		SafeTsQueryAll: func(query treesitter.Query) (*[]SafeTsQueryResult, error) {
			_, err := n.open(name)

			if err != nil {
				return nil, err
			}

			return n.safeTsQueryAll(query)
		},
		GetDefinitionLocations: func(line uint, character uint) (*[]Location, error) {
			_, err := n.open(name)

			if err != nil {
				return nil, err
			}

			return n.getDefinitionLocations(line, character)
		},
		Close: func() error {
			fmt.Println("Closing buffer", name)
			_, err := n.open(name)

			if err != nil {
				return err
			}

			err = n.deleteBuffer()

			if err != nil {
				return err
			}

			_, err = n.open(previous)

			if err != nil {
				return err
			}

			return nil
		},
	}

	_, err = n.open(name)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (n *Nvim) NewBuffer() (*Buffer, error) {
	counter += 1
	name := path.Join(n.Options().Config().Dir(), fmt.Sprintf("anydev.%d.lua", counter))

	return n.OpenBuffer(name)
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

func (n *Nvim) write() error {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_command",
		Params: []any{"write"},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to write buffer: %v\n", err)
	}

	return nil
}

func (n *Nvim) getBufferName() (string, error) {
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
}

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

func (n *Nvim) getBufferText(startRow int, startCol int, endRow int, endCol int) ([]string, error) {
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
}

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
