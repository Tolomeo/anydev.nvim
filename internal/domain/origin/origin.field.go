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
	) @field`, anyTypeAnnotationQuery),
}

type FieldOrigin struct {
	origin
}

func (fo *FieldOrigin) Definition() []string {
	root, _ := fo.definition.Match.Find("field")
	return strings.Split(root.Node.Text, "\n")
}

func (fo *FieldOrigin) GetName() string {
	fieldName, _ := fo.definition.Match.Find("field.name")
	return fieldName.Node.Text
}

func (fo *FieldOrigin) GetType() string {
	fieldType, _ := fo.definition.Match.Find("field.type")
	return fieldType.Node.Text
}

func NewFieldOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *FieldOrigin {
	return &FieldOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}
