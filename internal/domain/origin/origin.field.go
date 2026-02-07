package origin

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

var FieldAnnotationQuery = annotation.AtFieldQuery

type FieldAnnotationOrigin struct {
	origin
	fieldAnnotation annotation.AtField
}

func (fo *FieldAnnotationOrigin) Definition() []string {
	return strings.Split(fo.fieldAnnotation.Match(), "\n")
}

func (fo *FieldAnnotationOrigin) Name() *string {
	return fo.fieldAnnotation.Name()
}

func (fo *FieldAnnotationOrigin) Index() *string {
	return fo.fieldAnnotation.Index()
}

func (fo *FieldAnnotationOrigin) Private() bool {
	return fo.fieldAnnotation.Private()
}

func (fo *FieldAnnotationOrigin) Protected() bool {
	return fo.fieldAnnotation.Protected()
}

func (fo *FieldAnnotationOrigin) Package() bool {
	return fo.fieldAnnotation.Package()
}

func (fo *FieldAnnotationOrigin) Type() string {
	if fo.fieldAnnotation.Optional() {
		return fmt.Sprintf("(%s)?", fo.fieldAnnotation.Type())
	}

	return fo.fieldAnnotation.Type()
}

func NewFieldAnnotationOrigin(location nvim.Location, captures nvim.TsQueryMatch) *FieldAnnotationOrigin {
	fieldOrigin := FieldAnnotationOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: []string{},
		},
		fieldAnnotation: *annotation.NewAtField(captures),
	}

	return &fieldOrigin
}
