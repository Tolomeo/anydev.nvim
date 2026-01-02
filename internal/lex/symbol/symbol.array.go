package symbol

func NewArray(items Symbol) *Array {
	return &Array{
		Kind:  ArrayKindArray,
		Items: items,
	}
}
