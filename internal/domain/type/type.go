package type_

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

type Type interface {
	GetKind() string
}

type WithType interface {
	Type() Type
}

type withType struct {
	type_ Type
}

func (w withType) Type() Type {
	return w.type_
}

var _ WithType = (*withType)(nil)

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
