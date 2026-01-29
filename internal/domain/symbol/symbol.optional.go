package symbol

import type_ "github.com/Tolomeo/anydev.nvim/internal/domain/type"

const OptionalKind string = "optional"

type Optional struct {
	Kind string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Type type_.Type   `json:"type" yaml:"type" mapstructure:"type"`
}

func (o *Optional) GetKind() string {
	return o.Kind
}

var _ type_.Type = (*Optional)(nil)

func NewOptional(typ type_.Type) *Optional {
	return &Optional{
		Kind: OptionalKind,
		Type: typ,
	}
}
