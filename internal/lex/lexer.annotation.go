package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

type TypeAnnotation struct {
	Name string
}

type AtTypeAnnotation struct {
	Types         []TypeAnnotation
	Documentation []string
}

var atTypeAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(type_annotation
			"@type" . (%s) @type.type
			.
			("," . (%s) @type.type)*
			.
			(comment)? @type.documentation
			.
		) @type
	)`, anyTypeAnnotationQuery, anyTypeAnnotationQuery),
}

func (l *Lexer) lexAtTypeAnnotations(buffer *nvim.ScratchBuffer) (*AtTypeAnnotation, error) {
	captures, err := buffer.TsQueryOne(atTypeAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if captures == nil {
		return nil, nil
	}

	atType := AtTypeAnnotation{}

	for _, capture := range *captures {
		switch capture.Id {
		case "type.type":
			atType.Types = append(atType.Types, TypeAnnotation{capture.Node.Text})
		case "type.documentation":
			atType.Documentation = strings.Split(capture.Node.Text, "\n")
		}
	}

	return &atType, nil
}

type AtOverloadAnnotation struct {
	Type          TypeAnnotation
	Documentation []string
}

var atOverloadAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation 
		(overload_annotation 
			"@overload"
			.
			(%s) @type 
			.
			(comment)? @documentation
			.
		)
	)
`, typeAnnotationQueries["function_type"]),
}

func (l *Lexer) lexAtOverloadAnnotations(buffer *nvim.ScratchBuffer) ([]AtOverloadAnnotation, error) {
	atOverloads := []AtOverloadAnnotation{}
	matches, err := buffer.SafeTsQueryAll(atOverloadAnnotationQuery)

	// fmt.Printf("\n Overload matches: %+v\n", matches)

	if err != nil {
		return atOverloads, err
	}

	if matches == nil {
		return atOverloads, nil
	}

	for _, match := range *matches {
		if match.HasError {
			l.context.Logger().Warn(fmt.Sprintf("Skipping overload annotation in '%s' because it contains syntax errors", l.context.Target().Name()))
			continue
		}

		overload := AtOverloadAnnotation{}

		for _, capture := range match.Captures {
			switch capture.Id {
			case "documentation":
				overload.Documentation = []string{capture.Node.Text}
			case "type":
				overload.Type = TypeAnnotation{capture.Node.Text}
			}
		}

		atOverloads = append(atOverloads, overload)
	}

	return atOverloads, nil
}

type AtGenericAnnotation struct {
	Name  string
	Types []TypeAnnotation
}

var atGenericAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(generic_annotation
			"@generic"
			.
			(identifier) @generic.name
			.
			(":"
				.
				parent_type:
					(%s) @generic.type
			)?
			.
			(","
				.
				(identifier) @generic.name
				.
				(":"
					.
					parent_type:
						(%s) @generic.type
				)?
			)*
		) @generic
	)`, anyTypeAnnotationQuery, anyTypeAnnotationQuery),
}

func (l *Lexer) lexAtGenericAnnotations(buffer *nvim.ScratchBuffer) ([]AtGenericAnnotation, error) {
	atGenericAnnotations := []AtGenericAnnotation{}
	matches, err := buffer.TsQueryAll(atGenericAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if matches == nil {
		return atGenericAnnotations, nil
	}

	for _, matchCaptures := range *matches {
		lexedGeneric := AtGenericAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "generic.name":
				lexedGeneric.Name = capture.Node.Text
			case "generic.type":
				lexedGeneric.Types = append(lexedGeneric.Types, TypeAnnotation{capture.Node.Text})
			}
		}

		if lexedGeneric.Name == "" {
			return atGenericAnnotations, fmt.Errorf("Could not retrieve generic name for generic annotation '%v'", lexedGeneric)
		}

		atGenericAnnotations = append(atGenericAnnotations, lexedGeneric)
	}

	return atGenericAnnotations, nil
}

type AtParamAnnotation struct {
	Name          string
	Type          TypeAnnotation
	Optional      bool
	Documentation []string
}

var atParamAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(param_annotation
			"@param"
			.
			([("...") (identifier)]) @name
			.
			"?"? @optional
			.
			(%s) @type
			.
			(comment)? @documentation
			.
		) @param
	)`, anyTypeAnnotationQuery),
}

