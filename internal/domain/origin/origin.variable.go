package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
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
	assignmentRightCapture, _ := slicesx.FindFunc(vo.origin.definition.Match, func(capture treesitter.Capture) bool {
		return capture.Id == "assignment.right"
	})

	return assignmentRightCapture.Node.Text
}

func NewVariableOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *VariableOrigin {
	variableOrigin := VariableOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}

	for _, capture := range definition.Match {
		switch capture.Id {
		case "variable":
			variableOrigin.definition = capture.Node.Text
		case "assignment.right":
			variableOrigin.name = capture.Node.Text
		}
	}

	return &variableOrigin
}
