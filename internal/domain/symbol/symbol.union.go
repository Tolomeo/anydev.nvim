package symbol

import "github.com/Tolomeo/anydev.nvim/internal/domain/annotation"

const UnionKind string = "union"

type Union struct {
	Kind  string       `json:"kind" yaml:"kind" mapstructure:"kind"`
	Types []annotation.Type `json:"types" yaml:"types" mapstructure:"types"`
}

func (u *Union) GetKind() string {
	return u.Kind
}

var _ annotation.Type = (*Union)(nil)

func NewUnion(types []annotation.Type) *Union {
	return &Union{
		Kind:  UnionKind,
		Types: types,
	}
}
