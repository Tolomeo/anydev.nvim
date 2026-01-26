package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/cache"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

var lexedAnnotationsCache = cache.NewCache[AtAnnotations]()

type TypeAnnotation struct {
	Name string
}

type AtTypeAnnotation struct {
	Types         []TypeAnnotation
	Documentation []string
}

type AtGenericAnnotation struct {
	Name  string
	Types []TypeAnnotation
}

type AtReturnAnnotation struct {
	Name          string
	Type          TypeAnnotation
	Documentation []string
}

type AtParamAnnotation struct {
	Name          string
	Type          TypeAnnotation
	Optional      bool
	Documentation []string
}

type AtOverloadAnnotation struct {
	Type          TypeAnnotation
	Documentation []string
}

type AtFieldAnnotation struct {
	Name          string
	Type          TypeAnnotation
	Optional      bool
	Private       bool
	Protected     bool
	Package       bool
	Documentation []string
}

type AtAnnotations struct {
	AtType       *AtTypeAnnotation
	AtPrivate    bool
	AtProtected  bool
	AtPackage    bool
	AtDeprecated bool
	AtGenerics   []AtGenericAnnotation
	AtParams     map[string]AtParamAnnotation
	AtReturns    []AtReturnAnnotation
	AtOverloads  []AtOverloadAnnotation
	AtFields     map[string]AtFieldAnnotation
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

func (l *Lexer) lexAtTypeAnnotations(buffer *nvim.ScratchBuffer, annotations *AtAnnotations) (bool, error) {
	captures, err := buffer.TsQueryOne(atTypeAnnotationQuery)

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
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

	annotations.AtType = &atType
	return true, nil
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

func (l *Lexer) lexAtOverloadAnnotations(buffer *nvim.ScratchBuffer, annotations *AtAnnotations) (bool, error) {
	matches, err := buffer.SafeTsQueryAll(atOverloadAnnotationQuery)

	// fmt.Printf("\n Overload matches: %+v\n", matches)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
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

		annotations.AtOverloads = append(annotations.AtOverloads, overload)
	}

	return true, nil
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

func (l *Lexer) lexAtGenericAnnotations(buffer *nvim.ScratchBuffer, annotations *AtAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(atGenericAnnotationQuery)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
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
			return false, fmt.Errorf("Could not retrieve generic name for generic annotation '%v'", lexedGeneric)
		}

		annotations.AtGenerics = append(annotations.AtGenerics, lexedGeneric)
	}

	return true, nil
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

func (l *Lexer) lexAtParamAnnotations(buffer *nvim.ScratchBuffer, annotations *AtAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(atParamAnnotationQuery)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, err
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
			return false, fmt.Errorf("Could not retrieve param name for param annotation")
		}

		annotations.AtParams[lexedParam.Name] = lexedParam
	}

	return true, nil
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

func (l *Lexer) lexAtReturnAnnotations(buffer *nvim.ScratchBuffer, annotations *AtAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(atReturnAnnotationQuery)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
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

		annotations.AtReturns = append(annotations.AtReturns, returnAnnotation)
	}

	return true, nil
}

var atFieldAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
		(documentation
			(field_annotation
				"@field"
				.
				([
					(qualifier "public")
					(qualifier "private") @field.private
					(qualifier "protected") @field.protected
					(qualifier "package") @field.package
				 ])?
				.
				(identifier) @field.name
				.
				"?"? @field.optional
				.
				(%s) @field.type
				.
				(comment)? @field.documentation
				.
			) @field
		)`, anyTypeAnnotationQuery),
}

func (l *Lexer) lexAtFieldAnnotations(buffer *nvim.ScratchBuffer, annotations *AtAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(atFieldAnnotationQuery)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		field := AtFieldAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "field.name":
				field.Name = capture.Node.Text
			case "field.optional":
				field.Optional = true
			case "field.private":
				field.Private = true
			case "field.protected":
				field.Protected = true
			case "field.package":
				field.Package = true
			case "field.type":
				field.Type = TypeAnnotation{capture.Node.Text}
			case "field.documentation":
				field.Documentation = []string{capture.Node.Text}
			}
		}

		if field.Name == "" {
			return false, fmt.Errorf("Error lexing field annotation: could not find captured field name")
		}

		annotations.AtFields[field.Name] = field
	}

	return true, nil
}

func (l *Lexer) lexAtAnnotations(dockblock []string) (AtAnnotations, error) {
	if cachedAnnotations, cached := lexedAnnotationsCache.Get(dockblock...); cached {
		// fmt.Printf("\nUsing cached lexedAnnotations: %+v\n", cachedAnnotations)
		return cachedAnnotations, nil
	}

	annotations := AtAnnotations{
		AtType:      nil,
		AtGenerics:  []AtGenericAnnotation{},
		AtParams:    map[string]AtParamAnnotation{},
		AtReturns:   []AtReturnAnnotation{},
		AtOverloads: []AtOverloadAnnotation{},
		AtFields:    map[string]AtFieldAnnotation{},
	}

	buffer, err := l.context.Nvim().NewBuffer()

	if err != nil {
		return annotations, err
	}

	defer buffer.Close()

	err = buffer.SetLines(dockblock)

	if err != nil {
		return annotations, err
	}

	// Generics are lexed ahead of other annotations, which could read them
	_, err = l.lexAtGenericAnnotations(buffer, &annotations)

	if err != nil {
		return annotations, fmt.Errorf("Error lexing generic annotations: %w", err)
	}

	_, err = l.lexAtParamAnnotations(buffer, &annotations)

	if err != nil {
		return annotations, fmt.Errorf("Error lexing param annotation: %w", err)
	}

	_, err = l.lexAtOverloadAnnotations(buffer, &annotations)

	if err != nil {
		return annotations, fmt.Errorf("Error lexing overload annotation: %w", err)
	}

	_, err = l.lexAtReturnAnnotations(buffer, &annotations)

	if err != nil {
		return annotations, fmt.Errorf("Error lexing return annotation: %w", err)
	}

	_, err = l.lexAtTypeAnnotations(buffer, &annotations)

	if err != nil {
		return annotations, fmt.Errorf("Error lexing type annotation: %w", err)
	}

	_, err = l.lexAtFieldAnnotations(buffer, &annotations)

	if err != nil {
		return annotations, err
	}

	// fmt.Printf("\nLexedAnnotations: %+v\n", annotations)

	lexedAnnotationsCache.Set(annotations, dockblock...)
	return annotations, nil
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

func (l *Lexer) lexOptionalTypeAnnotation(buffer *nvim.ScratchBuffer, source string) (symbol.Type, error) {
	match, err := buffer.TsQueryOne(optionalTypeAnnotationQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
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

func (l *Lexer) lexBuiltinTypeAnnotation(source string) symbol.Type {
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
		Language: "luadoc",
		Query: `
		(documentation
			(type_annotation
				(table_literal_type . "{" . "}" . ) @table
			)
		)`,
	},
	"described": {
		Language: "luadoc",
		Query: fmt.Sprintf(`
		(documentation
			(type_annotation
				(table_literal_type
					"{"
					field: ([
						(
							"["
							.
							(_) @table.index.key
							.
							"]"
							.
							"?"? @table.index.optional
							.
							":"
							.
							(%s) @table.index.value
						) @table.index
						(
							(identifier) @table.field.name
							.
							"?"? @table.field.optional
							.
							":"
							.
							(%s) @table.field.value
						) @table.field
					]
					","?
					)+
					"}"
				) @table
				(comment)? @table.documentation
			)
		)`, anyTypeAnnotationQuery, anyTypeAnnotationQuery),
	},
}

func (l *Lexer) lexLiteralTableTypeAnnotation(buffer *nvim.ScratchBuffer) (*symbol.Table, error) {
	for _, query := range literalTableAnnotationQueries {
		match, err := buffer.TsQueryOne(query)

		switch {
		case err != nil:
			return nil, err
		case match == nil:
			continue
		}

		table := symbol.NewTable()

		for _, capture := range *match {
			switch capture.Id {

			case "table.field.name":
				table.Fields = append(table.Fields, symbol.NewSymbol(capture.Node.Text, symbol.Meta{}, symbol.Documentation{}, symbol.NewUnknown()))
			case "table.field.optional":
				// TODO: recover
				// table.Fields[len(table.Fields)-1].Optional = true
			case "table.field.value":
				lexedType, err := l.lexTypeAnnotation(TypeAnnotation{capture.Node.Text})
				if err != nil {
					return nil, err
				}
				table.Fields[len(table.Fields)-1].Type = lexedType

			case "table.index.key":
				table.Indexes = append(table.Indexes, *symbol.NewTableIndex())
				lexedKey, err := l.lexTypeAnnotation(TypeAnnotation{capture.Node.Text})
				if err != nil {
					return nil, err
				}
				table.Indexes[len(table.Indexes)-1].Key = lexedKey
			case "table.index.optional":
				table.Indexes[len(table.Indexes)-1].Optional = true
			case "table.index.value":
				lexedValue, err := l.lexTypeAnnotation(TypeAnnotation{capture.Node.Text})
				if err != nil {
					return nil, err
				}
				table.Indexes[len(table.Indexes)-1].Value = lexedValue

				/* case "table.documentation":
				table.Documentation = []string{capture.Node.Text} */

			}
		}

		return table, nil
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

	unionTypes := []symbol.Type{}

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

func (l *Lexer) lexParenthesizedTypeAnnotation(buffer *nvim.ScratchBuffer, source string) (symbol.Type, error) {
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

func (l *Lexer) lexTypeAnnotation(typ TypeAnnotation) (symbol.Type, error) {
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
