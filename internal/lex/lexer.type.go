package lex

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

var ErrNoMatch = errors.New("No match")

func newReferenceType(value string) *symbol.Reference {
	return &symbol.Reference{
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

var typeOptionalQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(documentation
		(type_annotation 
			(optional_type
				(_) @optional.type
			) @optional
		)
	)
	`,
}

func newOptional(typ *symbol.Symbol) *symbol.Optional {
	return &symbol.Optional{
		Kind: symbol.OptionalKindOptional,
		Type: typ,
	}
}

func (l *Lexer) lexOptional(buffer *nvim.Buffer, source string) (*symbol.Optional, error) {
	match, err := buffer.TsQueryOne(typeOptionalQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	for _, capture := range *match {
		switch capture.Id {
		case "optional.type":
			lexedType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return nil, err
			}

			return newOptional(&lexedType), nil
		}
	}

	return nil, fmt.Errorf("Could not retrieved type from optional type '%s'", source)
}

var typeFunctionQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
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
				(":"
					(%s) @return.type
					("," (%s) @return.type)*
				)? @return
			)
		)
	)
`, anyTypeQuery, anyTypeQuery, anyTypeQuery, anyTypeQuery, anyTypeQuery, anyTypeQuery),
}

func (l *Lexer) lexFunctionType(buffer *nvim.Buffer, function *symbol.Function) error {
	captures, err := buffer.TsQueryOne(typeFunctionQuery)

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

func (l *Lexer) lexTableType(buffer *nvim.Buffer, table *symbol.Table) error {
	matches, err := buffer.TsQueryAll(treesitter.Query{Language: "luadoc", Query: typeTableQuery})

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

func (l *Lexer) lexReference(name string) (*symbol.Reference, error) {
	if _, alreadyLexed := l.context.Result().Types[name]; alreadyLexed {
		l.context.Logger.Info(fmt.Sprintf("Skipping '%s': lexed type already found", name))
		return newReferenceType(name), nil
	}

	l.context.Result().Types[name] = newUnknownType()

	err := l.context.Fork(name, func(name string) error {
		typeSource, err := l.crawler.SourceType(name)

		if err != nil {
			return err
		}

		// fmt.Printf("\nReference '%s' source:\n%+v\n\n", name, typeSource.Origin)

		aliasType, err := l.lexAlias(typeSource)

		// fmt.Printf("\nLexed '%s' alias: %+v\n\n", name, aliasType)

		// fmt.Printf("\nType: %+v\n\n", typeSource.Origin.Definition)

		switch {
		case err != nil:
			return err
		case aliasType != nil:
			l.context.Result().Types[name] = aliasType
			return nil
		}

		l.context.Logger.Warn(fmt.Sprintf("Uknown type '%s' received", name))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return newReferenceType(name), nil
}

var typeArrayQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(documentation
		(type_annotation
			(array_type 
				(_) @array.itemstype 
			) @array
		)
	)
`}

func newArrayType(items symbol.Symbol) *symbol.Array {
	return &symbol.Array{
		Kind:  symbol.ArrayKindArray,
		Items: items,
	}
}

func (l *Lexer) lexArray(buffer *nvim.Buffer, source string) (*symbol.Array, error) {
	match, err := buffer.TsQueryOne(typeArrayQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	for _, matchCapture := range *match {
		switch matchCapture.Id {
		case "array.itemstype":
			itemsType, err := l.lexType(matchCapture.Node.Text)

			if err != nil {
				return nil, err
			}

			return newArrayType(itemsType), nil
		}
	}

	return nil, fmt.Errorf("Could not retrieve items type value from the array type '%s'", source)
}

var typeUnionQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(type_annotation
			(union_type
				(%s) @union.type
				(%s) @union.type
			) @union
		)
	)`, anyTypeQuery, anyTypeQuery)}

func newUnionType(types []symbol.Symbol) *symbol.Union {
	unionTypes := []symbol.UnionTypesElem{}

	for _, typ := range types {
		unionTypes = append(unionTypes, typ)
	}

	return &symbol.Union{
		Kind:  symbol.UnionKindUnion,
		Types: unionTypes,
	}
}

func (l *Lexer) lexUnion(buffer *nvim.Buffer, source string) (*symbol.Union, error) {
	match, err := buffer.TsQueryOne(typeUnionQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	unionTypes := []symbol.Symbol{}

	for _, matchCapture := range *match {
		switch matchCapture.Id {
		case "union.type":
			lexedType, err := l.lexType(matchCapture.Node.Text)

			if err != nil {
				return nil, err
			}

			// Flattening nested unions
			switch t := lexedType.(type) {
			case symbol.Union:
				for _, lexedUnionType := range t.Types {
					unionTypes = append(unionTypes, lexedUnionType)
				}
			default:
				unionTypes = append(unionTypes, lexedType)
			}
		}
	}

	if len(unionTypes) < 2 {
		return nil, fmt.Errorf("Could not retrieve all types in the union type '%s'", source)
	}

	return newUnionType(unionTypes), nil
}

var typeGroupQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(documentation
		(type_annotation
			(parenthesized_type 
				(_) @group.type
			) @group
		)
	)
`}

func (l *Lexer) lexGroup(buffer *nvim.Buffer, source string) (symbol.Symbol, error) {
	match, err := buffer.TsQueryOne(typeGroupQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	for _, matchCapture := range *match {
		switch matchCapture.Id {
		case "group.type":
			return l.lexType(matchCapture.Node.Text)
		}
	}

	return nil, fmt.Errorf("Could not retrieve the type value of the grouped type '%s'", source)
}

var typeStringLiteralQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(documentation
		(type_annotation
			(literal_type) @stringliteral
		)
	)
`}

func newStringLiteral(value string) *symbol.StringLiteral {
	return &symbol.StringLiteral{
		Kind:  symbol.StringLiteralKindStringliteral,
		Value: value,
	}
}

func (l *Lexer) lexStringLiteral(buffer *nvim.Buffer, source string) (*symbol.StringLiteral, error) {
	match, err := buffer.TsQueryOne(typeStringLiteralQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	for _, matchCapture := range *match {
		switch matchCapture.Id {
		case "stringliteral":
			return newStringLiteral(matchCapture.Node.Text), nil
		}
	}

	return nil, fmt.Errorf("Could not retrieve the value of the string literal type '%s'", source)
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

	defer buffer.Close()

	typeAnnotation := fmt.Sprintf("---@type %s", source)
	err = buffer.SetLines([]string{typeAnnotation})

	if err != nil {
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	}

	functionType := newFunctionType()
	err = l.lexFunctionType(buffer, functionType)

	switch {
	case errors.Is(ErrNoMatch, err):
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	default:
		return functionType, nil
	}

	lexedArray, err := l.lexArray(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedArray != nil:
		return lexedArray, nil
	}

	tableType := newTableType()
	err = l.lexTableType(buffer, tableType)

	switch {
	case errors.Is(ErrNoMatch, err):
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	default:
		return tableType, nil
	}

	lexedOptional, err := l.lexOptional(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedOptional != nil:
		return lexedOptional, nil
	}

	lexedUnion, err := l.lexUnion(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedUnion != nil:
		return lexedUnion, nil
	}

	lexedGroupedType, err := l.lexGroup(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedGroupedType != nil:
		return lexedGroupedType, nil
	}

	lexedStringLiteral, err := l.lexStringLiteral(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedStringLiteral != nil:
		return lexedStringLiteral, nil
	}

	lexedReference, err := l.lexReference(source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedReference != nil:
		return lexedReference, nil
	}

	l.context.Logger.Warn(fmt.Sprintf("Unknown type '%s' received", source))
	return newUnknownType(), nil
}
