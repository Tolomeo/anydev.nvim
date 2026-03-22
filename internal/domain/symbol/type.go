package symbol

type Type interface {
	GetKind() string
	Canonical() Type
}
