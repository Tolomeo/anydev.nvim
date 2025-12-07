package lex

import (
	"errors"
	"fmt"
	"strings"

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
		Kind: "function",
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
		Kind:  "builtin",
		Value: value,
	}
}

func newUnknownType() lexed.Unknown {
	return lexed.Unknown{
		Kind: "uknown",
	}
}

var NoMatch = errors.New("No match")

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
`

func (l *lexer) lexFunctionType() (lexed.Function, error) {
	captures, err := l.nvim.TsQuery("luadoc", typeFunctionQuery)

	switch {
	case errors.Is(nvim.ErrNotFound, err):
		return lexed.Function{}, NoMatch
	case err != nil:
		return lexed.Function{}, err
	}

	lexedFunction := newFunctionType()

	for _, capture := range captures {
		switch capture.Id {
		case "parameter":
			lexedFunction.Args = append(lexedFunction.Args, lexed.FunctionArg{})
		case "parameter.name":
			lexedFunction.Args[len(lexedFunction.Args)-1].Name = strings.Join(capture.Node.Text, "")
		case "parameter.type":
			parameterType, err := l.lexType(strings.Join(capture.Node.Text, ""))

			if err != nil {
				return lexed.Function{}, err
			}

			lexedFunction.Args[len(lexedFunction.Args)-1].Type = parameterType
		case "return.type":
			//TODO
		}
	}

	return lexedFunction, nil
}

func (l *lexer) lexType(source string) (lexed.Symbol, error) {
	switch source {
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

	functionType, err := l.lexFunctionType()

	switch {
	case errors.Is(NoMatch, err):
	case err != nil:
		return struct{}{}, fmt.Errorf("Error lexing type %s: %w", source, err)
	default:
		return functionType, nil
	}

	return newUnknownType(), nil
}
