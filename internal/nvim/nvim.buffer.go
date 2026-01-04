package nvim

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type Buffer struct {
	name                       string
	ReadLines                  func() ([]string, error)
	SetLines                   func([]string) error
	Close                      func() error
	GetTSNodeAt                func([]string, uint, uint) (*treesitter.TsNode, error)
	GetTsCommentBlockAt        func(uint, uint) (*treesitter.TsNode, error)
	TsQueryOne                 func(treesitter.Query) (*TsQueryMatch, error)
	TsQueryAll                 func(treesitter.Query) (*[]TsQueryMatch, error)
	SafeTsQueryAll             func(treesitter.Query) (*[]SafeTsQueryResult, error)
	GetDefinitionLocations     func(uint, uint) (*[]Location, error)
	GetTypeCompletion          func(uint, uint) ([]string, error)
	GetTypeDefinitionLocations func(uint, uint) (*[]TypeDefinitionLocation, error)
}

func (n *Nvim) OpenBuffer(name string) (*Buffer, error) {
	buffer := &Buffer{
		name: name,
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
		GetTypeCompletion: func(line, character uint) ([]string, error) {
			_, err := n.open(name)

			if err != nil {
				return nil, err
			}

			return n.getTypeCompletion(line, character)
		},
		GetTypeDefinitionLocations: func(line uint, character uint) (*[]TypeDefinitionLocation, error) {
			_, err := n.open(name)

			if err != nil {
				return nil, err
			}

			return n.getTypeDefinitionLocations(line, character)
		},
		Close: func() error {
			_, err := n.open(name)

			if err != nil {
				return err
			}

			err = n.deleteBuffer()

			if err != nil {
				return err
			}

			return nil
		},
	}

	_, err := n.open(name)

	if err != nil {
		return nil, err
	}

	return buffer, nil
}

var bufferCounter = 0

func (n *Nvim) NewBuffer() (*Buffer, error) {
	bufferCounter += 1
	name := path.Join(n.Options().Config().Dir(), fmt.Sprintf("anydev.%d.lua", bufferCounter))

	return n.OpenBuffer(name)
}
