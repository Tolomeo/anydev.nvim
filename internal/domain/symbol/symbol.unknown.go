package symbol

import type_ "github.com/Tolomeo/anydev.nvim/internal/domain/type"

const UnknownKindUnknown string = "unknown"

type Unknown struct {
	Kind          string        `json:"kind" yaml:"kind" mapstructure:"kind"`
	Documentation Documentation `json:"documentation" yaml:"documentation" mapstructure:"documentation"`
}

// GetKind implements Type.
func (u *Unknown) GetKind() string {
	return u.Kind
}

var _ type_.Type = (*Unknown)(nil)

func NewUnknown() *Unknown {
	return &Unknown{
		Kind: UnknownKindUnknown,
	}
}
