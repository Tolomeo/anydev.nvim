package origin

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

var typeAnnotationQueries = map[string]string{
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

var anyTypeAnnotationQuery = fmt.Sprintf(`[%s]`, strings.Join(mapx.Values(typeAnnotationQueries), " "))

type Origin interface {
	Url() string
	Line() uint
	Character() uint
	Definition() []string
	Documentation() []string
	Captures() nvim.TsQueryMatch
}

type origin struct {
	location   nvim.Location
	definition nvim.TsNodeQueryMatch
	docBlock   []string
}

func (l *origin) Url() string {
	return l.location.Url
}

func (l *origin) Line() uint {
	return l.location.StartLine()
}

func (l *origin) Character() uint {
	return l.location.StartCharacter()
}

func (l *origin) Type() string {
	return l.definition.Node.Type
}

func (l *origin) Documentation() []string {
	return l.docBlock
}

func (l *origin) Captures() nvim.TsQueryMatch {
	return l.definition.Match
}
