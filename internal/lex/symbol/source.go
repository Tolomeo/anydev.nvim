package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type Origin interface {
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
	Identifier() string
	GetOrigin() Origin
	SetOrigin(Origin)
}

type ValueOrigin struct {
	Location      nvim.Location
	Definition    treesitter.TsNode
	Documentation treesitter.TsNode
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

func (o *ValueOrigin) DefinitionText() string {
	return o.Definition.Text
}

func (o *ValueOrigin) DefinitionLines() []string {
	return strings.Split(o.Definition.Text, "\n")
}

func (o *ValueOrigin) DocumentationText() string {
	return o.Documentation.Text
}

func (o *ValueOrigin) DocumentationLines() []string {
	return strings.Split(o.Documentation.Text, "\n")
}

func (o *ValueOrigin) SetDocumentation(documentation treesitter.TsNode) {
	o.Documentation = documentation
}

type ValueSource struct {
	Path   string
	Origin Origin
}

func (s *ValueSource) Identifier() string {
	return s.Path
}

func (s *ValueSource) GetOrigin() Origin {
	return s.Origin
}

func (s *ValueSource) SetOrigin(origin Origin) {
	s.Origin = origin
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
	Path   string
	Origin Origin
}

func (s *TypeSource) Identifier() string {
	return s.Path
}

func (s *TypeSource) GetOrigin() Origin {
	return s.Origin
}

func (s *TypeSource) SetOrigin(origin Origin) {
	s.Origin = origin
}
