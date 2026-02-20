 package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type VirtualOrigin struct {
	origin
	node definition.Virtual
}

func (mo *VirtualOrigin) Definition() []string {
	return strings.Split(mo.node.Root().Text, "\n")
}

func NewVirtualOrigin(location nvim.Location, node definition.Virtual, annotations []string) *VirtualOrigin {
	virtualOrigin := VirtualOrigin{
		origin: origin{
			location:    location,
			annotations: annotations,
		},
		node: node,
	}

	return &virtualOrigin
}
