package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// F = require("T")
var RequireFunctionCallAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list)
		(expression_list
			value: (function_call
				name: (_) @require.call
				arguments: (arguments
					(string
						content: (string_content) @require.module
					)
				)
			)
		) @require
		(#any-of? @require.call "require")
	) @module`,
}

type RequireFunctionCallOrigin struct {
	origin
	definition string
	name       treesitter.TsNode
}

func (mro *RequireFunctionCallOrigin) Definition() []string {
	return strings.Split(mro.definition, "\n")
}

func (mro *RequireFunctionCallOrigin) RequiredModuleName() string {
	return mro.name.Text
}

func (mro *RequireFunctionCallOrigin) RequiredModuleNameRange() treesitter.Range {
	return mro.name.Range
}

func NewRequireFunctionCallOrigin(location nvim.Location, captures nvim.TsQueryMatch, documentation []string) *RequireFunctionCallOrigin {
	moduleOrigin := RequireFunctionCallOrigin{
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

var FunctionCallAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list)
		(expression_list
			value: (function_call
				name: (_) @function_call.name
			)
		)
		(#not-any-of? @function_call.name "require" "vim._defer_require")
	) @function_call`,
}

type FunctionCallOrigin struct {
	origin
	root         treesitter.TsNode
	functionName treesitter.TsNode
}

func (fco *FunctionCallOrigin) Definition() []string {
	return strings.Split(fco.root.Text, "\n")
}

func (fco *FunctionCallOrigin) FunctionName() []string {
	return strings.Split(fco.functionName.Text, "\n")
}

func NewFunctionCallOrigin(location nvim.Location, match nvim.TsQueryMatch, annotations []string) *FunctionCallOrigin {
	functionCallOrigin := FunctionCallOrigin{
		origin: origin{
			location:    location,
			captures:    match,
			annotations: annotations,
		},
	}

	for _, capture := range match {
		switch capture.Id {
		case "function_call":
			functionCallOrigin.root = capture.Node
		case "function_call.name":
			functionCallOrigin.functionName = capture.Node
		}
	}

	return &functionCallOrigin
}
