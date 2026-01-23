package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type Origin struct {
	location nvim.Location
	node     treesitter.TsNode
	docBlock []string
}

func (l *Origin) Url() string {
	return l.location.Url
}

func (l *Origin) Line() uint {
	return l.location.StartLine()
}

func (l *Origin) Character() uint {
	return l.location.StartCharacter()
}

func (l *Origin) Type() string {
	return l.node.Type
}

func (l *Origin) Definition() []string {
	return strings.Split(l.node.Text, "\n")
}

func (l *Origin) Documentation() []string {
	return l.docBlock
}

func NewOrigin(location nvim.Location, definition treesitter.TsNode, documentation []string) *Origin {
	return &Origin{
		location: location,
		node:     definition,
		docBlock: documentation,
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
