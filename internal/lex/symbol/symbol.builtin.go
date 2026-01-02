package symbol

func newBuiltinType(value BuiltinValue) *Builtin {
	return &Builtin{
		Kind:  BuiltinKindBuiltin,
		Value: value,
	}
}

func NewVoid() *Builtin {
	return newBuiltinType(BuiltinValueVoid)
}

func NewNil() *Builtin {
	return newBuiltinType(BuiltinValueNil)
}

func NewAny() *Builtin {
	return newBuiltinType(BuiltinValueAny)
}

func NewBoolean() *Builtin {
	return newBuiltinType(BuiltinValueBoolean)
}

func NewString() *Builtin {
	return newBuiltinType(BuiltinValueString)
}

func NewNumber() *Builtin {
	return newBuiltinType(BuiltinValueNumber)
}

func NewInteger() *Builtin {
	return newBuiltinType(BuiltinValueInteger)
}

func NewBuiltinFunction() *Builtin {
	return newBuiltinType(BuiltinValueFunction)
}

func NewBuiltinTable() *Builtin {
	return newBuiltinType(BuiltinValueTable)
}

func NewThread() *Builtin {
	return newBuiltinType(BuiltinValueTable)
}

func NewUserdata() *Builtin {
	return newBuiltinType(BuiltinValueUserdata)
}

func NewLightUserdata() *Builtin {
	return newBuiltinType(BuiltinValueLightuserdata)
}
