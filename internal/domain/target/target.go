package target

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

type TargetKind string

const (
	TargetKindValue TargetKind = "value"
	TargetKindType  TargetKind = "type"
)

type Target struct {
	kind    TargetKind
	parent  *Target
	name    string
	origins *origin.Origins
	meta    symbol.Meta
	type_   symbol.Type
}

func (t *Target) Kind() TargetKind {
	return t.kind
}

func (t *Target) ParentName() string {
	if t.parent == nil {
		return ""
	}

	return t.parent.Name()
}

func (t *Target) Name() string {
	return t.name
}

func (t *Target) Identifier() string {
	if t.parent == nil {
		return t.name
	}

	return fmt.Sprintf("%s.%s", t.parent.Identifier(), t.name)
}

func (t *Target) Origins() *origin.Origins {
	return t.origins
}

func (t *Target) SetOrigins(o *origin.Origins) {
	t.origins = o
}

func (t *Target) Origin() origin.Origin {
	return t.origins.Last()
}

func (t *Target) Meta() symbol.Meta {
	return t.meta
}

func (t *Target) SetMeta(m symbol.Meta) {
	t.meta = m
}

func (t *Target) Type() symbol.Type {
	return t.type_
}

func (t *Target) SetType(ty symbol.Type) {
	t.type_ = ty
}

func (t *Target) NewChild(name string) *Target {
	return &Target{
		parent: t,
		kind:   t.kind,
		name:   name,
	}
}

func NewTarget(kind TargetKind, name string) *Target {
	return &Target{
		kind: kind,
		name: name,
	}
}
