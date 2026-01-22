package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type Origin struct {
	location      nvim.Location
	definition    treesitter.TsNode
	documentation treesitter.TsNode
}

func (l *Origin) Url() string {
	return l.location.Url
}

func (l *Origin) Line() uint {
	return uint(l.location.TargetRange.Start.Line)
}

func (l *Origin) Character() uint {
	return uint(l.location.TargetRange.Start.Character)
}

func (l *Origin) Type() string {
	return l.definition.Type
}

func (l *Origin) Definition() []string {
	return strings.Split(l.definition.Text, "\n")
}

func (l *Origin) Documentation() []string {
	return strings.Split(l.documentation.Text, "\n")
}

func (l *Origin) SetDocumentation(d treesitter.TsNode) {
	l.documentation = d
}

func NewOrigin(location nvim.Location, definition treesitter.TsNode) *Origin {
	return &Origin{
		location:   location,
		definition: definition,
	}
}

type Origins []*Origin

func (o *Origins) Last() *Origin {
	return (*o)[len(*o)-1]
}

func (o *Origins) First() *Origin {
	return (*o)[0]
}

func (o *Origins) Merge(o2 *Origins) {
	*o = append(*o, *o2...)
}

func NewOrigins(origins ...*Origin) *Origins {
	t := Origins{}

	for _, l := range origins {
		t = append(t, l)
	}

	return &t
}
