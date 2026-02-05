package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

var ClassAnnotationQuery = annotation.ClassAnnotationQuery

type ClassOrigin struct {
	origin
	classAnnotation annotation.AtClass
}

func (co *ClassOrigin) Definition() []string {
	return strings.Split(co.classAnnotation.Match(), "\n")
}

func (co *ClassOrigin) Name() string {
	return co.classAnnotation.Name()
}

func NewClassOrigin(location nvim.Location, captures nvim.TsQueryMatch) *ClassOrigin {
	classOrigin := ClassOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: []string{},
		},
		classAnnotation: *annotation.NewClass(captures),
	}

	return &classOrigin
}
