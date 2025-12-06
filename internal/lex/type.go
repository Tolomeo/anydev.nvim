package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
)

/*
	 func newReference(value string) lexed.Reference {
		return lexed.Reference{
			Type:  "reference",
			Value: value,
		}
	}
*/
func newBuiltin(value lexed.BuiltinValue) lexed.Builtin {
	return lexed.Builtin{
		Type:  "builtin",
		Value: value,
	}
}

func newUnknown() lexed.Unknown {
	return lexed.Unknown {
		Type: "uknown",
	}
}

func (l *lexer) lexType(source string) (lexed.Symbol, error) {
	fmt.Println(source)

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

	err := l.scratch([]string{"@type " + source})

	if err != nil {
		return struct{}{}, fmt.Errorf("Error lexing type %s: %w", source, err)
	}

	return newUnknown(), nil
}
