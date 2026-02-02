package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

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
		name: (identifier) @field.name
		value: [
			(number) @value.number
			(string) @value.string
			(true) @value.boolean
			(false) @value.boolean
		] @field.value
	) @value`,
}

var ValueFieldIndexAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(field
		name: (string
			content: (string_content) @field.name
		)
		value: [
			(number) @value.number
			(string) @value.string
			(true) @value.boolean
			(false) @value.boolean
		] @value.value
	) @value`,
}

type ValueOrigin struct {
	origin
	definition string
	value      string
	type_      string
}

func (vo *ValueOrigin) Definition() []string {
	return strings.Split(vo.definition, "\n")
}

func (vo *ValueOrigin) Value() string {
	return vo.value
}

func (vo *ValueOrigin) Type() string {
	return vo.type_
}

func NewValueOrigin(location nvim.Location, captures nvim.TsQueryMatch, annotations []string) *ValueOrigin {
	valueOrigin := ValueOrigin{}

	for _, capture := range captures {
		switch capture.Id {
		case "value":
			valueOrigin.definition = capture.Node.Text
		case "value.value":
			valueOrigin.value = capture.Node.Text
		case "value.number":
			valueOrigin.type_ = "number"
		case "value.string":
			valueOrigin.type_ = "string"
		case "value.boolean":
			valueOrigin.type_ = "boolean"
		}
	}

	return &valueOrigin

}
