package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type FunctionOrigin struct {
	origin
	node definition.Function
}

func (fo *FunctionOrigin) Definition() []string {
	return strings.Split(fo.node.Root().Text, "\n")
}

func (fo *FunctionOrigin) Name() string {
	return fo.node.Name().Text
}

func (fo *FunctionOrigin) Static() bool {
	return fo.node.Static() != nil
}

func (fo *FunctionOrigin) Args() []string {
	args, _ := slicesx.MapFunc(fo.node.Args(), func(argNode treesitter.TsNode) (string, error) {
		return argNode.Text, nil
	})

	return args
}

func NewFunctionOrigin(location nvim.Location, node definition.Function, annotations []string) *FunctionOrigin {
	functionOrigin := FunctionOrigin{
		origin: origin{
			location:    location,
			annotations: annotations,
		},
		node: node,
	}

	return &functionOrigin
}
