package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// F = require("T")
var ModuleRequireAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list)
		(expression_list
			value: (function_call
				name: (identifier) @require.call
				arguments: (arguments
					(string
						content: (string_content) @require.module
					)
				)
			)
		) @require
		(#eq? @require.call "require")
	) @module`,
}

type ModuleRequireOrigin struct {
	origin
	definition string
	name       treesitter.TsNode
}

func (mro *ModuleRequireOrigin) Name() string {
	return mro.name.Text
}

func (mro *ModuleRequireOrigin) NameRange() treesitter.Range {
	return mro.name.Range
}

func (mro *ModuleRequireOrigin) Definition() []string {
	return strings.Split(mro.definition, "\n")
}

func NewModuleOrigin(location nvim.Location, captures nvim.TsQueryMatch, documentation []string) *ModuleRequireOrigin {
	moduleOrigin := ModuleRequireOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: documentation,
		},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "module":
			moduleOrigin.definition = capture.Node.Text
		case "require.module":
			moduleOrigin.name = capture.Node
		}
	}

	return &moduleOrigin
}
