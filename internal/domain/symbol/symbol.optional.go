package symbol

import "github.com/Tolomeo/anydev.nvim/internal/domain/annotation"

const OptionalKind string = "optional"

type Optional struct {
	Kind string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Type annotation.Type   `json:"type" yaml:"type" mapstructure:"type"`
}

func (o *Optional) GetKind() string {
	return o.Kind
}

var _ annotation.Type = (*Optional)(nil)

func NewOptional(typ annotation.Type) *Optional {
	return &Optional{
		Kind: OptionalKind,
		Type: typ,
	}
}
