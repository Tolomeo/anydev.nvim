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
}

func (vo *VariableOrigin) Definition() []string {
	root, _ := vo.definition.Match.Find("variable")
	return strings.Split(root.Node.Text, "\n")
}

func (vo *VariableOrigin) GetAssignedName() string {
	assignmentRightCapture, _ := slicesx.FindFunc(vo.definition.Match, func(capture treesitter.Capture) bool {
		return capture.Id == "assignment.right"
	})

	return assignmentRightCapture.Node.Text
}

func NewVariableOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *VariableOrigin {
	return &VariableOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}
