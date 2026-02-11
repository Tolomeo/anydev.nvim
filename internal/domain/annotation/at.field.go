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
		([
			(indexed_field (identifier) @field.index)
			(identifier) @field.name
		])
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
	root      treesitter.TsNode
	name      *treesitter.TsNode
	index     *treesitter.TsNode
	private   *treesitter.TsNode
	protected *treesitter.TsNode
	package_  *treesitter.TsNode
	optional  *treesitter.TsNode
	type_     treesitter.TsNode
}

func (af *AtField) Root() treesitter.TsNode {
	return af.root
}

func (af *AtField) Name() *treesitter.TsNode {
	return af.name
}

func (af *AtField) Index() *treesitter.TsNode {
	return af.index
}

func (af *AtField) Private() *treesitter.TsNode {
	return af.private
}

func (af *AtField) Protected() *treesitter.TsNode {
	return af.protected
}

func (af *AtField) Package() *treesitter.TsNode {
	return af.package_
}

func (af *AtField) Optional() *treesitter.TsNode {
	return af.optional
}

func (af *AtField) Type() treesitter.TsNode {
	return af.type_
}

func NewAtField(match nvim.TsQueryMatch) *AtField {
	atField := AtField{}

	for _, capture := range match {
		switch capture.Id {
		case "field":
			atField.root = capture.Node
		case "field.name":
			atField.name = &capture.Node
		case "field.index":
			atField.index = &capture.Node
		case "field.private":
			atField.private = &capture.Node
		case "field.protected":
			atField.protected = &capture.Node
		case "field.package":
			atField.package_ = &capture.Node
		case "field.optional":
			atField.optional = &capture.Node
		case "field.type":
			atField.type_ = capture.Node
		}
	}

	return &atField
}
