package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (l *Lexer) lexValue() (symbol.Symbol, error) {
	target := l.context.Target()
	buffer, err := l.context.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	sourceDefinition := target.Origin().DefinitionLines()
	err = buffer.SetLines(sourceDefinition)

	if err != nil {
		return nil, err
	}

	table, err := l.lexTableValue()

	switch {
	case err != nil:
		return nil, err
	case table != nil:
		return table, nil
	}

	function, err := l.lexFunctionValue()

	switch {
	case err != nil:
		return nil, err
	case function != nil:
		return function, nil
	}

	meta, err := l.lexMetaValue()

	switch {
	case err != nil:
		return nil, err
	case meta != nil:
		return meta, nil
	}

	return nil, fmt.Errorf("Error lexing %s: unknown origin <%+v>", target.Identifier(), target.Origin())
}

var functionQueries = map[string]string{
	/* function fn() end
	function fn(arg1) end
	function fn(arg1, arg2) end
	function fn(arg1, arg2, ...) end
		function fn(...) end */
	"declaration": `
		(function_declaration
			name: (identifier) @name
			parameters: (parameters
				(identifier)? @arg
				("," (identifier) @arg)*
				("," (vararg_expression) @vararg)?
				(vararg_expression)? @vararg
			)
		)
	`,
	/* local T = function() end
	local M = function(arg) end
	local D = function(arg, ...) end
	local E = function(...) end */
	"assignment": `
		(variable_declaration
			(assignment_statement
				(variable_list
					name: (identifier) @name
				)
				(expression_list
					value: (function_definition
						parameters: (parameters
							(identifier)? @arg
							("," (identifier) @arg)*
							("," (vararg_expression) @vararg)?
							(vararg_expression)? @vararg
						)
					)
				)
			)
		)
	`,
	/* function api:fn() end
	function api:fn(name) end
	function api:fn(name, value) end
	function api:fn(name, value, ...) end
	function api:fn(...) end */
	"methodDeclaration": `
		(function_declaration
			name: (method_index_expression
				method: (identifier) @name
			) @access.instance
			parameters: (parameters
				(identifier)? @arg
				("," (identifier) @arg)*
				("," (vararg_expression) @vararg)?
				(vararg_expression)? @vararg
			)
		)
	`,
	/* function api.fn() end
	function api.fn(name) end
	function api.fn(name, value) end
	function api.fn(name, value, ...) end
	function api.fn(...) end */
	"methodIndexDeclaration ": `
		(function_declaration
			name: (dot_index_expression
				field: (identifier) @name
			) @access.class
			parameters: (parameters
				(identifier)? @arg
				("," (identifier) @arg)*
				("," (vararg_expression) @vararg)?
				(vararg_expression)? @vararg
			)
		)
	`,
	/* api.fn = function() end
	api.fn = function(name) end
	api.fn = function(name, value) end
	api.fn = function(name, value, ...) end
	api.fn = function(...) end */
	"methodAssignment": `
		(assignment_statement
			(variable_list
				name: (dot_index_expression
					field: (identifier) @name
				) @access.class
			)
			(expression_list
				value: (function_definition
					parameters: (parameters
						(identifier)? @arg
						("," (identifier) @arg)*
						("," (vararg_expression) @vararg)?
						(vararg_expression)? @vararg
					)
				)
			)
		)
	`,
	/* api['fn'] = function() end
	api['fn'] = function(name) end
	api['fn'] = function(name, value) end
	api['fn'] = function(name, value, ...) end
	api['fn'] = function(...) end */
	"methodIndexAssignment": `
		(assignment_statement
			(variable_list
				name: (bracket_index_expression
					table: (_)
					field: (string
						content: (string_content) @name
					)
				) @access.class
			)
			(expression_list
				value: (function_definition
					parameters: (parameters
						(identifier)? @arg
						("," (identifier) @arg)*
						("," (vararg_expression) @vararg)?
						(vararg_expression)? @vararg
					)
				)
			)
		)
	`,
}

