package lex

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

var ErrNoMatch = errors.New("No match")

func newReferenceType(value string) symbol.Reference {
	return symbol.Reference{
		Kind:  symbol.ReferenceKindReference,
		Value: value,
	}
}

func newFunctionType() *symbol.Function {
	return &symbol.Function{
		Kind: symbol.FunctionKindFunction,
	}
}

func newFunctionTypeArg(name string) symbol.FunctionArg {
	return symbol.FunctionArg{
		Name:     name,
		Type:     newUnknownType(),
		Optional: false,
	}
}

func newBuiltinType(value symbol.BuiltinValue) symbol.Builtin {
	return symbol.Builtin{
		Kind:  symbol.BuiltinKindBuiltin,
		Value: value,
	}
}

func newUnknownType() *symbol.Unknown {
	return &symbol.Unknown{
		Kind: symbol.UnknownKindUnknown,
	}
}

func newTableType() *symbol.Table {
	return &symbol.Table{
		Kind: symbol.TableKindTable,
	}
}

var typeQueries = map[string]string{
	"builtin_type":         "(builtin_type)",
	"identifier":           "(identifier)",
	"array_type":           "(array_type)",
	"table_type":           "(table_type)",
	"table_literal_type":   "(table_literal_type)",
	"union_type":           "(union_type)",
	"parenthesized_type":   "(parenthesized_type)",
	"tuple_type":           "(tuple_type)",
	"function_type":        "(function_type)",
	"member_type":          "(member_type)",
	"optional_type":        "(optional_type)",
	"literal_type":         "(literal_type)",
	"numeric_literal_type": "(numeric_literal_type)",
	"custom_type":          "(custom_type)",
}

var anyTypeQuery = fmt.Sprintf(`[%s]`, strings.Join(mapx.Values(typeQueries), " "))

var typeFunctionQuery string = fmt.Sprintf(`
	(documentation
		(type_annotation
			(function_type
				(parameter
					(identifier) @parameter.name
					":"
					%s @parameter.type
				) @parameter
				("," (parameter
					(identifier) @parameter.name
					":"
					%s @parameter.type
				) @parameter)*
				("," (parameter
					"..." @parameter.name
					":"
					%s @parameter.type
				) @parameter)?
				(parameter
					"..." @parameter.name
					":"
					%s @parameter.type
				)? @parameter
				":"?
				%s? @return.type
			)
		)
	)
`, anyTypeQuery, anyTypeQuery, anyTypeQuery, anyTypeQuery, anyTypeQuery)

func (l *Lexer) lexFunctionType(function *symbol.Function) error {
	captures, err := l.context.Nvim.TsQueryOne(treesitter.Query{Language: "luadoc", Query: typeFunctionQuery})

	switch {
	case err != nil:
		return err
	case captures == nil:
		return ErrNoMatch
	}

	args := []symbol.FunctionArg{}
	returns := []symbol.FunctionReturn{}

	for _, capture := range *captures {
		switch capture.Id {
		case "parameter":
			args = append(args, symbol.FunctionArg{})
		case "parameter.name":
			args[len(args)-1].Name = capture.Node.Text
		case "parameter.type":
			parameterType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return err
			}

			args[len(args)-1].Type = parameterType
		case "return.type":
			returnType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return err
			}

			returns = append(returns, symbol.FunctionReturn{
				Type: returnType,
			})
		}
	}

	function.Args = append(function.Args, args...)
	function.Return = append(function.Return, returns...)

	return nil
}

var typeTableQuery string = `
	(documentation
		(type_annotation
			(table_type
				key: (builtin_type) @key
				value: (builtin_type) @value
			) @table
			(comment)? @documentation
		)
	)
`

