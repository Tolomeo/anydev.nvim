package symbol

const arrayKind string = "array"

type Array struct {
	Kind  string  `json:"kind" yaml:"kind" mapstructure:"kind"`
	Items Type `json:"items" yaml:"items" mapstructure:"items"`
}

// Kind implements Type.
func (a *Array) GetKind() string {
	return a.Kind
}

var _ Type = (*Array)(nil)

type ArrayItems Symbol

func NewArray(items Type) *Array {
	return &Array{
		Kind:  arrayKind,
		Items: items,
	}
}
