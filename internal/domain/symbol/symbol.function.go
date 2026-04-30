package symbol

const functionKind string = "function"

type Function struct {
	Kind      string             `json:"kind" yaml:"kind" mapstructure:"kind"`
	Arguments []FunctionArgument `json:"arguments" yaml:"arguments" mapstructure:"arguments"`
	Generics  []FunctionGeneric  `json:"generics" yaml:"generics" mapstructure:"generics"`
	Name      string             `json:"name" yaml:"name" mapstructure:"name"`
	Overloads []FunctionOverload `json:"overloads" yaml:"overloads" mapstructure:"overloads"`
	Returns   []FunctionReturn   `json:"returns" yaml:"returns" mapstructure:"returns"`
}

func (f *Function) Canonical() Type {
	return f
}

func (f *Function) GetKind() string {
	return f.Kind
}

var _ Type = (*Function)(nil)

type FunctionArgument struct {
	Documentation Documentation `json:"documentation" yaml:"documentation" mapstructure:"documentation"`
	Name          string        `json:"name" yaml:"name" mapstructure:"name"`
	Optional      bool          `json:"optional" yaml:"optional" mapstructure:"optional"`
	Type          Type          `json:"type" yaml:"type" mapstructure:"type"`
}

type FunctionArgumentType Symbol

type FunctionGeneric struct {
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	Type Type   `json:"type" yaml:"type" mapstructure:"type"`
}

type FunctionGenericTypesElem Symbol

type FunctionOverload struct {
	Arguments     []FunctionArgument `json:"arguments" yaml:"arguments" mapstructure:"arguments"`
	Documentation Documentation      `json:"documentation" yaml:"documentation" mapstructure:"documentation"`
	Generics      []FunctionGeneric  `json:"generics" yaml:"generics" mapstructure:"generics"`
	Returns       []FunctionReturn   `json:"returns" yaml:"returns" mapstructure:"returns"`
}

type FunctionReturn struct {
	Documentation Documentation `json:"documentation" yaml:"documentation" mapstructure:"documentation"`
	Name          string        `json:"name" yaml:"name" mapstructure:"name"`
	Type          Type          `json:"type" yaml:"type" mapstructure:"type"`
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

func NewFunctionGeneric(name string, typ Type) *FunctionGeneric {
	return &FunctionGeneric{
		Name: name,
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
