package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/cache"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

var lexedAnnotationsCache = cache.NewCache[*lexedAnnotations]()

type lexedType struct {
	Name string
}

type lexedTypeAnnotation struct {
	Type          lexedType
	Documentation []string
}

type lexedGenericAnnotation struct {
	Name  string
	Types []lexedType
}

type lexedReturnAnnotation struct {
	Name          *string
	Type          lexedType
	Documentation []string
}

type lexedParamAnnotation struct {
	Name          string
	Type          lexedType
	Optional      bool
	Documentation []string
}

type lexedOverloadAnnotation struct {
	Type          lexedType
	Documentation []string
}

type lexedAliasAnnotation struct {
	Name          string
	Type          lexedType
	Documentation []string
}

type lexedClassAnnotation struct {
	Name    string
	Parents []string
}

type lexedAnnotations struct {
	Type      *lexedTypeAnnotation
	Private   bool
	Protected bool
	Generics  []lexedGenericAnnotation
	Params    map[string]lexedParamAnnotation
	Returns   []lexedReturnAnnotation
	Overloads []lexedOverloadAnnotation
	Aliases   map[string]lexedAliasAnnotation
	Classes   map[string]lexedClassAnnotation
}

// TODO: support for multiple comma-separated types
var typeAnnotationQuery string = fmt.Sprintf(`
	(documentation
		(type_annotation
			(%s) @type.type
			(comment)? @type.documentation
		)
	) @type
`, anyTypeQuery)

func (l *Lexer) lexTypeAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	captures, err := buffer.TsQueryOne(treesitter.Query{Language: "luadoc", Query: typeAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	lexed := lexedTypeAnnotation{}

	for _, capture := range *captures {
		switch capture.Id {
		case "type.type":
			lexed.Type = lexedType{capture.Node.Text}
		case "type.documentation":
			lexed.Documentation = strings.Split(capture.Node.Text, "\n")
		}
	}

	annotations.Type = &lexed
	return true, nil
}

var overloadAnnotationQuery string = fmt.Sprintf(`
	(documentation 
		(overload_annotation 
			(%s) @type 
			(comment)? @documentation
		)
	)
`, typeQueries["function_type"])

func (l *Lexer) lexOverloadAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, err := buffer.SafeTsQueryAll(treesitter.Query{Language: "luadoc", Query: overloadAnnotationQuery})

	// fmt.Printf("\n Overload matches: %+v\n", matches)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, match := range *matches {
		if match.HasError {
			l.context.Logger.Warn(fmt.Sprintf("Skipping overload annotation in '%s' because it contains syntax errors", l.context.Current()))
			continue
		}

		overload := lexedOverloadAnnotation{}

		for _, capture := range match.Captures {
			switch capture.Id {
			case "documentation":
				overload.Documentation = []string{capture.Node.Text}
			case "type":
				overload.Type = lexedType{capture.Node.Text}
			}
		}

		annotations.Overloads = append(annotations.Overloads, overload)
	}

	return true, nil
}

var genericAnnotationQuery string = fmt.Sprintf(`
	(documentation 
		(generic_annotation
			(identifier) @generic.name
			parent_type: 
				(%s)? @generic.type
		)
	) @generic
`, anyTypeQuery)

func (l *Lexer) lexGenericAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(treesitter.Query{Language: "luadoc", Query: genericAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		lexedGeneric := lexedGenericAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "generic.name":
				lexedGeneric.Name = capture.Node.Text
			case "generic.type":
				lexedGeneric.Types = append(lexedGeneric.Types, lexedType{capture.Node.Text})
			}
		}

		if lexedGeneric.Name == "" {
			return false, fmt.Errorf("Could not retrieve generic name for generic annotation '%v'", lexedGeneric)
		}

		annotations.Generics = append(annotations.Generics, lexedGeneric)
	}

	return true, nil
}

var paramAnnotationQueries = map[string]string{
	"arg": fmt.Sprintf(`
		(documentation
			(param_annotation
				(identifier) @name
				"?"? @optional
				(%s) @type
				(comment)? @documentation
			)
		) @param
	`, anyTypeQuery),
	"vararg": fmt.Sprintf(`
		(documentation
			(param_annotation 
				"..." @name
				(%s) @type
				(comment)? @documentation
			)
		) @param
	`, anyTypeQuery),
}

