package symbol

const TypeReferenceKind string = "typereference"

type TypeReference struct {
	Kind  string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

func (r *TypeReference) Canonical() Type {
	return r
}

// GetKind implements Type.
func (r *TypeReference) GetKind() string {
	return r.Kind
}

var _ Type = (*TypeReference)(nil)

func NewTypeReference(value string) *TypeReference {
	return &TypeReference{
		Kind:  TypeReferenceKind,
		Value: value,
	}
}

const ModuleReferenceKind string = "modulereference"

type ModuleReference struct {
	Kind  string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

func (r *ModuleReference) Canonical() Type {
	return r
}

// GetKind implements Module.
func (r *ModuleReference) GetKind() string {
	return r.Kind
}

var _ Type = (*ModuleReference)(nil)

func NewModuleReference(value string) *ModuleReference {
	return &ModuleReference{
		Kind:  ModuleReferenceKind,
		Value: value,
	}
}
