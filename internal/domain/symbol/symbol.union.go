package symbol

import type_ "github.com/Tolomeo/anydev.nvim/internal/domain/type"

const UnionKind string = "union"

type Union struct {
	Kind  string       `json:"kind" yaml:"kind" mapstructure:"kind"`
	Types []type_.Type `json:"types" yaml:"types" mapstructure:"types"`
}

func (u *Union) GetKind() string {
	return u.Kind
}

var _ type_.Type = (*Union)(nil)

func NewUnion(types []type_.Type) *Union {
	return &Union{
		Kind:  UnionKind,
		Types: types,
	}
}
