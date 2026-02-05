package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AtFieldQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(field_annotation
		"@field"
		.
		([
			(qualifier "public")
			(qualifier "private") @field.private
			(qualifier "protected") @field.protected
			(qualifier "package") @field.package
		 ])?
		.
		(identifier) @field.name
		.
		"?"? @field.optional
		.
		(%s) @field.type
		.
		(comment)? @field.documentation
		.
	) @field`, AnyTypeQuery),
}

type AtField struct {
	match     string
	name      string
	private   bool
	protected bool
	package_  bool
	optional  bool
	type_     string
}

func (af *AtField) Match() string {
	return af.match
}

func (af *AtField) Name() string {
	return af.name
}

func (af *AtField) Private() bool {
	return af.private
}

func (af *AtField) Protected() bool {
	return af.protected
}

func (af *AtField) Package() bool {
	return af.package_
}

func (af *AtField) Optional() bool {
	return af.optional
}

func (af *AtField) Type() string {
	if af.optional {
		return fmt.Sprintf("(%s)?", af.type_)
	}

	return af.type_
}

func NewAtField(match nvim.TsQueryMatch) *AtField {
	atField := AtField{}

	for _, capture := range match {
		switch capture.Id {
		case "field":
			atField.match = capture.Node.Text
		case "field.name":
			atField.name = capture.Node.Text
		case "field.private":
			atField.private = true
		case "field.protected":
			atField.protected = true
		case "field.package":
			atField.package_ = true
		case "field.optional":
			atField.optional = true
		case "field.type":
			atField.type_ = capture.Node.Text
		}
	}

	return &atField
}
