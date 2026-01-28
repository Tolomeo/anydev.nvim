package symbol

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
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
}

func (ao *ModuleOrigin) GetModuleName() string {
	moduleNameCapture, _ := slicesx.FindFunc(ao.definition.Match, func(capture treesitter.Capture) bool {
		return capture.Id == "require.module"
	})

	return moduleNameCapture.Node.Text
}

func (mo *ModuleOrigin) Definition() []string {
	root, _ := mo.definition.Match.Find("module")
	return strings.Split(root.Node.Text, "\n")
}

func NewModuleOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *ModuleOrigin {
	return &ModuleOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}