func (l *Lexer) lexTableType(table *symbol.Table) error {
	matches, err := l.context.Nvim.TsQueryAll(treesitter.Query{Language: "luadoc", Query: typeTableQuery})

	switch {
	case err != nil:
		return err
	case matches == nil:
		return ErrNoMatch
	}

	// TODO: here match one
	for _, matchCaptures := range *matches {
		for _, capture := range matchCaptures {
			switch capture.Id {
			case "table":
				table.Fields = append(table.Fields, symbol.TableField{})
			case "key":
				table.Fields[len(table.Fields)-1].Name = capture.Node.Text
			case "value":
				valueType, err := l.lexType(capture.Node.Text)

				if err != nil {
					return err
				}

				table.Fields[len(table.Fields)-1].Value = valueType
			}
		}
	}

	return nil
}

func (l *Lexer) lexBuiltinType(source string) symbol.Symbol {
	switch source {
	case "void":
		return newBuiltinType(symbol.BuiltinValueVoid)
	case "nil":
		return newBuiltinType(symbol.BuiltinValueNil)
	case "any":
		return newBuiltinType(symbol.BuiltinValueAny)
	case "boolean":
		return newBuiltinType(symbol.BuiltinValueBoolean)
	case "string":
		return newBuiltinType(symbol.BuiltinValueString)
	case "number":
		return newBuiltinType(symbol.BuiltinValueNumber)
	case "integer", "int":
		return newBuiltinType(symbol.BuiltinValueInteger)
	case "function":
		return newBuiltinType(symbol.BuiltinValueFunction)
	case "table":
		return newBuiltinType(symbol.BuiltinValueTable)
	case "thread":
		return newBuiltinType(symbol.BuiltinValueTable)
	case "userdata":
		return newBuiltinType(symbol.BuiltinValueUserdata)
	case "lightuserdata":
		return newBuiltinType(symbol.BuiltinValueLightuserdata)
	}

	return nil
}

/* func (l *Lexer) lexTypeReference(source string) {
	typeAnnotation := fmt.Sprintf("---@type %s", source)
	_ = l.scratch([]string{typeAnnotation})

	locations, _ := l.context.Nvim.GetDefinitionLocation(0, uint(len(typeAnnotation)))

	fmt.Printf("%+v", locations)

} */

func (l *Lexer) lexReference(name string) error {
	if _, exists := l.context.Result().Types[name]; exists {
		l.context.Logger.Info("Skipping '%s': lexed type already found")
		return nil
	}

	l.context.Result().Types[name] = struct{}{}

	l.context.Fork(name, func(name string) error {
		typeSource, err := l.crawler.SourceType(name)

		if err != nil {
			fmt.Printf("Reference error: %v", err)
		}

		fmt.Printf("\nReference source:\n%+v\n\n", typeSource.Origin)

		aliasType, err := l.lexAlias(typeSource)

		fmt.Printf("\nType: %+v\n\n", typeSource.Origin.Definition)

		switch {
		case err != nil:
			return err
		case aliasType != nil:
			l.context.Result().Types[name] = aliasType
			return nil
		}

		l.context.Logger.Warn(fmt.Sprintf("Uknown type '%s' received", name))
		l.context.Result().Types[name] = newUnknownType()
		return nil
	})

	return nil
}

func (l *Lexer) lexType(source string) (symbol.Symbol, error) {
	builtinType := l.lexBuiltinType(strings.TrimSpace(source))

	if builtinType != nil {
		return builtinType, nil
	}

	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	}

	defer buffer.Delete()

	typeAnnotation := fmt.Sprintf("---@type %s", source)
	err = buffer.SetLines([]string{typeAnnotation})

	if err != nil {
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	}

	functionType := newFunctionType()
	err = l.lexFunctionType(functionType)

	switch {
	case errors.Is(ErrNoMatch, err):
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	default:
		return functionType, nil
	}

	tableType := newTableType()
	err = l.lexTableType(tableType)

	switch {
	case errors.Is(ErrNoMatch, err):
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	default:
		return tableType, nil
	}

	l.lexReference(source)

	l.context.Logger.Warn(fmt.Sprintf("Uknown type '%s' received", source))
	return newUnknownType(), nil
}