func (l *Lexer) lexParamAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, hasMatches, err := slicesx.MapFindFunc(
		mapx.Values(paramAnnotationQueries),
		func(paramAnnotationQuery string) (*[]nvim.TsQueryMatch, bool, error) {
			paramMatches, err := buffer.TsQueryAll(treesitter.Query{Language: "luadoc", Query: paramAnnotationQuery})

			switch {
			case err != nil:
				return nil, false, err
			case paramMatches == nil:
				return nil, false, nil
			}

			return paramMatches, true, nil
		},
	)

	switch {
	case err != nil:
		return false, err
	case !hasMatches:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		lexedParam := lexedParamAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "name":
				lexedParam.Name = capture.Node.Text
			case "optional":
				lexedParam.Optional = true
			case "documentation":
				lexedParam.Documentation = []string{capture.Node.Text}
			case "type":
				lexedParam.Type = lexedType{capture.Node.Text}
			}
		}

		if lexedParam.Name == "" {
			return false, fmt.Errorf("Could not retrieve param name for param annotation")
		}

		annotations.Params[lexedParam.Name] = lexedParam
	}

	return true, nil
}

var returnAnnotationQuery string = fmt.Sprintf(`
	(documentation
		(return_annotation
			(%s) @return.type
			(comment)? @return.documentation
		)
	) @return
`, anyTypeQuery)

func (l *Lexer) lexReturnAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(treesitter.Query{Language: "luadoc", Query: returnAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		returnAnnotation := lexedReturnAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "return.name":
				returnAnnotation.Name = &capture.Node.Text
			case "return.documentation":
				returnAnnotation.Documentation = []string{capture.Node.Text}
			case "return.type":
				returnAnnotation.Type = lexedType{capture.Node.Text}
			}
		}

		annotations.Returns = append(annotations.Returns, returnAnnotation)
	}

	return true, nil
}

var privateAnnotationQuery string = `
	(documentation 
		(qualifier_annotation) @qualifier
		(#match? @qualifier "\\@private")
	)
`

