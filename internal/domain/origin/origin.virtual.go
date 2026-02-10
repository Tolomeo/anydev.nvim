package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// local F = ...
var VirtualVariableAssignmentQuery = treesitter.Query{
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
	) @virtual`,
}

var VirtualFieldAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(field
		name: (identifier) @virtual.name
		value: (vararg_expression)
	) @virtual
	(#not-has-ancestor? @virtual "field") ; Avoiding nested matches
	`,
}

type VirtualOrigin struct {
	origin
	definition string
}

func (mo *VirtualOrigin) Definition() []string {
	return strings.Split(mo.definition, "\n")
}

func NewVirtualOrigin(location nvim.Location, captures nvim.TsQueryMatch, annotations []string) *VirtualOrigin {
	virtualOrigin := VirtualOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: annotations,
		},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "virtual":
			virtualOrigin.definition = capture.Node.Text
		}
	}

	return &virtualOrigin
}
