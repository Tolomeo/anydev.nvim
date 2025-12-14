package crawl

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/ts"
)

type origin struct {
	location      nvim.Location
	node          ts.TsNode
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

func (o *origin) Definition() string {
	return o.node.Text
}

func (o *origin) Documentation() []string {
	return o.documentation
}

type Source struct {
	path   string
	origin *origin
	fields []*Source
}

func (s *Source) Path() string {
	return s.path
}

func (s *Source) Origin() *origin {
	return s.origin
}

func (s *Source) Fields() []*Source {
	return s.fields
}
