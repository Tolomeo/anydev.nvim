package symbol

func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{
		Kind:  StringLiteralKindStringliteral,
		Value: value,
	}
}

func NewNumericLiteral(value string) *NumericLiteral {
	return &NumericLiteral{
		Kind: NumericLiteralKindNumericliteral,
		Value: value,
	}
}
