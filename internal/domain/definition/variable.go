package definition

import (
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
			.
			value: (identifier) @name
			.
		) @assignment.right
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
			.
			value: (dot_index_expression) @name
			.
		) @assignment.right
	) @variable`,
}

/*
	F = {
		...
	    T = X,
		...
	}
*/
var VariableFieldAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(field
		name: (identifier) @assignment.left
		value: [
			(identifier) @name
			(dot_index_expression
				table: (_)
				field: (identifier) @name
			)
		] @assignment.right
	) @variable`,
}

type Variable struct {
	root treesitter.TsNode
	name treesitter.TsNode
}

func (v *Variable) Root() treesitter.TsNode {
	return v.root
}

func (v *Variable) Name() treesitter.TsNode {
	return v.name
}

func NewVariable(captures nvim.TsQueryMatch) *Variable {
	variable := Variable{}

	for _, capture := range captures {
		switch capture.Id {
		case "variable":
			variable.root = capture.Node
		case "name":
			variable.name = capture.Node
		/* case "assignment.right":
			variable.name = capture.Node */
		}
	}

	return &variable
}