func (l *Lexer) lexFunctionValue() (*symbol.Function, error) {
	origin := l.context.Target().Origin()
	buffer, err := l.context.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.DefinitionLines())

	if err != nil {
		return nil, err
	}

	var match *nvim.TsQueryMatch

	for _, query := range functionQueries {
		captures, err := buffer.TsQueryOne(treesitter.Query{Language: "lua", Query: query})

		if err != nil {
			return nil, err
		}

		if captures != nil {
			match = captures
			break
		}
	}

	if match == nil {
		return nil, nil
	}

	function := symbol.NewFunction()
	function.Documentation = origin.DocumentationLines()

	for _, capture := range *match {
		switch capture.Id {
		case "name":
			function.Name = &capture.Node.Text
		case "access.class":
			function.Access = &symbol.FunctionClassAccess
		case "access.instance":
			function.Access = &symbol.FunctionIstanceAccess
		case "arg":
			function.Arguments = append(function.Arguments, *symbol.NewFunctionArgument(capture.Node.Text))
		case "vararg":
			function.Arguments = append(function.Arguments, *symbol.NewFunctionArgument(capture.Node.Text))
		}
	}

	annotations, err := l.lexAtAnnotations(function.Documentation)

	if err != nil {
		return nil, fmt.Errorf("Error lexing function %s: %w", *function.Name, err)
	}

	for _, genericAnnotation := range annotations.AtGenerics {
		genericName := genericAnnotation.Name
		genericTypes, err := slicesx.MapFunc(genericAnnotation.Types, func(genericType TypeAnnotation) (symbol.Symbol, error) {
			return l.lexTypeAnnotation(genericType)
		})

		if err != nil {
			return nil, err
		}

		function.Generics = append(function.Generics, *symbol.NewFunctionGeneric(genericName, genericTypes...))
	}

	for argIndex := range function.Arguments {
		name := function.Arguments[argIndex].Name

		paramAnnotation, hasParamAnnotation := annotations.AtParams[name]

		if !hasParamAnnotation {
			l.context.Logger().Info(fmt.Sprintf("Using '%s' for undocumented argument type '%s'", function.Arguments[argIndex].Type, function.Arguments[argIndex].Name))
			continue
		}

		if functionGeneric, isGeneric := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
			return generic.Name == paramAnnotation.Type.Name
		}); isGeneric {
			function.Arguments[argIndex].Type = symbol.NewReference(functionGeneric.Name)
		} else {
			argumentType, err := l.lexTypeAnnotation(paramAnnotation.Type)

			if err != nil {
				return nil, err
			}

			function.Arguments[argIndex].Type = argumentType
		}

		function.Arguments[argIndex].Optional = paramAnnotation.Optional
		function.Arguments[argIndex].Documentation = paramAnnotation.Documentation
	}

	for _, returnAnnotation := range annotations.AtReturns {
		functionReturn := symbol.NewFunctionReturn()
		functionReturn.Name = returnAnnotation.Name
		functionReturn.Documentation = returnAnnotation.Documentation

		if functionGeneric, isGeneric := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
			return generic.Name == returnAnnotation.Type.Name
		}); isGeneric {
			functionReturn.Type = symbol.NewReference(functionGeneric.Name)
		} else {
			typ, err := l.lexTypeAnnotation(returnAnnotation.Type)

			if err != nil {
				return nil, err
			}

			functionReturn.Type = typ
		}

		function.Returns = append(function.Returns, *functionReturn)
	}

	for _, overloadAnnotation := range annotations.AtOverloads {
		overloadType, err := l.lexTypeAnnotation(overloadAnnotation.Type)

		if err != nil {
			return nil, fmt.Errorf("Error lexing function overload annotation type: %w", err)
		}

		overloadFunctionType, isFunctionType := overloadType.(*symbol.Function)

		if !isFunctionType {
			return nil, fmt.Errorf("Error lexing function overload annotation type: lexed type '%+v' is not a function", overloadFunctionType)
		}

		functionOverload := symbol.NewFunctionOverload()
		functionOverload.Generics = overloadFunctionType.Generics
		functionOverload.Arguments = overloadFunctionType.Arguments
		functionOverload.Documentation = overloadFunctionType.Documentation
		functionOverload.Returns = overloadFunctionType.Returns

		function.Overloads = append(function.Overloads, *functionOverload)
	}

	return function, nil
}

