package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// F.T = X
var VariableAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: (identifier) @assignment.right 
		)
	) @variable`,
}

// F.T = X.V
var VariableDotFieldAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: (dot_index_expression) @assignment.right 
		)
	) @variable`,
}

type VariableOrigin struct {
	origin
	definition string
	name       string
}

func (vo *VariableOrigin) Definition() []string {
	return strings.Split(vo.definition, "\n")
}

func (vo *VariableOrigin) Name() string {
	return vo.name
}

func NewVariableOrigin(location nvim.Location, captures nvim.TsQueryMatch, documentation []string) *VariableOrigin {
	variableOrigin := VariableOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: documentation,
		},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "variable":
			variableOrigin.definition = capture.Node.Text
		case "assignment.right":
			variableOrigin.name = capture.Node.Text
		}
	}

	return &variableOrigin
}
