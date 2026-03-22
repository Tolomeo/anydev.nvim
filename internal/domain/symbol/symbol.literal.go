package symbol

const StringLiteralKind string = "stringliteral"

type StringLiteral struct {
	Kind  string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

func (s *StringLiteral) Canonical() Type {
	return s
}

// GetKind implements Type.
func (s *StringLiteral) GetKind() string {
	return s.Kind
}

var _ Type = (*StringLiteral)(nil)

func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{
		Kind:  StringLiteralKind,
		Value: value,
	}
}

const NumericLiteralKind string = "numericliteral"

type NumericLiteral struct {
	Kind  string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

func (n *NumericLiteral) Canonical() Type {
	return n
}

// GetKind implements Type.
func (n *NumericLiteral) GetKind() string {
	return n.Kind
}

var _ Type = (*NumericLiteral)(nil)

func NewNumericLiteral(value string) *NumericLiteral {
	return &NumericLiteral{
		Kind:  NumericLiteralKind,
		Value: value,
	}
}

const BooleanLiteralKind string = "booleanliteral"

type BooleanLiteral struct {
	Kind  string              `json:"kind" yaml:"kind" mapstructure:"kind"`
	Value BooleanLiteralValue `json:"value" yaml:"value" mapstructure:"value"`
}

func (b *BooleanLiteral) Canonical() Type {
	return b
}

// GetKind implements Type.
func (b *BooleanLiteral) GetKind() string {
	return b.Kind
}

var _ Type = (*BooleanLiteral)(nil)

type BooleanLiteralValue string

const BooleanLiteralValueFalse BooleanLiteralValue = "false"
const BooleanLiteralValueTrue BooleanLiteralValue = "true"

func NewBooleanLiteral(value string) *BooleanLiteral {
	return &BooleanLiteral{
		Kind:  BooleanLiteralKind,
		Value: BooleanLiteralValue(value),
	}
}