var tableQueries = map[string]string{
	"tableDeclaration": `
	(variable_declaration
		(assignment_statement
			(variable_list
				name: (identifier)
			) @table.name
			(expression_list
				value: (table_constructor)
			) @table.value
		)
	)
`,
	"tableFieldDeclaration": `
	(assignment_statement
		(variable_list
			name: (dot_index_expression
				table: (_)
				field: (identifier) @table.name
			)
		)
		(expression_list
			value: (table_constructor) @table.value
		)
	)
`,
	"tableIndexFieldDeclaration": `
		(assignment_statement
			(variable_list
				name: (bracket_index_expression
					table: (_)
					field: (string
						content: (string_content) @table.name
					)
				)
			)
			(expression_list
				value: (table_constructor) @table.value
			)
		)
`}

func (l *Lexer) lexTableValue() (*symbol.Table, error) {
	target := l.context.Target()
	origin := target.Origin()
	buffer, err := l.context.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.DefinitionLines())

	if err != nil {
		return nil, err
	}

	var match *nvim.TsQueryMatch

	for _, query := range tableQueries {
		captures, err := buffer.TsQueryOne(treesitter.Query{Language: "lua", Query: query})

		if err != nil {
			return nil, err
		}

		if captures != nil {
			match = captures
			break
		}
	}

	if match == nil {
		return nil, nil
	}

	table := symbol.NewTable()

	for _, capture := range *match {
		switch capture.Id {
		case "table.name":
			table.Name = capture.Node.Text
		}
	}

	tableFields, err := l.context.Nvim().GetValueCompletion(target.Identifier())

	if err != nil {
		return nil, err
	}

	for _, fieldName := range tableFields {
		fieldValue, err := l.context.ExtractChild(fieldName)

		if err != nil {
			return nil, err
		}

		/* annotations, err := l.lexAtAnnotations(tableFieldSource.GetOrigin().DefinitionLines())
		tableField.Private = annotations.AtPrivate
		tableField.Protected = annotations.AtProtected
		tableField.Package = annotations.AtPackage
		tableField.Deprecated = annotations.AtDeprecated
		tableField.Protected = annotations.AtProtected */
		table.Fields = append(table.Fields, symbol.TableField{
			Name:  fieldName,
			Value: fieldValue,
		})
	}

	return table, nil
}

var metaQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: [
				(vararg_expression) @assignment.right
			] 
		)
	) @assignment`,
}

func (l *Lexer) lexMetaValue() (symbol.Symbol, error) {
	target := l.context.Target()
	origin := target.Origin()
	buffer, err := l.context.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.DefinitionLines())

	if err != nil {
		return nil, err
	}

	captures, err := buffer.TsQueryOne(metaQuery)

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	unknown := symbol.NewUnknown()
	unknown.Documentation = origin.DocumentationLines()

	annotations, err := l.lexAtAnnotations(origin.DocumentationLines())

	switch {
	case err != nil:
		return nil, err
	case annotations.AtType == nil:
		l.context.Logger().Warn(fmt.Sprintf("Unknown meta type '%s' received", l.context.Target().Name()))
		return unknown, nil
	case len(annotations.AtType.Types) < 1:
		return nil, fmt.Errorf("Error lexing @type annotations for meta type '%s': no type annotations found", l.context.Target().Name())
	}

	lexedType, err := l.lexTypeAnnotation(annotations.AtType.Types[0])

	if err != nil {
		return nil, err
	}

	return lexedType, nil
}
