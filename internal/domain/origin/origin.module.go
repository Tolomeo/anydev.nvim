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

type ModuleOrigin struct {
	origin
	definition string
	name       string
}

func (ao *ModuleOrigin) Name() string {
	return ao.name
}

func (mo *ModuleOrigin) Definition() []string {
	return strings.Split(mo.definition, "\n")
}

func NewModuleOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *ModuleOrigin {
	moduleOrigin := ModuleOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}

	for _, capture := range definition.Match {
		switch capture.Id {
		case "module":
			moduleOrigin.definition = capture.Node.Text
		case "require.module":
			moduleOrigin.name = capture.Node.Text
		}
	}

	return &moduleOrigin
}
