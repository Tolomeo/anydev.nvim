package definition

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// F = 0
var ValueAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			.
			name: (identifier) @value.name
			.
		)
		(expression_list
			.
			value: [
				(number) @value.number
				(string) @value.string
				(true) @value.boolean
				(false) @value.boolean
			] @value.value
			.
		)
	) @value`,
}

// F.X = 0
var ValueFieldDotIndexAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			.
			name: (dot_index_expression
					table: (identifier) @value.parent
					field: (identifier) @value.name
				)
			.
		)
		(expression_list
			.
			value: [
				(number) @value.number
				(string) @value.string
				(true) @value.boolean
				(false) @value.boolean
			] @value.value
			.
		)
	) @value`,
}

/*
	{
		...
	    fieldName = 0,
		...
	}
*/
var ValueFieldAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(field
		name: (identifier) @value.name
		value: [
			(number) @value.number
			(string) @value.string
			(true) @value.boolean
			(false) @value.boolean
		] @value.value
	) @value`,
}

/*
	{
		...
	    ['fieldName'] = 0,
		...
	}
*/
var ValueFieldIndexAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(field
		name: (string
			content: (string_content) @value.name
		)
		value: [
			(number) @value.number
			(string) @value.string
			(true) @value.boolean
			(false) @value.boolean
		] @value.value
	) @value`,
}

type Value struct {
	root  treesitter.TsNode
	value treesitter.TsNode
	type_ string
}

func (v *Value) Root() treesitter.TsNode {
	return v.root
}

func (v *Value) Value() treesitter.TsNode {
	return v.value
}

func (v *Value) Type() string {
	return v.type_
}

func NewValue(match nvim.TsQueryMatch) *Value {
	value := Value{}

	for _, capture := range match {
		switch capture.Id {
		case "value":
			value.root = capture.Node
		case "value.value":
			value.value = capture.Node
		case "value.number":
			value.type_ = "number"
		case "value.string":
			value.type_ = "string"
		case "value.boolean":
			value.type_ = "boolean"
		}
	}

	return &value
}
