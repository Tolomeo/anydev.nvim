package lex

import (
	"errors"
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

/*
	 func newReferenceType(value string) lexed.Reference {
		return lexed.Reference{
			Type:  "reference",
			Value: value,
		}
	}
*/

func newFunctionType() lexed.Function {
	return lexed.Function{
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

func newUnknownType() lexed.Unknown {
	return lexed.Unknown{
		Kind: lexed.UnknownKindUnknown,
	}
}

func newTableType() lexed.Table {
	return lexed.Table{
		Kind: lexed.TableKindTable,
	}
}

var NoMatch = errors.New("No match")
var SyntaxError = errors.New("Syntax error")

var typeFunctionQuery string = `
	(documentation
		(type_annotation
			(function_type
				(parameter
					(identifier) @parameter.name
					[
						(builtin_type)
						(table_type)
						(member_type)
					] @parameter.type
				) @parameter
				(builtin_type)? @return.type
			)
		)
	)
	(ERROR) @error
`

func (l *lexer) lexFunctionType(function *lexed.Function) error {
	captures, err := l.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: typeFunctionQuery})

	switch {
	case errors.Is(nvim.ErrNotFound, err):
		return NoMatch
	case err != nil:
		return err
	}

	lexedFunction := newFunctionType()

	for _, capture := range captures {
		switch capture.Id {
		case "parameter":
			lexedFunction.Args = append(lexedFunction.Args, lexed.FunctionArg{})
		case "parameter.name":
			lexedFunction.Args[len(lexedFunction.Args)-1].Name = capture.Node.Text
		case "parameter.type":
			parameterType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return err
			}

			lexedFunction.Args[len(lexedFunction.Args)-1].Type = parameterType
		case "return.type":
			//TODO
		case "error":
			//TODO: trace error
			return SyntaxError
		}
	}

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
	(ERROR) @error
`

func (l *lexer) lexTableType(table *lexed.Table) error {
	captures, err := l.nvim.TsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: typeTableQuery})

	switch {
	case errors.Is(nvim.ErrNotFound, err):
		return NoMatch
	case err != nil:
		return err
	}

	fmt.Printf("\n\n%+v\n\n", captures)

	for _, capture := range captures {
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
		case "documentation":
			table.Fields[len(table.Fields)-1].Documentation = []string{capture.Node.Text}
		case "error":
			// TODO: trace that there was an error while lexing the table
			return SyntaxError
		}
	}

	return nil
}

func (l *lexer) lexType(source string) (lexed.Symbol, error) {
	switch source {
	case "void":
		return newBuiltinType(lexed.BuiltinValueVoid), nil
	case "nil":
		return newBuiltinType(lexed.BuiltinValueNil), nil
	case "any":
		return newBuiltinType(lexed.BuiltinValueAny), nil
	case "boolean":
		return newBuiltinType(lexed.BuiltinValueBoolean), nil
	case "string":
		return newBuiltinType(lexed.BuiltinValueString), nil
	case "number":
		return newBuiltinType(lexed.BuiltinValueNumber), nil
	case "integer", "int":
		return newBuiltinType(lexed.BuiltinValueInteger), nil
	case "function":
		return newBuiltinType(lexed.BuiltinValueFunction), nil
	case "table":
		return newBuiltinType(lexed.BuiltinValueTable), nil
	case "thread":
		return newBuiltinType(lexed.BuiltinValueTable), nil
	case "userdata":
		return newBuiltinType(lexed.BuiltinValueUserdata), nil
	case "lightuserdata":
		return newBuiltinType(lexed.BuiltinValueLightuserdata), nil
	}

	lines := []string{"@type " + source}
	err := l.scratch(lines)

	if err != nil {
		return struct{}{}, fmt.Errorf("Error lexing type %s: %w", source, err)
	}

	functionType := newFunctionType()
	err = l.lexFunctionType(&functionType)

	switch {
	case errors.Is(NoMatch, err):
	case err != nil:
		return struct{}{}, fmt.Errorf("Error lexing type %s: %w", source, err)
	default:
		return functionType, nil
	}

	tableType := newTableType()
	err = l.lexTableType(&tableType)

	switch {
	case errors.Is(NoMatch, err):
	case err != nil:
		return struct{}{}, fmt.Errorf("Error lexing type %s: %w", source, err)
	default:
		return tableType, nil
	}

	return newUnknownType(), nil
}
