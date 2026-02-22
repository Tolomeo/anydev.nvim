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

// NOTE: this does not support function calls without parens
var FunctionCallAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			.
			name: [
				(identifier) @name
				(dot_index_expression
					table: (identifier) @parent
					field: (identifier) @name
				)
			]
			.
		)
		(expression_list
			value: (function_call
				name: (_) @function_call.name
				arguments: (
					(arguments
						.
						"("
						.
						((_) * @function_call.argument . ","? )*
						.
						")"
						.
					)
				)
			)
		)
		(#not-any-of? @function_call.name "require" "memoize")
	) @function_call
	(#not-has-ancestor? @function_call "function_call") ; Avoiding nested matches
	`,
}

type FunctionCall struct {
	root              treesitter.TsNode
	parent            *treesitter.TsNode
	name              treesitter.TsNode
	functionName      treesitter.TsNode
	functionArguments []treesitter.TsNode
}

func (fc *FunctionCall) Root() treesitter.TsNode {
	return fc.root
}

func (fc *FunctionCall) Parent() *treesitter.TsNode {
	return fc.parent
}

func (fc *FunctionCall) Name() treesitter.TsNode {
	return fc.name
}

func (fc *FunctionCall) FunctionName() treesitter.TsNode {
	return fc.functionName
}

func (fc *FunctionCall) FunctionArguments() []treesitter.TsNode {
	return fc.functionArguments
}

func NewFunctionCall(match nvim.TsQueryMatch) *FunctionCall {
	functionCallOrigin := FunctionCall{
		functionArguments: []treesitter.TsNode{},
	}

	for _, capture := range match {
		switch capture.Id {
		case "function_call":
			functionCallOrigin.root = capture.Node
		case "name":
			functionCallOrigin.name = capture.Node
		case "parent":
			functionCallOrigin.parent = &capture.Node
		case "function_call.name":
			functionCallOrigin.functionName = capture.Node
		case "function_call.argument":
			functionCallOrigin.functionArguments = append(functionCallOrigin.functionArguments, capture.Node)
		}
	}

	return &functionCallOrigin
}
