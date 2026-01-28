package origin

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var ClassAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(class_annotation
		"@class"
		.
		"(exact)"?
		.
		(identifier) @class.name
		.
		(":" 
			. (%s) @class.parent
			("," (%s) @class.parent)*
		)?
	) @class`, anyTypeQuery, anyTypeQuery),
}

type ClassOrigin struct {
	origin
	definition string
	name       string
	parents    []string
}

func (co *ClassOrigin) Definition() []string {
	return strings.Split(co.definition, "\n")
}

func (co *ClassOrigin) Name() string {
	return co.name
}

func NewClassOrigin(location nvim.Location, captures nvim.TsQueryMatch, documentation []string) *ClassOrigin {
	classOrigin := ClassOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: documentation,
		},
		parents: []string{},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "class":
			classOrigin.definition = capture.Node.Text
		case "class.name":
			classOrigin.name = capture.Node.Text
		case "class.parent":
			classOrigin.parents = append(classOrigin.parents, capture.Node.Text)
		}
	}

	return &classOrigin
}
