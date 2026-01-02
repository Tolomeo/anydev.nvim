package symbol

func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{
		Kind:  StringLiteralKindStringliteral,
		Value: value,
	}
}
