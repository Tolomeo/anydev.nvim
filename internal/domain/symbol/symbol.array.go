package symbol

import type_ "github.com/Tolomeo/anydev.nvim/internal/domain/type"

const arrayKind string = "array"

type Array struct {
	Kind  string  `json:"kind" yaml:"kind" mapstructure:"kind"`
	Items type_.Type `json:"items" yaml:"items" mapstructure:"items"`
}

// Kind implements Type.
func (a *Array) GetKind() string {
	return a.Kind
}

var _ type_.Type = (*Array)(nil)

type ArrayItems Symbol

func NewArray(items type_.Type) *Array {
	return &Array{
		Kind:  arrayKind,
		Items: items,
	}
}
