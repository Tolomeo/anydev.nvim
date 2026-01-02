package symbol

func NewReference(value string) *Reference {
	return &Reference{
		Kind:  ReferenceKindReference,
		Value: value,
	}
}