func (l *Lexer) lexAtParamAnnotations(buffer *nvim.ScratchBuffer) (map[string]AtParamAnnotation, error) {
	params := map[string]AtParamAnnotation{}
	matches, err := buffer.TsQueryAll(atParamAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if matches == nil {
		return params, nil
	}

	for _, matchCaptures := range *matches {
		lexedParam := AtParamAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "name":
				lexedParam.Name = capture.Node.Text
			case "optional":
				lexedParam.Optional = true
			case "documentation":
				lexedParam.Documentation = []string{capture.Node.Text}
			case "type":
				lexedParam.Type = TypeAnnotation{capture.Node.Text}
			}
		}

		if lexedParam.Name == "" {
			return params, fmt.Errorf("Could not retrieve param name for param annotation")
		}

		params[lexedParam.Name] = lexedParam
	}

	return params, nil
}

type AtReturnAnnotation struct {
	Name          string
	Type          TypeAnnotation
	Documentation []string
}

var atReturnAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(return_annotation
			"@return"
			.
			(%s) @return.type
			.
			(comment)? @return.documentation
			.
		) @return
	)`, anyTypeAnnotationQuery),
}

func (l *Lexer) lexAtReturnAnnotations(buffer *nvim.ScratchBuffer) ([]AtReturnAnnotation, error) {
	atReturns := []AtReturnAnnotation{}
	matches, err := buffer.TsQueryAll(atReturnAnnotationQuery)

	if err != nil {
		return atReturns, err
	}

	if matches == nil {
		return atReturns, nil
	}

	for _, matchCaptures := range *matches {
		returnAnnotation := AtReturnAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "return.name":
				returnAnnotation.Name = capture.Node.Text
			case "return.documentation":
				returnAnnotation.Documentation = []string{capture.Node.Text}
			case "return.type":
				returnAnnotation.Type = TypeAnnotation{capture.Node.Text}
			}
		}

		atReturns = append(atReturns, returnAnnotation)
	}

	return atReturns, nil
}

var typeAnnotationQueries = map[string]string{
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

var anyTypeAnnotationQuery = fmt.Sprintf(`[%s]`, strings.Join(mapx.Values(typeAnnotationQueries), " "))

var optionalTypeAnnotationQuery = treesitter.Query{
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

func (l *Lexer) lexOptionalTypeAnnotation(buffer *nvim.ScratchBuffer, source string) (annotation.Type, error) {
	match, err := buffer.TsQueryOne(optionalTypeAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if match == nil {
		return nil, nil
	}

	for _, capture := range *match {
		switch capture.Id {
		case "optional.type":
			optionalType, err := l.lexTypeAnnotation(TypeAnnotation{capture.Node.Text})

			if err != nil {
				return nil, err
			}

			return symbol.NewOptional(optionalType), nil
		}
	}

	return nil, fmt.Errorf("Could not retrieved type from optional type '%s'", source)
}

var functionTypeAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(type_annotation
			(function_type
				(parameter
					(identifier) @parameter.name
					":"
					(%s) @parameter.type
				)? @parameter
				("," (parameter
					(identifier) @parameter.name
					":"
					(%s) @parameter.type
				) @parameter)*
				("," (parameter
					"..." @parameter.name
					":"
					(%s) @parameter.type
				) @parameter)?
				(parameter
					"..." @parameter.name
					":"
					(%s) @parameter.type
				)? @parameter
				(":"
					(%s) @return.type
					("," (%s) @return.type)*
				)? @return
				(comment)? @function.documentation
			)
		)
	)
`, anyTypeAnnotationQuery, anyTypeAnnotationQuery, anyTypeAnnotationQuery, anyTypeAnnotationQuery, anyTypeAnnotationQuery, anyTypeAnnotationQuery),
}

