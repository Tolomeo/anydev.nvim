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
