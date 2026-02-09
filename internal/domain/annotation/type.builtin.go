package annotation

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

const (
	BuiltinVoid          string = "void"
	BuiltinNil           string = "nil"
	BuiltinAny           string = "any"
	BuiltinBoolean       string = "boolean"
	BuiltinString        string = "string"
	BuiltinNumber        string = "number"
	BuiltinInteger       string = "integer"
	BuiltinInt           string = "int"
	BuiltinFunction      string = "function"
	BuiltinTable         string = "table"
	BuiltinThread        string = "thread"
	BuiltinUserdata      string = "userdata"
	BuiltinLightUserdata string = "lightuserdata"
)

var TypeQueries = map[string]string{
	"builtin_type":         "(builtin_type)",
	"identifier":           "(identifier)",
	"array_type":           "(array_type)",
	"table_type":           "(table_type)",
	"table_literal_type":   "(table_literal_type)",
	"union_type":           "(union_type)",
	"parenthesized_type":   "(parenthesized_type)",
	"tuple_type":           "(tuple_type)",
	"function_type":        "(function_type)",
	"member_type":          "(member_type)",
	"optional_type":        "(optional_type)",
	"literal_type":         "(literal_type)",
	"numeric_literal_type": "(numeric_literal_type)",
	"custom_type":          "(custom_type)",
}

var AnyTypeQuery = fmt.Sprintf(`[%s]`, strings.Join(mapx.Values(TypeQueries), " "))
