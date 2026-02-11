package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type ValueOrigin struct {
	origin
	node definition.Value
}

func (vo *ValueOrigin) Definition() []string {
	return strings.Split(vo.node.Root().Text, "\n")
}

func (vo *ValueOrigin) Value() string {
	return vo.node.Value().Text
}

func (vo *ValueOrigin) Type() string {
	return vo.node.Type()
}

func NewValueOrigin(location nvim.Location, node definition.Value, annotations []string) *ValueOrigin {
	valueOrigin := ValueOrigin{
		origin: origin{
			location:    location,
			annotations: annotations,
		},
		node: node,
	}

	return &valueOrigin

}
