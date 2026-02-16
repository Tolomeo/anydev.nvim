package definition

import (
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

type RequireFunctionCall struct {
	root treesitter.TsNode
	name treesitter.TsNode
}

func (rfc *RequireFunctionCall) Root() treesitter.TsNode {
	return rfc.root
}

func (rfc *RequireFunctionCall) RequiredModuleName() treesitter.TsNode {
	return rfc.name
}

func NewRequireFunctionCall(match nvim.TsQueryMatch) *RequireFunctionCall {
	requireFnCall := RequireFunctionCall{}

	for _, capture := range match {
		switch capture.Id {
		case "module":
			requireFnCall.root = capture.Node
		case "require.module":
			requireFnCall.name = capture.Node
		}
	}

	return &requireFnCall
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
		(#not-any-of? @function_call.name "require" "vim._defer_require" "memoize" "create_option_accessor")
	) @function_call
	(#not-has-ancestor? @function_call "function_call") ; Avoiding nested matches
	`,
}

type FunctionCall struct {
	root         treesitter.TsNode
	functionName treesitter.TsNode
}

func (fc *FunctionCall) Root() treesitter.TsNode {
	return fc.root
}

func (fc *FunctionCall) FunctionName() treesitter.TsNode {
	return fc.functionName
}

func NewFunctionCall(match nvim.TsQueryMatch) *FunctionCall {
	functionCallOrigin := FunctionCall{}

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
