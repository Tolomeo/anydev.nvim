package lex

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

/*
	 func newReference(value string) lexed.Reference {
		return lexed.Reference{
			Type:  "reference",
			Value: value,
		}
	}
*/
func newFunction() lexed.Function {
	return lexed.Function{
		Type: "function",
	}
}

func newBuiltin(value lexed.BuiltinValue) lexed.Builtin {
	return lexed.Builtin{
		Type:  "builtin",
		Value: value,
	}
}

func newUnknown() lexed.Unknown {
	return lexed.Unknown{
		Type: "uknown",
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

	lexedFunction := newFunction()

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
		return newBuiltin(lexed.BuiltinValueNil), nil
	case "any":
		return newBuiltin(lexed.BuiltinValueAny), nil
	case "boolean":
		return newBuiltin(lexed.BuiltinValueBoolean), nil
	case "string":
		return newBuiltin(lexed.BuiltinValueString), nil
	case "number":
		return newBuiltin(lexed.BuiltinValueNumber), nil
	case "integer", "int":
		return newBuiltin(lexed.BuiltinValueInteger), nil
	case "function":
		return newBuiltin(lexed.BuiltinValueFunction), nil
	case "table":
		return newBuiltin(lexed.BuiltinValueTable), nil
	case "thread":
		return newBuiltin(lexed.BuiltinValueTable), nil
	case "userdata":
		return newBuiltin(lexed.BuiltinValueUserdata), nil
	case "lightuserdata":
		return newBuiltin(lexed.BuiltinValueLightuserdata), nil
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

	return newUnknown(), nil
}