func (l *Lexer) lexPrivateAnnotation(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	captures, err := buffer.TsQueryOne(treesitter.Query{Language: "luadoc", Query: privateAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	annotations.Private = true

	return true, nil
}

var protectedAnnotationQuery string = `
	(documentation 
		(qualifier_annotation) @qualifier
		(#match? @qualifier "\\@protected")
	)
`

func (l *Lexer) lexProtectedAnnotation(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	captures, err := buffer.TsQueryOne(treesitter.Query{Language: "luadoc", Query: protectedAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	annotations.Protected = true

	return true, nil
}

// Luadoc matches an empty type node even when the type is not present
// So those false positives are excluded with the not-eq predicate
var simpleAliasQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(alias_annotation
		(identifier) @alias.name
		(%s) @alias.type
		(comment)? @alias.documentation
		(#not-eq? @alias.type "")
	)`, anyTypeQuery),
}

func (l *Lexer) lexSimpleAliasAnnotation(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(simpleAliasQuery)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		alias := lexedAliasAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "alias.name":
				alias.Name = capture.Node.Text
			case "alias.type":
				alias.Type = lexedType{capture.Node.Text}
			case "alias.documentation":
				alias.Documentation = append(alias.Documentation, capture.Node.Text)
			}
		}

		if alias.Name == "" {
			return false, fmt.Errorf("Error lexing simple alias annotation: could not find captured alias name")
		}

		if alias.Type.Name == "" {
			return false, fmt.Errorf("Error lexing simple alias annotation: could not find captured alias type")
		}

		annotations.Aliases[alias.Name] = alias
	}

	return true, nil
}

// Luadoc matches an empty type node even when the type is not present
// That means that enum aliases have an empty type node defined
var enumAliasQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(alias_annotation
		(identifier) @alias.name
		(%s) @alias.type
		(comment)? @alias.documentation
		(#eq? @alias.type "")
	)`, anyTypeQuery),
}
var enumAliasMemberQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(continuation
		(%s) @alias.type
	)
`, anyTypeQuery)}

func (l *Lexer) lexEnumAliasAnnotation(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(enumAliasQuery)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		enumAlias := lexedAliasAnnotation{}
		enumAliasMembers := []string{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "alias.name":
				enumAlias.Name = capture.Node.Text
			}
		}

		if enumAlias.Name == "" {
			return false, fmt.Errorf("Error lexing enum alias annotation: could not find captured alias name")
		}

		nextLines, err := buffer.NextLineIterator(uint(matchCaptures.LineRange().Start + 1))

		if err != nil {
			return false, err
		}

		for line, err := range nextLines {
			if err != nil {
				return false, err
			}

			enumAliasMemberMatch, err := line.TsQueryOne(enumAliasMemberQuery)

			if err != nil {
				return false, nil
			}

			if enumAliasMemberMatch == nil {
				break
			}

			enumAliasMemberTypeCapture, found := slicesx.FindFunc(*enumAliasMemberMatch, func(capture treesitter.Capture) bool {
				return capture.Id == "alias.type"
			})

			if !found {
				return false, fmt.Errorf("Error lexing enum alias annotation: could not find captured alias member type")
			}

			enumAliasMembers = append(enumAliasMembers, enumAliasMemberTypeCapture.Node.Text)
		}

		if len(enumAliasMembers) < 1 {
			return false, fmt.Errorf("Could not retrieve enum members from enum alias")
		}

		enumAlias.Type = lexedType{strings.Join(enumAliasMembers, "|")}

		annotations.Aliases[enumAlias.Name] = enumAlias
	}

	return true, nil
}

func (l *Lexer) lexAliasAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	found, err := l.lexSimpleAliasAnnotation(buffer, annotations)

	switch {
	case err != nil:
		return false, err
	case found:
		return true, nil
	}

	found, err = l.lexEnumAliasAnnotation(buffer, annotations)

	switch {
	case err != nil:
		return false, err
	case found:
		return true, nil
	}

	return false, nil
}

var classAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
		(documentation
			(class_annotation
				(identifier) @class.name
				(":" 
					. (%s) @class.parent
					("," (%s) @class.parent)*
				)?
			) @class
		)`, anyTypeQuery, anyTypeQuery),
}

/* var classFieldAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
		(documentation
			(field_annotation
				(identifier) @field.name
				(%s) @field.type
				(comment)? @field.documentation
			) @field
		)`, anyTypeQuery),
} */

func (l *Lexer) lexClassAnnotation(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(classAnnotationQuery)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		class := lexedClassAnnotation{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "class.name":
				class.Name = capture.Node.Text
			case "class.parent":
				class.Parents = append(class.Parents, capture.Node.Text)
			}
		}

		if class.Name == "" {
			return false, fmt.Errorf("Error lexing class annotation: could not find captured class name")
		}

		// class := symbol.NewTable()

		annotations.Classes[class.Name] = class
	}

	/* fieldMatches, err := buffer.TsQueryAll(classFieldAnnotationQuery)

	switch {
	case err != nil:
		return false, err
	case fieldMatches == nil:
		l.context.Logger.Warn(fmt.Sprintf("Class '%s' with no fields", l.context.Current()))
		return true, nil
	}

	for _, fieldCaptures := range *fieldMatches {
		classField := symbol.NewTableField()

		for _, fieldCapture := range fieldCaptures {
			// TODO: field.documentation
			switch fieldCapture.Id {
			case "field.name":
				classField.Name = fieldCapture.Node.Text
			case "field.type":
				lexedFieldType, err := l.lexType(fieldCapture.Node.Text)

				if err != nil {
					return true, err
				}

				classField.Value = lexedFieldType
			}
		}

		class.Fields = append(class.Fields, *classField)
	}

	annotations.classes = class */
	return true, nil
}

func (l *Lexer) lexAnnotations(dockblock []string) (*lexedAnnotations, error) {
	if cachedAnnotations, cached := lexedAnnotationsCache.Get(dockblock...); cached {
		fmt.Printf("\nUsing cached lexedAnnotations: %+v\n", cachedAnnotations)
		return cachedAnnotations, nil
	}

	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	annotations := lexedAnnotations{
		Type:      nil,
		Private:   false,
		Protected: false,
		Generics:  []lexedGenericAnnotation{},
		Params:    map[string]lexedParamAnnotation{},
		Returns:   []lexedReturnAnnotation{},
		Overloads: []lexedOverloadAnnotation{},
		Aliases:   map[string]lexedAliasAnnotation{},
		Classes:   map[string]lexedClassAnnotation{},
	}

	err = buffer.SetLines(dockblock)

	if err != nil {
		return nil, err
	}

	// Generics are lexed ahead of other annotations, which could read them
	_, err = l.lexGenericAnnotations(buffer, &annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing generic annotations: %w", err)
	}

	_, err = l.lexPrivateAnnotation(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing private annotation: %w", err)
	}

	_, err = l.lexProtectedAnnotation(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing private annotation: %w", err)
	}

	_, err = l.lexParamAnnotations(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing param annotation: %w", err)
	}

	_, err = l.lexOverloadAnnotations(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing overload annotation: %w", err)
	}

	_, err = l.lexReturnAnnotations(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing return annotation: %w", err)
	}

	_, err = l.lexTypeAnnotations(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type annotation: %w", err)
	}

	_, err = l.lexAliasAnnotations(buffer, &annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing alias annotations: %w", err)
	}

	_, err = l.lexClassAnnotation(buffer, &annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing alias annotations: %w", err)
	}

	fmt.Printf("\nLexedAnnotations: %+v\n", annotations)

	lexedAnnotationsCache.Set(&annotations, dockblock...)
	return &annotations, nil
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
			optionalType, err := l.lexType(lexedType{capture.Node.Text})

			if err != nil {
				return nil, err
			}

			return symbol.NewOptional(&optionalType), nil
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
					(%s) @parameter.type
				) @parameter
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
			parameterType, err := l.lexType(lexedType{capture.Node.Text})

			if err != nil {
				return nil, err
			}

			args[len(args)-1].Type = parameterType
		case "return.type":
			returnType, err := l.lexType(lexedType{capture.Node.Text})

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

