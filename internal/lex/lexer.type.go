package lex

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

var ErrNoMatch = errors.New("No match")

func newReferenceType(value string) lexed.Reference {
	return lexed.Reference{
		Kind:  lexed.ReferenceKindReference,
		Value: value,
	}
}

func newFunctionType() *lexed.Function {
	return &lexed.Function{
		Kind: lexed.FunctionKindFunction,
	}
}

func newFunctionTypeArg(name string) lexed.FunctionArg {
	return lexed.FunctionArg{
		Name:     name,
		Type:     newUnknownType(),
		Optional: false,
	}
}

func newBuiltinType(value lexed.BuiltinValue) lexed.Builtin {
	return lexed.Builtin{
		Kind:  lexed.BuiltinKindBuiltin,
		Value: value,
	}
}

func newUnknownType() *lexed.Unknown {
	return &lexed.Unknown{
		Kind: lexed.UnknownKindUnknown,
	}
}

func newTableType() *lexed.Table {
	return &lexed.Table{
		Kind: lexed.TableKindTable,
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

func (l *Lexer) lexFunctionType(function *lexed.Function) error {
	captures, err := l.context.nvim.TsQueryOne(nvim.TsQueryConfig{Language: "luadoc", Query: typeFunctionQuery})

	switch {
	case err != nil:
		return err
	case captures == nil:
		return ErrNoMatch
	}

	args := []lexed.FunctionArg{}
	returns := []lexed.FunctionReturn{}

	for _, capture := range *captures {
		switch capture.Id {
		case "parameter":
			args = append(args, lexed.FunctionArg{})
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

			returns = append(returns, lexed.FunctionReturn{
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

func (l *Lexer) lexTableType(table *lexed.Table) error {
	matches, err := l.context.nvim.TsQueryAll(nvim.TsQueryConfig{Language: "luadoc", Query: typeTableQuery})

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
				table.Fields = append(table.Fields, lexed.TableField{})
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

func (l *Lexer) lexBuiltinType(source string) lexed.Symbol {
	switch source {
	case "void":
		return newBuiltinType(lexed.BuiltinValueVoid)
	case "nil":
		return newBuiltinType(lexed.BuiltinValueNil)
	case "any":
		return newBuiltinType(lexed.BuiltinValueAny)
	case "boolean":
		return newBuiltinType(lexed.BuiltinValueBoolean)
	case "string":
		return newBuiltinType(lexed.BuiltinValueString)
	case "number":
		return newBuiltinType(lexed.BuiltinValueNumber)
	case "integer", "int":
		return newBuiltinType(lexed.BuiltinValueInteger)
	case "function":
		return newBuiltinType(lexed.BuiltinValueFunction)
	case "table":
		return newBuiltinType(lexed.BuiltinValueTable)
	case "thread":
		return newBuiltinType(lexed.BuiltinValueTable)
	case "userdata":
		return newBuiltinType(lexed.BuiltinValueUserdata)
	case "lightuserdata":
		return newBuiltinType(lexed.BuiltinValueLightuserdata)
	}

	return nil
}

/* func (l *Lexer) lexTypeReference(source string) {
	typeAnnotation := fmt.Sprintf("---@type %s", source)
	_ = l.scratch([]string{typeAnnotation})

	locations, _ := l.context.nvim.GetDefinitionLocation(0, uint(len(typeAnnotation)))

	fmt.Printf("%+v", locations)

} */

func (l *Lexer) lexType(source string) (lexed.Symbol, error) {
	builtinType := l.lexBuiltinType(strings.TrimSpace(source))

	if builtinType != nil {
		return builtinType, nil
	}

	buffer, err := l.context.nvim.Buffer()

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

	// l.lexTypeReference(source)

	l.context.logger.Warn(fmt.Sprintf("Uknown type '%s' received", source))
	return newUnknownType(), nil
}
