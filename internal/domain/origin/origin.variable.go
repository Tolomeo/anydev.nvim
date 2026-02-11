package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type VariableOrigin struct {
	origin
	node definition.Variable
}

func (vo *VariableOrigin) Definition() []string {
	return strings.Split(vo.node.Root().Text, "\n")
}

func (vo *VariableOrigin) Name() string {
	return vo.node.Name().Text
}

func (vo *VariableOrigin) NameRange() treesitter.Range {
	return vo.node.Name().Range
}

func NewVariableOrigin(location nvim.Location, node definition.Variable, annotations []string) *VariableOrigin {
	variableOrigin := VariableOrigin{
		origin: origin{
			location:    location,
			annotations: annotations,
		},
		node: node,
	}

	return &variableOrigin
}
