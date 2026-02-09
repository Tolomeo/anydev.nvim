package symbol

const OptionalKind string = "optional"

type Optional struct {
	Kind string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Type Type   `json:"type" yaml:"type" mapstructure:"type"`
}

func (o *Optional) GetKind() string {
	return o.Kind
}

var _ Type = (*Optional)(nil)

func NewOptional(typ Type) *Optional {
	return &Optional{
		Kind: OptionalKind,
		Type: typ,
	}
}
