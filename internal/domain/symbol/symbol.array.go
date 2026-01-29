package symbol

import "github.com/Tolomeo/anydev.nvim/internal/domain/annotation"

const arrayKind string = "array"

type Array struct {
	Kind  string  `json:"kind" yaml:"kind" mapstructure:"kind"`
	Items annotation.Type `json:"items" yaml:"items" mapstructure:"items"`
}

// Kind implements Type.
func (a *Array) GetKind() string {
	return a.Kind
}

var _ annotation.Type = (*Array)(nil)

type ArrayItems Symbol

func NewArray(items annotation.Type) *Array {
	return &Array{
		Kind:  arrayKind,
		Items: items,
	}
}
