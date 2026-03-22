package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type ClassOrigin struct {
	origin
	node annotation.AtClass
}

func (co *ClassOrigin) Definition() []string {
	return strings.Split(co.node.Root().Text, "\n")
}

func (co *ClassOrigin) Name() string {
	return co.node.Name().Text
}

func NewClassOrigin(location nvim.Location, node annotation.AtClass) *ClassOrigin {
	classOrigin := ClassOrigin{
		origin: origin{
			location: location,
		},
		node: node,
	}

	return &classOrigin
}
