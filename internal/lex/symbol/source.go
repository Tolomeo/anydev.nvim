package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type ValueOrigin struct {
	Location      nvim.Location
	Definition    treesitter.TsNode
	Documentation *treesitter.TsNode
}

func (o *ValueOrigin) Url() string {
	return o.Location.Url
}

func (o *ValueOrigin) Line() uint {
	return uint(o.Location.TargetRange.Start.Line)
}

func (o *ValueOrigin) Character() uint {
	return uint(o.Location.TargetRange.Start.Character)
}

func (o *ValueOrigin) Type() string {
	return o.Definition.Type
}

func (o *ValueOrigin) DefinitionLines() []string {
	/* if o.Definition == nil {
		return []string{}
	} */

	return strings.Split(o.Definition.Text, "\n")
}

func (o *ValueOrigin) DocumentationLines() []string {
	return strings.Split(o.Documentation.Text, "\n")
}

type ValueSource struct {
	Path   string
	Origin *ValueOrigin
}

type TypeOrigin struct {
	Location      nvim.Location
	Definition    treesitter.TsNode
	Documentation treesitter.TsNode
}

func (o *TypeOrigin) Url() string {
	return o.Location.Url
}

func (o *TypeOrigin) Line() uint {
	return uint(o.Location.TargetRange.Start.Line)
}

func (o *TypeOrigin) Character() uint {
	return uint(o.Location.TargetRange.Start.Character)
}

func (o *TypeOrigin) Type() string {
	return o.Definition.Type
}

func (o *TypeOrigin) DefinitionText() string {
	return o.Definition.Text
}

func (o *TypeOrigin) DefinitionLines() []string {
	/* if o.Definition == nil {
		return []string{}
	} */

	return strings.Split(o.DefinitionText(), "\n")
}

func (o *TypeOrigin) DocumentationText() string {
	return o.Documentation.Text
}

func (o *TypeOrigin) DocumentationLines() []string {
	return strings.Split(o.DocumentationText(), "\n")
}

type TypeSource struct {
	Path   string
	Origin *TypeOrigin
}
