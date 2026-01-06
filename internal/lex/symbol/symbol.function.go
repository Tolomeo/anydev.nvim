package symbol

var (
	FunctionClassAccess   string = "class"
	FunctionIstanceAccess string = "instance"
)

func NewFunction() *Function {
	return &Function{
		Kind: FunctionKindFunction,
		Generics: []FunctionGeneric{},
		Returns: []FunctionReturn{},
	}
}

func NewFunctionArgument(name string) *FunctionArgument {
	return &FunctionArgument{
		Name:     name,
		Type:     NewUnknown(),
		Optional: false,
	}
}

func NewFunctionGeneric(name string, types ...Symbol) *FunctionGeneric {
	genericTypes := []FunctionGenericTypesElem{}

	for _, typ := range types {
		genericTypes = append(genericTypes, typ)
	}

	return &FunctionGeneric{
		Name:  name,
		Types: genericTypes,
	}
}

