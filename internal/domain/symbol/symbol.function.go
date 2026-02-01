package symbol

import "github.com/Tolomeo/anydev.nvim/internal/domain/annotation"

var (
	FunctionClassAccess   string = "class"
	FunctionIstanceAccess string = "instance"
)

const functionKind string = "function"

type Function struct {
	Kind      string             `json:"kind" yaml:"kind" mapstructure:"kind"`
	Access    *string            `json:"access,omitempty" yaml:"access,omitempty" mapstructure:"access,omitempty"`
	Arguments []FunctionArgument `json:"arguments" yaml:"arguments" mapstructure:"arguments"`
	Generics  []FunctionGeneric  `json:"generics" yaml:"generics" mapstructure:"generics"`
	Name      string             `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Overloads []FunctionOverload `json:"overloads,omitempty" yaml:"overloads,omitempty" mapstructure:"overloads,omitempty"`
	Returns   []FunctionReturn   `json:"returns" yaml:"returns" mapstructure:"returns"`
}

// GetKind implements Type.
func (f *Function) GetKind() string {
	return f.Kind
}

var _ annotation.Type = (*Function)(nil)

type FunctionArgument struct {
	Documentation Documentation   `json:"documentation,omitempty" yaml:"documentation,omitempty" mapstructure:"documentation,omitempty"`
	Name          string          `json:"name" yaml:"name" mapstructure:"name"`
	Optional      bool            `json:"optional" yaml:"optional" mapstructure:"optional"`
	Type          annotation.Type `json:"type" yaml:"type" mapstructure:"type"`
}

type FunctionArgumentType Symbol

type FunctionGeneric struct {
	Name string          `json:"name" yaml:"name" mapstructure:"name"`
	Type annotation.Type `json:"type" yaml:"type" mapstructure:"type"`
}

type FunctionGenericTypesElem Symbol

type FunctionOverload struct {
	Arguments     []FunctionArgument `json:"arguments" yaml:"arguments" mapstructure:"arguments"`
	Documentation Documentation      `json:"documentation" yaml:"documentation" mapstructure:"documentation"`
	Generics      []FunctionGeneric  `json:"generics" yaml:"generics" mapstructure:"generics"`
	Returns       []FunctionReturn   `json:"returns" yaml:"returns" mapstructure:"returns"`
}

type FunctionReturn struct {
	Documentation Documentation   `json:"documentation,omitempty" yaml:"documentation,omitempty" mapstructure:"documentation,omitempty"`
	Name          string          `json:"name" yaml:"name" mapstructure:"name"`
	Type          annotation.Type `json:"type" yaml:"type" mapstructure:"type"`
}

type FunctionReturnType Symbol

func NewFunction() *Function {
	return &Function{
		Kind:      functionKind,
		Generics:  []FunctionGeneric{},
		Arguments: []FunctionArgument{},
		Returns:   []FunctionReturn{},
		Overloads: []FunctionOverload{},
	}
}

func NewFunctionGeneric(name string, typ annotation.Type) *FunctionGeneric {
	return &FunctionGeneric{
		Name:  name,
		Type: typ,
	}
}

func NewFunctionArgument(name string) *FunctionArgument {
	return &FunctionArgument{
		Name:     name,
		Type:     NewUnknown(),
		Optional: false,
	}
}

func NewFunctionReturn() *FunctionReturn {
	return &FunctionReturn{}
}

func NewFunctionOverload() *FunctionOverload {
	return &FunctionOverload{}
}
