package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// local F = ...
var MetaVariableAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: [
				(vararg_expression) @assignment.right
			] 
		)
	) @meta`,
}

type MetaOrigin struct {
	origin
	definition string
}

func (mo *MetaOrigin) Definition() []string {
	return strings.Split(mo.definition, "\n")
}

func NewMetaOrigin(location nvim.Location, captures nvim.TsQueryMatch, documentation []string) *MetaOrigin {
	metaOrigin := MetaOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: documentation,
		},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "meta":
			metaOrigin.definition = capture.Node.Text
		}
	}

	return &metaOrigin
}
