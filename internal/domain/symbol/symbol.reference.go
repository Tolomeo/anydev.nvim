package symbol

import "github.com/Tolomeo/anydev.nvim/internal/domain/annotation"

const ReferenceKind string = "reference"

type Reference struct {
	Kind  string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

// GetKind implements Type.
func (r *Reference) GetKind() string {
	return r.Kind
}

var _ annotation.Type = (*Reference)(nil)

func NewReference(value string) *Reference {
	return &Reference{
		Kind:  ReferenceKind,
		Value: value,
	}
}
