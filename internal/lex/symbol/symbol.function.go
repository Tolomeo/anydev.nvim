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

func NewFunctionArg(name string) *FunctionArg {
	return &FunctionArg{
		Name:     name,
		Type:     NewUnknown(),
		Optional: false,
	}
}
