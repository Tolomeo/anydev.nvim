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
}

func (co *ClassOrigin) Definition() []string {
	root, _ := co.definition.Match.Find("class")
	return strings.Split(root.Node.Text, "\n")
}

func NewClassOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *ClassOrigin {
	return &ClassOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

