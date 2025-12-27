package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type Origin struct {
	Location      nvim.Location
	Definition    treesitter.TsNode
	Documentation *treesitter.TsNode
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
	return o.Definition.Type
}

func (o *Origin) DefinitionLines() []string {
	/* if o.Definition == nil {
		return []string{}
	} */

	return strings.Split(o.Definition.Text, "\n")
}

func (o *Origin) DocumentationLines() []string {
	return strings.Split(o.Documentation.Text, "\n")
}

type Source struct {
	Path   string
	Origin *Origin
}
