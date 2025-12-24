package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type Origin struct {
	Location      nvim.Location
	Node          treesitter.TsNode
	Documentation []string
}

func (o *Origin) Url() string {
	return o.Location.Url
}

func (o *Origin) Line() uint {
	return uint(o.Location.TargetRange.Start.Line)
}

func (o *Origin) Character() uint {
	return uint(o.Location.TargetRange.Start.Character)
}

func (o *Origin) Type() string {
	return o.Node.Type
}

func (o *Origin) Definition() []string {
	return strings.Split(o.Node.Text, "\n")
}

type Source struct {
	Path   string
	Origin *Origin
}
