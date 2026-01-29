package type_

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var LiteralTableTypeEmptyQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(table_literal_type . "{" . "}" . ) @table`,
}

var LiteralTableTypeQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(table_literal_type
		"{"
		field: ([
			(
				"["
				.
				(_) @table.index.key
				.
				"]"
				.
				"?"? @table.index.optional
				.
				":"
				.
				(%s) @table.index.value
			) @table.index
			(
				(identifier) @table.field.name
				.
				"?"? @table.field.optional
				.
				":"
				.
				(%s) @table.field.value
			) @table.field
		]
		","?
		)+
		"}"
	) @table`, AnyTypeQuery, AnyTypeQuery),
}

type literalTableTypeField struct {
	key      string
	value    string
	optional bool
}

func (lttf literalTableTypeField) Key() string {
	return lttf.key
}

func (lttf literalTableTypeField) Value() string {
	if lttf.optional {
		return fmt.Sprintf("(%s)?", lttf.value)
	}

	return lttf.value
}

type literalTableTypeIndex struct {
	key      string
	value    string
	optional bool
}

func (ltti literalTableTypeIndex) Key() string {
	return ltti.key
}

func (ltti literalTableTypeIndex) Value() string {
	if ltti.optional {
		return fmt.Sprintf("(%s)?", ltti.value)
	}

	return ltti.value
}

type LiteralTableType struct {
	fields  []literalTableTypeField
	indexes []literalTableTypeIndex
}

func (ltt *LiteralTableType) Fields() []literalTableTypeField {
	return ltt.fields
}

func (ltt *LiteralTableType) Indexes() []literalTableTypeIndex {
	return ltt.indexes
}

func NewLiteralTableType(captures nvim.TsQueryMatch) *LiteralTableType {
	literalTableType := LiteralTableType{
		fields:  []literalTableTypeField{},
		indexes: []literalTableTypeIndex{},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "table.field.name":
			field := literalTableTypeField{key: capture.Node.Text, value: "any"}
			literalTableType.fields = append(literalTableType.fields, field)
		case "table.field.optional":
			literalTableType.fields[len(literalTableType.fields)-1].optional = true
		case "table.field.value":
			literalTableType.fields[len(literalTableType.fields)-1].value = capture.Node.Text

		case "table.index.key":
			index := literalTableTypeIndex{key: capture.Node.Text, value: "any"}
			literalTableType.indexes = append(literalTableType.indexes, index)
		case "table.index.optional":
			literalTableType.indexes[len(literalTableType.indexes)-1].optional = true
		case "table.index.value":
			literalTableType.indexes[len(literalTableType.indexes)-1].value = capture.Node.Text
		}
	}

	return &literalTableType
}
