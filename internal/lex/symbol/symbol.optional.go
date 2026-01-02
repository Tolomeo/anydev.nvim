package symbol

func NewOptional(typ *Symbol) *Optional {
	return &Optional{
		Kind: OptionalKindOptional,
		Type: typ,
	}
}
