package symbol

const UnionKind string = "union"

type Union struct {
	Kind  string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Types []Type `json:"types" yaml:"types" mapstructure:"types"`
}

func (u *Union) GetKind() string {
	return u.Kind
}

var _ Type = (*Union)(nil)

func NewUnion(types []Type) *Union {
	return &Union{
		Kind:  UnionKind,
		Types: types,
	}
}
