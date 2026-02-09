package symbol

const EnumeratorKind = "enumerator"

type Enumerator struct {
	Fields []EnumeratorField `json:"fields" yaml:"fields" mapstructure:"fields"`
	Kind   string            `json:"kind" yaml:"kind" mapstructure:"kind"`
	Name   string            `json:"name" yaml:"name" mapstructure:"name"`
}

func (e *Enumerator) GetKind() string {
	return e.Kind
}

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
