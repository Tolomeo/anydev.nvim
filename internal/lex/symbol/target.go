package symbol

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
	Origin() *Origin
}