func (l *Lexer) lexFunctionTypeAnnotation(buffer *nvim.ScratchBuffer) (*symbol.Function, error) {
	captures, err := buffer.TsQueryOne(functionTypeAnnotationQuery)

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	function := symbol.NewFunction()
	args := []symbol.FunctionArgument{}
	returns := []symbol.FunctionReturn{}

	for _, capture := range *captures {
		switch capture.Id {
		case "parameter":
			args = append(args, symbol.FunctionArgument{})
		case "parameter.name":
			args[len(args)-1].Name = capture.Node.Text
		case "parameter.type":
			parameterType, err := l.lexTypeAnnotation(TypeAnnotation{capture.Node.Text})

			if err != nil {
				return nil, err
			}

			args[len(args)-1].Type = parameterType
		case "return.type":
			returnType, err := l.lexTypeAnnotation(TypeAnnotation{capture.Node.Text})

			if err != nil {
				return nil, err
			}

			returns = append(returns, symbol.FunctionReturn{
				Type: returnType,
			})
		}
	}

	function.Arguments = append(function.Arguments, args...)
	function.Returns = append(function.Returns, returns...)

	return function, nil
}

// TODO: check if it is possible to mark value as optional
var tableTypeAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(type_annotation
			(table_type
				key: (%s) @key
				value: (%s) @value
			) @table
			(comment)? @documentation
		)
	)`, anyTypeAnnotationQuery, anyTypeAnnotationQuery),
}

func (l *Lexer) lexTableTypeAnnotation(buffer *nvim.ScratchBuffer) (*symbol.Table, error) {
	match, err := buffer.TsQueryOne(tableTypeAnnotationQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	table := symbol.NewTable()

	for _, capture := range *match {
		switch capture.Id {
		case "table":
			table.Indexes = append(table.Indexes, *symbol.NewTableIndex())
		case "key":
			keyType, err := l.lexTypeAnnotation(TypeAnnotation{capture.Node.Text})

			if err != nil {
				return nil, err
			}

			table.Indexes[len(table.Indexes)-1].Key = keyType
		case "value":
			valueType, err := l.lexTypeAnnotation(TypeAnnotation{capture.Node.Text})

			if err != nil {
				return nil, err
			}

			table.Indexes[len(table.Indexes)-1].Value = valueType
		}
	}

	return table, nil
}

func (l *Lexer) lexBuiltinTypeAnnotation(source string) annotation.Type {
	switch source {
	case "void":
		return symbol.NewVoid()
	case "nil":
		return symbol.NewNil()
	case "any":
		return symbol.NewAny()
	case "boolean":
		return symbol.NewBoolean()
	case "string":
		return symbol.NewString()
	case "number":
		return symbol.NewNumber()
	case "integer", "int":
		return symbol.NewInteger()
	case "function":
		return symbol.NewBuiltinFunction()
	case "table":
		return symbol.NewBuiltinTable()
	case "thread":
		return symbol.NewThread()
	case "userdata":
		return symbol.NewUserdata()
	case "lightuserdata":
		return symbol.NewLightUserdata()
	}

	return nil
}

func (l *Lexer) lexReferenceTypeAnnotation(name string) (*symbol.Reference, error) {
	err := l.context.Extract("type", name)

	if err != nil {
		return nil, err
	}

	return symbol.NewReference(name), nil
}

var arrayTypeAnnotationQuery = treesitter.Query{
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

func (l *Lexer) lexArrayTypeAnnotation(buffer *nvim.ScratchBuffer, source string) (*symbol.Array, error) {
	match, err := buffer.TsQueryOne(arrayTypeAnnotationQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	for _, matchCapture := range *match {
		switch matchCapture.Id {
		case "array.itemstype":
			itemsType, err := l.lexTypeAnnotation(TypeAnnotation{matchCapture.Node.Text})

			if err != nil {
				return nil, err
			}

			return symbol.NewArray(itemsType), nil
		}
	}

	return nil, fmt.Errorf("Could not retrieve items type value from the array type '%s'", source)
}

var literalTableAnnotationQueries = map[string]treesitter.Query{
	"empty": {
		Language: annotation.LiteralTableTypeEmptyQuery.Language,
		Query: fmt.Sprintf(`
		(documentation
			(type_annotation
				(%s)
				(comment)? @table.documentation
			)
		)`, annotation.LiteralTableTypeEmptyQuery.Query),
	},
	"described": {
		Language: annotation.LiteralTableTypeQuery.Language,
		Query: fmt.Sprintf(`
		(documentation
			(type_annotation
				(%s)
				(comment)? @table.documentation
			)
		)`, annotation.LiteralTableTypeQuery.Query),
	},
}

func (l *Lexer) lexLiteralTableTypeAnnotation(buffer *nvim.ScratchBuffer) (*symbol.Table, error) {
	for _, query := range literalTableAnnotationQueries {
		match, err := buffer.TsQueryOne(query)

		if err != nil {
			return nil, err
		}

		if match == nil {
			continue
		}

		literalTableType := annotation.NewLiteralTableType(*match)
		tableSymbol := symbol.NewTable()

		for _, fieldType := range literalTableType.Fields() {
			tableField := symbol.NewTableField()
			tableField.Name = fieldType.Key()
			tableFieldType, err := l.lexTypeAnnotation(TypeAnnotation{fieldType.Value()})

			if err != nil {
				return nil, err
			}

			tableField.Type = tableFieldType
			tableSymbol.Fields = append(tableSymbol.Fields, *tableField)
		}

		for _, indexType := range literalTableType.Indexes() {
			tableIndex := symbol.NewTableIndex()
			tableIndexKey, err := l.lexTypeAnnotation(TypeAnnotation{indexType.Key()})

			if err != nil {
				return nil, err
			}

			tableIndexValue, err := l.lexTypeAnnotation(TypeAnnotation{indexType.Value()})

			if err != nil {
				return nil, err
			}

			tableIndex.Key = tableIndexKey
			tableIndex.Value = tableIndexValue

			tableSymbol.Indexes = append(tableSymbol.Indexes, *tableIndex)
		}

		return tableSymbol, nil
	}

	return nil, nil
}

var unionTypeAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(type_annotation
			(union_type
				(%s) @union.type
				(%s) @union.type
			) @union
		)
	)`, anyTypeAnnotationQuery, anyTypeAnnotationQuery)}

