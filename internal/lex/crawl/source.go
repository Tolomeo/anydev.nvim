package crawl

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type origin struct {
	location      nvim.Location
	node          treesitter.TsNode
	documentation []string
}

func (o *origin) Url() string {
	return o.location.Url
}

func (o *origin) Line() uint {
	return uint(o.location.TargetRange.Start.Line)
}

func (o *origin) Character() uint {
	return uint(o.location.TargetRange.Start.Character)
}

func (o *origin) Type() string {
	return o.node.Type
}

func (o *origin) Definition() []string {
	return strings.Split(o.node.Text, "\n")
}

func (o *origin) Documentation() []string {
	return o.documentation
}

type Source struct {
	path   string
	origin *origin
}

func (s *Source) Path() string {
	return s.path
}

func (s *Source) Origin() *origin {
	return s.origin
}
