package symbol

type Symbol struct {
	withMeta
	Name string
	Type Type
}

func NewSymbol(name string, type_ Type) Symbol {
	return Symbol{
		Name: name,
		Type: type_,
	}
}
