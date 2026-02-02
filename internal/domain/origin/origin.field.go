package origin

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
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
	) @field`, annotation.AnyTypeQuery),
}

type FieldAnnotationOrigin struct {
	origin
	definition string
	name       string
	private    bool
	protected  bool
	package_   bool
	optional   bool
	type_      string
}

func (fo *FieldAnnotationOrigin) Definition() []string {
	return strings.Split(fo.definition, "\n")
}

func (fo *FieldAnnotationOrigin) Name() string {
	return fo.name
}

func (fo *FieldAnnotationOrigin) Private() bool {
	return fo.private
}

func (fo *FieldAnnotationOrigin) Protected() bool {
	return fo.protected
}

func (fo *FieldAnnotationOrigin) Package() bool {
	return fo.package_
}

func (fo *FieldAnnotationOrigin) Type() string {
	if fo.optional {
		return fmt.Sprintf("(%s)?", fo.type_)
	}

	return fo.type_
}

func NewFieldAnnotationOrigin(location nvim.Location, captures nvim.TsQueryMatch) *FieldAnnotationOrigin {
	fieldOrigin := FieldAnnotationOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: []string{},
		},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "field":
			fieldOrigin.definition = capture.Node.Text
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
