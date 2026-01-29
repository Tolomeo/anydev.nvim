package symbol

import "github.com/Tolomeo/anydev.nvim/internal/domain/annotation"

const BuiltinValueAny string = "any"
const BuiltinValueBoolean string = "boolean"
const BuiltinValueFunction string = "function"
const BuiltinValueInteger string = "integer"
const BuiltinValueLightuserdata string = "lightuserdata"
const BuiltinValueNil string = "nil"
const BuiltinValueNumber string = "number"
const BuiltinValueString string = "string"
const BuiltinValueTable string = "table"
const BuiltinValueThread string = "thread"
const BuiltinValueUserdata string = "userdata"
const BuiltinValueVoid string = "void"

// https://luals.github.io/wiki/annotations/#documenting-types
type Builtin struct {
	Kind  string `json:"kind" yaml:"kind" mapstructure:"kind"`
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

func (b *Builtin) GetKind() string {
	return b.Kind
}

var _ annotation.Type = (*Builtin)(nil)

type BuiltinKind string

type BuiltinValue string

func newBuiltinType(value string) *Builtin {
	return &Builtin{
		Kind:  "builtin",
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
