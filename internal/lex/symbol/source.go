package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type Origin interface{}

type Source interface {
	Identifier() string
	Url() string
	Line() uint
	Character() uint
	Type() uint
	DefinitionLines() []string
	DocumentationLines() []string
}

type ValueOrigin struct {
	Location      nvim.Location
	Definition    treesitter.TsNode
	Documentation *treesitter.TsNode
}

type ValueSource struct {
	Path   string
	Origin *ValueOrigin
}

func (s *ValueSource) Identifier() string {
	return s.Path
}

func (s *ValueSource) Url() string {
	return s.Origin.Location.Url
}

func (s *ValueSource) Line() uint {
	return uint(s.Origin.Location.TargetRange.Start.Line)
}

func (s *ValueSource) Character() uint {
	return uint(s.Origin.Location.TargetRange.Start.Character)
}

func (s *ValueSource) Type() string {
	return s.Origin.Definition.Type
}

func (s *ValueSource) DefinitionLines() []string {
	/* if o.Definition == nil {
		return []string{}
	} */

	return strings.Split(s.Origin.Definition.Text, "\n")
}

func (s *ValueSource) DocumentationLines() []string {
	if s.Origin.Documentation == nil {
		return []string{}
	}

	return strings.Split(s.Origin.Documentation.Text, "\n")
}

type TypeOrigin struct {
	Location      nvim.Location
	Definition    treesitter.TsNode
	Documentation treesitter.TsNode
}

type TypeSource struct {
	Path   string
	Origin *TypeOrigin
}

func (s *TypeSource) Identifier() string {
	return s.Path
}

func (s *TypeSource) Url() string {
	return s.Origin.Location.Url
}

func (s *TypeSource) Line() uint {
	return uint(s.Origin.Location.TargetRange.Start.Line)
}

func (s *TypeSource) Character() uint {
	return uint(s.Origin.Location.TargetRange.Start.Character)
}

func (s *TypeSource) Type() string {
	return s.Origin.Definition.Type
}

func (s *TypeSource) DefinitionText() string {
	return s.Origin.Definition.Text
}

func (s *TypeSource) DefinitionLines() []string {
	/* if o.Definition == nil {
		return []string{}
	} */

	return strings.Split(s.DefinitionText(), "\n")
}

func (s *TypeSource) DocumentationText() string {
	return s.Origin.Documentation.Text
}

func (s *TypeSource) DocumentationLines() []string {
	return strings.Split(s.DocumentationText(), "\n")
}
