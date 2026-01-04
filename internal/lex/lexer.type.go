package lex

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

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

func (l *Lexer) lexOptionalType(buffer *nvim.Buffer, source string) (*symbol.Optional, error) {
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

			return symbol.NewOptional(&lexedType), nil
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

func (l *Lexer) lexFunctionType(buffer *nvim.Buffer) (*symbol.Function, error) {
	captures, err := buffer.TsQueryOne(typeFunctionQuery)

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
			parameterType, err := l.lexType(capture.Node.Text)

			if err != nil {
				return nil, err
			}

			args[len(args)-1].Type = parameterType
		case "return.type":
			returnType, err := l.lexType(capture.Node.Text)

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

func (l *Lexer) lexTableType(buffer *nvim.Buffer) (*symbol.Table, error) {
	matches, err := buffer.TsQueryAll(treesitter.Query{Language: "luadoc", Query: typeTableQuery})

	switch {
	case err != nil:
		return nil, err
	case matches == nil:
		return nil, nil
	}

	table := symbol.NewTable()

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
					return nil, err
				}

				table.Fields[len(table.Fields)-1].Value = valueType
			}
		}
	}

	return table, nil
}

func (l *Lexer) lexBuiltinType(source string) symbol.Symbol {
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

func (l *Lexer) lexClassType(source *symbol.TypeSource) (*symbol.Table, error) {
	// Replacing all dots in the alias name with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	name := source.Identifier()
	origin := source.GetOrigin()

	patchedName := strings.ReplaceAll(name, ".", "_")
	definitionText := origin.DefinitionText()
	patchedDefinitionText := strings.Replace(definitionText, name, patchedName, 1)
	documentationText := origin.DocumentationText()
	patchedDocumentationLines := strings.Split(
		strings.Replace(documentationText, definitionText, patchedDefinitionText, 1),
		"\n",
	)

	lexedAnnotations, err := l.lexAnnotations(patchedDocumentationLines)

	switch {
	case err != nil:
		return nil, err
	case lexedAnnotations.class == nil:
		return nil, nil
	}

	class := lexedAnnotations.class
	// Replacing the name which was captured as patched with the original one
	class.Name = name
	class.Documentation = origin.DocumentationLines()

	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	annotation := fmt.Sprintf("---@type %s", name)
	ref := "local ref"
	completion := "ref."
	err = buffer.SetLines([]string{annotation, ref, completion})

	if err != nil {
		return nil, err
	}

	fmt.Println(name)
	classFields, err := buffer.GetTypeCompletion(2, uint(len(completion)))
	fmt.Println(classFields)

	if err != nil {
		return nil, err
	}

	for _, fieldName := range classFields {
		if found := slices.ContainsFunc(class.Fields, func(field symbol.TableField) bool {
			return field.Name == fieldName
		}); found {
			l.context.Logger.Info(fmt.Sprintf("Skipping '%s' field '%s': already lexed", name, fieldName))
			continue
		}

		err := l.context.Push(fieldName, func(path string) error {
			classField := symbol.TableField{Name: fieldName}
			source, err := l.crawler.SourceTypeMember(name, fieldName)
			origin := source.GetOrigin()

			switch {
			case err != nil:
				return err
			case source == nil:
				l.context.Logger.Warn(fmt.Sprintf("Using unknown for '%s' field '%s', with no origin", name, fieldName))
				classField.Value = symbol.NewUnknown()
				return nil
			}

			annotations, err := l.lexAnnotations(origin.DocumentationLines())
			classField.Private = annotations.private
			classField.Protected = annotations.protected
			classFieldValue, err := l.lexValue(source)

			if err != nil {
				return err
			}

			classField.Value = classFieldValue

			class.Fields = append(class.Fields, classField)
			return nil
		})

		if err != nil {
			return nil, err
		}

	}

	return class, nil
}

func (l *Lexer) lexAliasType(source *symbol.TypeSource) (symbol.Symbol, error) {
	// Replacing all dots in the alias name with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	name := source.Path
	origin := source.GetOrigin()

	definitionText := origin.DefinitionText()
	patchedDefinitionText := strings.Replace(definitionText, name, strings.ReplaceAll(name, ".", "_"), 1)
	documentationText := origin.DocumentationText()
	patchedDocumentationLines := strings.Split(
		strings.Replace(documentationText, definitionText, patchedDefinitionText, 1),
		"\n",
	)

	lexedAnnotations, err := l.lexAnnotations(patchedDocumentationLines)

	switch {
	case err != nil:
		return nil, err
	case lexedAnnotations.alias == nil:
		return nil, nil
	}

	// TODO: replace the name of the returned type with the original name
	// TODO: attach documentation
	alias := lexedAnnotations.alias

	return alias, nil
}

func (l *Lexer) lexReferenceType(name string) (*symbol.Reference, error) {
	if _, alreadyLexed := l.context.Result().Types[name]; alreadyLexed {
		l.context.Logger.Info(fmt.Sprintf("Skipping '%s': lexed type already found", name))
		return symbol.NewReference(name), nil
	}

	l.context.Result().Types[name] = symbol.NewUnknown()

	err := l.context.Fork(name, func(name string) error {
		typeSource, err := l.crawler.SourceType(name)

		if err != nil {
			return err
		}

		// fmt.Printf("\nReference '%s' source:\n%+v\n\n", name, typeSource.Origin)

		aliasType, err := l.lexAliasType(typeSource)

		// fmt.Printf("\nLexed '%s' alias: %+v\n\n", name, aliasType)

		// fmt.Printf("\nType: %+v\n\n", typeSource.Origin.Definition)

		switch {
		case err != nil:
			return err
		case aliasType != nil:
			l.context.Result().Types[name] = aliasType
			return nil
		}

		classType, err := l.lexClassType(typeSource)

		switch {
		case err != nil:
			return err
		case classType != nil:
			l.context.Result().Types[name] = classType
			return nil
		}

		l.context.Logger.Warn(fmt.Sprintf("Uknown type '%s' received", name))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return symbol.NewReference(name), nil
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

func (l *Lexer) lexArrayType(buffer *nvim.Buffer, source string) (*symbol.Array, error) {
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

			return symbol.NewArray(itemsType), nil
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

func (l *Lexer) lexUnionType(buffer *nvim.Buffer, source string) (*symbol.Union, error) {
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

	return symbol.NewUnion(unionTypes), nil
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

func (l *Lexer) lexGroupType(buffer *nvim.Buffer, source string) (symbol.Symbol, error) {
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

func (l *Lexer) lexStringLiteralType(buffer *nvim.Buffer, source string) (*symbol.StringLiteral, error) {
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
			return symbol.NewStringLiteral(matchCapture.Node.Text), nil
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

	functionType, err := l.lexFunctionType(buffer)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case functionType != nil:
		return functionType, nil
	}

	lexedArray, err := l.lexArrayType(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedArray != nil:
		return lexedArray, nil
	}

	tableType, err := l.lexTableType(buffer)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case tableType != nil:
		return tableType, nil
	}

	lexedOptional, err := l.lexOptionalType(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedOptional != nil:
		return lexedOptional, nil
	}

	lexedUnion, err := l.lexUnionType(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedUnion != nil:
		return lexedUnion, nil
	}

	lexedGroupedType, err := l.lexGroupType(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedGroupedType != nil:
		return lexedGroupedType, nil
	}

	lexedStringLiteral, err := l.lexStringLiteralType(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedStringLiteral != nil:
		return lexedStringLiteral, nil
	}

	lexedReference, err := l.lexReferenceType(source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedReference != nil:
		return lexedReference, nil
	}

	l.context.Logger.Warn(fmt.Sprintf("Unknown type '%s' received", source))
	return symbol.NewUnknown(), nil
}
