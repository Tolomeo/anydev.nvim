package definition

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// F = ...
var VirtualVariableAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (_) @virtual.name
		)
		(expression_list
			value: (vararg_expression)
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

type Virtual struct {
	root treesitter.TsNode
	name treesitter.TsNode
}

func (v *Virtual) Root() treesitter.TsNode {
	return v.root
}

func (v *Virtual) Name() treesitter.TsNode {
	return v.name
}

func NewVirtual(match nvim.TsQueryMatch) *Virtual {
	virtual := Virtual{}

	for _, capture := range match {
		switch capture.Id {
		case "virtual":
			virtual.root = capture.Node
		case "virtual.name":
			virtual.name = capture.Node
		}
	}

	return &virtual
}