func (l *Lexer) lexUnionTypeAnnotation(buffer *nvim.ScratchBuffer, source string) (*symbol.Union, error) {
	match, err := buffer.TsQueryOne(unionTypeAnnotationQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	unionTypes := []annotation.Type{}

	for _, matchCapture := range *match {
		switch matchCapture.Id {
		case "union.type":
			lexedType, err := l.lexTypeAnnotation(TypeAnnotation{matchCapture.Node.Text})

			if err != nil {
				return nil, err
			}

			// Flattening nested unions
			switch t := lexedType.(type) {
			case *symbol.Union:
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

	return symbol.NewUnion(unionTypes), nil
}

var parenthesizedTypeAnnotationQuery = treesitter.Query{
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

func (l *Lexer) lexParenthesizedTypeAnnotation(buffer *nvim.ScratchBuffer, source string) (annotation.Type, error) {
	match, err := buffer.TsQueryOne(parenthesizedTypeAnnotationQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	for _, matchCapture := range *match {
		switch matchCapture.Id {
		case "group.type":
			return l.lexTypeAnnotation(TypeAnnotation{matchCapture.Node.Text})
		}
	}

	return nil, fmt.Errorf("Could not retrieve the type value of the grouped type '%s'", source)
}

var literalNumberTypeAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(documentation
		(type_annotation
			(numeric_literal_type) @literal.number
		) @literal
	)`,
}

func (l *Lexer) lexLiteralNumberTypeAnnotation(buffer *nvim.ScratchBuffer, source string) (*symbol.NumericLiteral, error) {
	match, err := buffer.TsQueryOne(literalNumberTypeAnnotationQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	for _, matchCapture := range *match {
		switch matchCapture.Id {
		case "literal.number":
			return symbol.NewNumericLiteral(matchCapture.Node.Text), nil
		}
	}

	return nil, fmt.Errorf("Could not retrieve the value of the string literal type '%s'", source)
}

var literalBooleanTypeAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(documentation
		(type_annotation
			(identifier) @literal.value
		) @literal
		(#any-of? @literal.value "true" "false")
	)`,
}

func (l *Lexer) lexLiteralBooleanTypeAnnotation(buffer *nvim.ScratchBuffer) (*symbol.BooleanLiteral, error) {
	match, err := buffer.TsQueryOne(literalBooleanTypeAnnotationQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	for _, matchCapture := range *match {
		switch matchCapture.Id {
		case "literal.value":
			return symbol.NewBooleanLiteral(matchCapture.Node.Text), nil
		}
	}

	return nil, fmt.Errorf("Could not retrieve the value of the string literal type")
}

var literalStringTypeAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(documentation
		(type_annotation
			(literal_type) @stringliteral
		)
	)
`}

func (l *Lexer) lexLiteralStringTypeAnnotation(buffer *nvim.ScratchBuffer, source string) (*symbol.StringLiteral, error) {
	match, err := buffer.TsQueryOne(literalStringTypeAnnotationQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	for _, matchCapture := range *match {
		switch matchCapture.Id {
		case "stringliteral":
			return symbol.NewStringLiteral(matchCapture.Node.Text), nil
		}
	}

	return nil, fmt.Errorf("Could not retrieve the value of the string literal type '%s'", source)
}

func (l *Lexer) lexTypeAnnotation(typ TypeAnnotation) (annotation.Type, error) {
	source := strings.TrimSpace(typ.Name)
	builtinType := l.lexBuiltinTypeAnnotation(strings.TrimSpace(source))

	if builtinType != nil {
		return builtinType, nil
	}

	buffer, err := l.context.Nvim().NewBuffer()

	if err != nil {
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	}

	defer buffer.Close()

	typeAnnotation := fmt.Sprintf("---@type %s", source)
	err = buffer.SetLines([]string{typeAnnotation})

	if err != nil {
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	}

	functionType, err := l.lexFunctionTypeAnnotation(buffer)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case functionType != nil:
		return functionType, nil
	}

	lexedArray, err := l.lexArrayTypeAnnotation(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedArray != nil:
		return lexedArray, nil
	}

	tableType, err := l.lexTableTypeAnnotation(buffer)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case tableType != nil:
		return tableType, nil
	}

	literalTableType, err := l.lexLiteralTableTypeAnnotation(buffer)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case literalTableType != nil:
		return literalTableType, nil
	}

	lexedOptional, err := l.lexOptionalTypeAnnotation(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedOptional != nil:
		return lexedOptional, nil
	}

	lexedUnion, err := l.lexUnionTypeAnnotation(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedUnion != nil:
		return lexedUnion, nil
	}

	lexedGroupedType, err := l.lexParenthesizedTypeAnnotation(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedGroupedType != nil:
		return lexedGroupedType, nil
	}

	lexedNumericLiteral, err := l.lexLiteralNumberTypeAnnotation(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedNumericLiteral != nil:
		return lexedNumericLiteral, nil
	}

	lexedBooleanLiteral, err := l.lexLiteralBooleanTypeAnnotation(buffer)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedBooleanLiteral != nil:
		return lexedBooleanLiteral, nil
	}

	lexedStringLiteral, err := l.lexLiteralStringTypeAnnotation(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedStringLiteral != nil:
		return lexedStringLiteral, nil
	}

	lexedReference, err := l.lexReferenceTypeAnnotation(source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedReference != nil:
		return lexedReference, nil
	}

	l.context.Logger().Warn(fmt.Sprintf("Unknown type '%s' received", source))
	return symbol.NewUnknown(), nil
}
