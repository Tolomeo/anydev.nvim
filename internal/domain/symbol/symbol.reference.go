package symbol

import type_ "github.com/Tolomeo/anydev.nvim/internal/domain/type"

const ReferenceKind string = "reference"

type Reference struct {
	Kind  string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

// GetKind implements Type.
func (r *Reference) GetKind() string {
	return r.Kind
}

var _ type_.Type = (*Reference)(nil)

func NewReference(value string) *Reference {
	return &Reference{
		Kind:  ReferenceKind,
		Value: value,
	}
}
