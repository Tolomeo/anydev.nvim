package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type TableOrigin struct {
	origin
	node definition.Table
}

func (to *TableOrigin) Definition() []string {
	return strings.Split(to.node.Root().Text, "\n")
}

func (to *TableOrigin) Name() string {
	return to.node.Name().Text
}

func NewTableOrigin(location nvim.Location, node definition.Table, annotations []string) *TableOrigin {
	tableOrigin := TableOrigin{
		origin: origin{
			location:    location,
			annotations: annotations,
		},
		node: node,
	}

	return &tableOrigin
}
