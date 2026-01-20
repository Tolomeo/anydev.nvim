package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

/* type Origin interface {
	Url() string
	Line() uint
	Character() uint
	Type() string
	DefinitionText() string
	DefinitionLines() []string
	DocumentationText() string
	DocumentationLines() []string
	SetDocumentation(treesitter.TsNode)
}

type Source interface {
	Name() string
	Identifier() string
	GetOrigin() Origin
	SetOrigin(Origin)
} */

type Origin struct {
	Location      nvim.Location
	Definition    treesitter.TsNode
	Documentation treesitter.TsNode
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

func (o *Origin) DefinitionText() string {
	return o.Definition.Text
}

func (o *Origin) DefinitionLines() []string {
	return strings.Split(o.Definition.Text, "\n")
}

func (o *Origin) DocumentationText() string {
	return o.Documentation.Text
}

func (o *Origin) DocumentationLines() []string {
	return strings.Split(o.Documentation.Text, "\n")
}

func (o *Origin) SetDocumentation(documentation treesitter.TsNode) {
	o.Documentation = documentation
}

/* type ValueSource struct {
	name   string
	origin Origin
}

func (s *ValueSource) Name() string {
	return s.name
}

func (s *ValueSource) Identifier() string {
	return s.name
}

func (s *ValueSource) GetOrigin() Origin {
	return s.origin
}

func (s *ValueSource) SetOrigin(origin Origin) {
	s.origin = origin
}

func NewValueSource(name string) *ValueSource {
	return &ValueSource{
		name: name,
	}
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
	return strings.Split(o.DefinitionText(), "\n")
}

func (o *TypeOrigin) DocumentationText() string {
	return o.Documentation.Text
}

func (o *TypeOrigin) DocumentationLines() []string {
	return strings.Split(o.DocumentationText(), "\n")
}

func (o *TypeOrigin) SetDocumentation(documentation treesitter.TsNode) {
	o.Documentation = documentation
}

type TypeSource struct {
	parent string
	name   string
	origin Origin
}

func (s *TypeSource) ParentName() string {
	return s.parent
}

func (s *TypeSource) Name() string {
	return s.name
}

func (s *TypeSource) Identifier() string {
	if s.parent == "" {
		return s.name
	}
	return strings.Join([]string{s.parent, s.name}, ".")
}

func (s *TypeSource) GetOrigin() Origin {
	return s.origin
}

func (s *TypeSource) SetOrigin(origin Origin) {
	s.origin = origin
}

func NewTypeSource(name string, parent string) *TypeSource {
	return &TypeSource{
		parent: parent,
		name:   name,
	}
} */
