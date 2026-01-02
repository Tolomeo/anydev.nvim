package symbol

func NewUnion(types []Symbol) *Union {
	unionTypes := []UnionTypesElem{}

	for _, typ := range types {
		unionTypes = append(unionTypes, typ)
	}

	return &Union{
		Kind:  UnionKindUnion,
		Types: unionTypes,
	}
}
