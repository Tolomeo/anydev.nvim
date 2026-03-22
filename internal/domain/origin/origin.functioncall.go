package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type FunctionCallOrigin struct {
	origin
	node definition.FunctionCall
}

func (fco *FunctionCallOrigin) Definition() []string {
	return strings.Split(fco.node.Root().Text, "\n")
}

func (fco *FunctionCallOrigin) Name() string {
	return fco.node.Name().Text
}

func (fco *FunctionCallOrigin) FunctionName() string {
	return fco.node.FunctionName().Text
}

func (fco *FunctionCallOrigin) FunctionArguments() []string {
	args := []string{}

	for _, node := range fco.node.FunctionArguments() {
		args = append(args, node.Text)
	}

	return args
}

func (fco *FunctionCallOrigin) FunctionNameRange() treesitter.Range {
	return fco.node.FunctionName().Range
}

func NewFunctionCallOrigin(location nvim.Location, node definition.FunctionCall, annotations []string) *FunctionCallOrigin {
	functionCallOrigin := FunctionCallOrigin{
		origin: origin{
			location:    location,
			annotations: annotations,
		},
		node: node,
	}

	return &functionCallOrigin
}
