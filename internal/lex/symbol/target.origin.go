package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type Origin interface {
	Url() string
	Line() uint
	Character() uint
	Type() string
	Definition() []string
	Documentation() []string
}

type origin struct {
	location   nvim.Location
	definition nvim.TsNodeQueryMatch
	docBlock   []string
}

func (l *origin) Url() string {
	return l.location.Url
}

func (l *origin) Line() uint {
	return l.location.StartLine()
}

func (l *origin) Character() uint {
	return l.location.StartCharacter()
}

func (l *origin) Type() string {
	return l.definition.Node.Type
}

func (l *origin) Definition() []string {
	return strings.Split(l.definition.Node.Text, "\n")
}

func (l *origin) Documentation() []string {
	return l.docBlock
}

type FunctionOrigin struct {
	origin
}

func NewFunctionOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *FunctionOrigin {
	return &FunctionOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type TableOrigin struct {
	origin
}

func NewTableOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *TableOrigin {
	return &TableOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type VariableOrigin struct {
	origin
}

func NewVariableOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *VariableOrigin {
	return &VariableOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type ModuleOrigin struct {
	origin
}

func NewModuleOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *ModuleOrigin {
	return &ModuleOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type ClassOrigin struct {
	origin
}

func NewClassOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *ClassOrigin {
	return &ClassOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type AliasOrigin struct {
	origin
}

func NewAliasOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *AliasOrigin {
	return &AliasOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type FieldOrigin struct {
	origin
}

func NewFieldOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *FieldOrigin {
	return &FieldOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type MetaOrigin struct {
	origin
}

func NewMetaOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *MetaOrigin {
	return &MetaOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type Origins []Origin

func (o *Origins) Last() Origin {
	return (*o)[len(*o)-1]
}

func (o *Origins) First() Origin {
	return (*o)[0]
}

func (o *Origins) Merge(o2 *Origins) {
	*o = append(*o, *o2...)
}

func NewOrigins(origins ...Origin) *Origins {
	t := Origins{}

	for _, l := range origins {
		t = append(t, l)
	}

	return &t
}
