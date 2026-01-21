package symbol

type Origin interface {
	Url() string
	Line() uint
	Character() uint
	Type() string
	Definition() []string
	Documentation() []string
}

type TargetKind string

const (
	TargetKindValue TargetKind = "value"
	TargetKindType  TargetKind = "type"
)

type Target interface {
	Kind() TargetKind
	ParentName() string
	Identifier() string
	Name() string
	Origin() Origin
}
