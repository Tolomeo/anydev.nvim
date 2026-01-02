package symbol

func NewUnknown() *Unknown {
	return &Unknown{
		Kind: UnknownKindUnknown,
	}
}
