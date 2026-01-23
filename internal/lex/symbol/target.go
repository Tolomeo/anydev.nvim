package symbol

import "fmt"

type TargetKind string

const (
	TargetKindValue TargetKind = "value"
	TargetKindType  TargetKind = "type"
)

type Target struct {
	kind    TargetKind
	parent  *Target
	name    string
	origins *Origins
	meta    Meta
	type_   Type
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

func (t *Target) Origins() *Origins {
	return t.origins
}

func (t *Target) SetOrigins(o *Origins) {
	t.origins = o
}

func (t *Target) Origin() Origin {
	return t.origins.Last()
}

func (t *Target) Meta() Meta {
	return t.meta
}

func (t *Target) SetMeta(m Meta) {
	t.meta = m
}

func (t *Target) Type() Type {
	return t.type_
}

func (t *Target) SetType(ty Type) {
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
