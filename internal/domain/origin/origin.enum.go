package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

var EnumAnnotationQuery = annotation.AtEnumQuery

type EnumOrigin struct {
	origin
	definition string
	name       string
}

func (eo *EnumOrigin) Definition() []string {
	return strings.Split(eo.definition, "\n")
}

func (eo *EnumOrigin) Name() string {
	return eo.name
}

func NewEnumAnnotationOrigin(location nvim.Location, captures nvim.TsQueryMatch) *EnumOrigin {
	enumAnnotation := annotation.NewAtEnum(captures)

	return &EnumOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: []string{},
		},
		definition: enumAnnotation.Text(),
		name:       enumAnnotation.Name(),
	}
}
