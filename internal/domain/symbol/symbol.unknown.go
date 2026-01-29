package symbol

import "github.com/Tolomeo/anydev.nvim/internal/domain/annotation"

const UnknownKindUnknown string = "unknown"

type Unknown struct {
	Kind          string        `json:"kind" yaml:"kind" mapstructure:"kind"`
	Documentation Documentation `json:"documentation" yaml:"documentation" mapstructure:"documentation"`
}

// GetKind implements Type.
func (u *Unknown) GetKind() string {
	return u.Kind
}

var _ annotation.Type = (*Unknown)(nil)

func NewUnknown() *Unknown {
	return &Unknown{
		Kind: UnknownKindUnknown,
	}
}