var typeTableQuery = treesitter.Query{
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
	)`, anyTypeQuery, anyTypeQuery),
}

func (l *Lexer) lexTableType(buffer *nvim.Buffer) (*symbol.Table, error) {
	matches, err := buffer.TsQueryAll(typeTableQuery)

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
				table.Fields = append(table.Fields, *symbol.NewTableField())
			case "key":
				table.Fields[len(table.Fields)-1].Name = capture.Node.Text
			case "value":
				valueType, err := l.lexType(lexedType{capture.Node.Text})

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

func (l *Lexer) lexReferenceType(name string) (*symbol.Reference, error) {
	err := l.LexType(name, l.context)

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
			itemsType, err := l.lexType(lexedType{matchCapture.Node.Text})

			if err != nil {
				return nil, err
			}

			return symbol.NewArray(itemsType), nil
		}
	}

	return nil, fmt.Errorf("Could not retrieve items type value from the array type '%s'", source)
}

var tableLiteralQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(documentation
		(type_annotation
			(table_literal_type
				field: ((%s) @table.key 
				(%s) @table.value) @table.field 
				(","
					field: ((%s) @table.key 
					(%s) @table.value)
				) @table.field
			) @table
		)
	)`, anyTypeQuery, anyTypeQuery, anyTypeQuery, anyTypeQuery),
}

func (l *Lexer) lexLiteralTableType(buffer *nvim.Buffer) (*symbol.Table, error) {
	match, err := buffer.TsQueryOne(tableLiteralQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	table := symbol.NewTable()

	for _, capture := range *match {
		switch capture.Id {
		case "table.field":
			table.Fields = append(table.Fields, *symbol.NewTableField())
		case "table.key":
			table.Fields[len(table.Fields)-1].Name = capture.Node.Text
		case "table.value":
			value, err := l.lexType(lexedType{capture.Node.Text})

			if err != nil {
				return nil, err
			}

			table.Fields[len(table.Fields)-1].Value = value
		}
	}

	return table, nil
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
			lexedType, err := l.lexType(lexedType{matchCapture.Node.Text})

			if err != nil {
				return nil, err
			}

			// Flattening nested unions
			switch t := lexedType.(type) {
			case symbol.Union:
				for _, lexedUnionType := range t.Types {
					unionTypes = append(unionTypes, lexedUnionType)
				}
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
			return l.lexType(lexedType{matchCapture.Node.Text})
		}
	}

	return nil, fmt.Errorf("Could not retrieve the type value of the grouped type '%s'", source)
}

var typeNumericLiteralQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(documentation
		(type_annotation
			(numeric_literal_type) @literal.number
		) @literal
	)`,
}

func (l *Lexer) lexNumericLiteralType(buffer *nvim.Buffer, source string) (*symbol.NumericLiteral, error) {
	match, err := buffer.TsQueryOne(typeNumericLiteralQuery)

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

func (l *Lexer) lexType(typ lexedType) (symbol.Symbol, error) {
	source := strings.TrimSpace(typ.Name)
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

	literalTableType, err := l.lexLiteralTableType(buffer)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case literalTableType != nil:
		return literalTableType, nil
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

	lexedNumericLiteral, err := l.lexNumericLiteralType(buffer, source)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type %s: %w", source, err)
	case lexedNumericLiteral != nil:
		return lexedNumericLiteral, nil
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
