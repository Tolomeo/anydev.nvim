package symbol

var (
	FunctionClassAccess   string = "class"
	FunctionIstanceAccess string = "instance"
)

func NewFunction() *Function {
	return &Function{
		Kind: FunctionKindFunction,
	}
}

func NewFunctionArgument(name string) *FunctionArgument {
	return &FunctionArgument{
		Name:     name,
		Type:     NewUnknown(),
		Optional: false,
	}
}
