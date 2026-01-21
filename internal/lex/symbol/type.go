package symbol

type Type interface {
	GetKind() string
}

type WithType interface {
	Type() Type
}

type withType struct {
	type_ Type
}

func (w withType) Type() Type {
	return w.type_
}

var _ WithType = (*withType)(nil)
