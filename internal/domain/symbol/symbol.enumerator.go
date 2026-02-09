package symbol

const EnumeratorKind = "enumerator"

type Enumerator struct {
	Kind   string            `json:"kind" yaml:"kind" mapstructure:"kind"`
	Name   string            `json:"name" yaml:"name" mapstructure:"name"`
	Fields []EnumeratorField `json:"fields" yaml:"fields" mapstructure:"fields"`
}

func (e *Enumerator) Canonical() Type {
	return e
}

func (e *Enumerator) GetKind() string {
	return e.Kind
}

var _ Type = (*Enumerator)(nil)

type EnumeratorField struct {
	Name  string `json:"name" yaml:"name" mapstructure:"name"`
	Value Type   `json:"value" yaml:"value" mapstructure:"value"`
}

func NewEnumerator(name string) *Enumerator {
	return &Enumerator{
		Kind:   EnumeratorKind,
		Name:   name,
		Fields: []EnumeratorField{},
	}
}
