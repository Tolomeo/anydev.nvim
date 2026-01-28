package origin

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var FieldAnnotationQuery = treesitter.Query{
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
	) @field`, anyTypeQuery),
}

type FieldOrigin struct {
	origin
	name      string
	private   bool
	protected bool
	package_  bool
	optional  bool
	type_     string
}

func (fo *FieldOrigin) Definition() []string {
	root, _ := fo.definition.Match.Find("field")
	return strings.Split(root.Node.Text, "\n")
}

func (fo *FieldOrigin) Name() string {
	return fo.name
}

func (fo *FieldOrigin) Private() bool {
	return fo.private
}

func (fo *FieldOrigin) Protected() bool {
	return fo.protected
}

func (fo *FieldOrigin) Package() bool {
	return fo.package_
}

func (fo *FieldOrigin) Type() string {
	if fo.optional {
		return fmt.Sprintf("(%s)?", fo.type_)
	}

	return fo.type_
}

func NewFieldOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *FieldOrigin {
	fieldOrigin := FieldOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}

	for _, capture := range definition.Match {
		switch capture.Id {
		case "field.name":
			fieldOrigin.name = capture.Node.Text
		case "field.private":
			fieldOrigin.private = true
		case "field.protected":
			fieldOrigin.protected = true
		case "field.package":
			fieldOrigin.package_ = true
		case "field.optional":
			fieldOrigin.optional = true
		case "field.type":
			fieldOrigin.type_ = capture.Node.Text
		}
	}

	return &fieldOrigin
}
