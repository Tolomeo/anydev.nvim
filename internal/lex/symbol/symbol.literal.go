package symbol

func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{
		Kind:  StringLiteralKindStringliteral,
		Value: value,
	}
}

func NewNumericLiteral(value string) *NumericLiteral {
	return &NumericLiteral{
		Kind:  NumericLiteralKindNumericliteral,
		Value: value,
	}
}

func NewBooleanLiteral(value string) *BooleanLiteral {

	return &BooleanLiteral{
		Kind:  BooleanLiteralKindBooleanliteral,
		Value: BooleanLiteralValue(value),
	}
}
