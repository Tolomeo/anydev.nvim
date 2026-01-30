package target

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

type TargetKind string

const (
	TargetKindValue TargetKind = "value"
	TargetKindType  TargetKind = "type"
)

type Target struct {
	kind          TargetKind
	parent        *Target
	name          string
	originChain   *origin.OriginChain
	meta          symbol.Meta
	documentation symbol.Documentation
	type_         annotation.Type
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

func (t *Target) OriginChain() *origin.OriginChain {
	return t.originChain
}

func (t *Target) SetOriginChain(o *origin.OriginChain) {
	t.originChain = o
}

func (t *Target) Origin() origin.Origin {
	return t.originChain.Last()
}

func (t *Target) Meta() symbol.Meta {
	return t.meta
}

func (t *Target) SetMeta(m symbol.Meta) {
	t.meta = m
}

func (t *Target) Type() annotation.Type {
	return t.type_
}

func (t *Target) SetType(typ annotation.Type) {
	t.type_ = typ
}

func (t *Target) Documentation() symbol.Documentation {
	return t.documentation
}

func (t *Target) SetDocumentation(documentation symbol.Documentation) {
	t.documentation = documentation
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
		kind:          kind,
		name:          name,
		meta:          symbol.Meta{},
		documentation: symbol.Documentation{},
	}
}
