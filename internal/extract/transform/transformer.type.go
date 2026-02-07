package transform

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type AtTypeAnnotation struct {
	Types         []string
	Documentation []string
}

func (tr *Transformer) getAtTypeAnnotations(buffer *nvim.ScratchBuffer) (*AtTypeAnnotation, error) {
	captures, err := buffer.TsQueryOne(annotation.AtTypeQuery)

	if err != nil {
		return nil, err
	}

	if captures == nil {
		return nil, nil
	}

	atTypeAnnotation := annotation.NewAtType(*captures)
	atType := AtTypeAnnotation{}

	for _, typ := range atTypeAnnotation.Types() {
		atType.Types = append(atType.Types, typ)
	}

	return &atType, nil
}

var optionalTypeAnnotationQuery = annotation.OptionalQuery.Extend(func(query string) string {
	return fmt.Sprintf(`
	(documentation
		(type_annotation 
			(%s)
		)
	)`, query)
})

func (tr *Transformer) getOptionalType(buffer *nvim.ScratchBuffer, _ string) (annotation.Type, error) {
	match, err := buffer.TsQueryOne(optionalTypeAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if match == nil {
		return nil, nil
	}

	optionalAnnotation := annotation.NewOptional(*match)
	optionalType, err := tr.getType(optionalAnnotation.Type())

	if err != nil {
		return nil, err
	}

	return symbol.NewOptional(optionalType), nil
}

var functionTypeAnnotationQuery = annotation.FunctionQuery.Extend(func(query string) string {
	return fmt.Sprintf(`
	(documentation
		(type_annotation
			(%s)
		)
	)
`, query)
})

func (tr *Transformer) getFunctionType(buffer *nvim.ScratchBuffer) (*symbol.Function, error) {
	captures, err := buffer.TsQueryOne(functionTypeAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if captures == nil {
		return nil, nil
	}

	functionAnnotation := annotation.NewFunction(*captures)

	function := symbol.NewFunction()
	args := []symbol.FunctionArgument{}
	returns := []symbol.FunctionReturn{}

	for _, annotationArg := range functionAnnotation.Arguments() {
		argType, err := tr.getType(annotationArg.Type())

		if err != nil {
			return nil, err
		}

		args = append(args, symbol.FunctionArgument{
			Name: annotationArg.Name(),
			Type: argType,
		})
	}

	for _, annotationReturn := range functionAnnotation.Returns() {
		returnType, err := tr.getType(annotationReturn.Type())

		if err != nil {
			return nil, err
		}

		returns = append(returns, symbol.FunctionReturn{
			Type: returnType,
		})
	}

	function.Arguments = append(function.Arguments, args...)
	function.Returns = append(function.Returns, returns...)

	return function, nil
}

var tableTypeAnnotationQueries = []treesitter.Query{
	annotation.TableQuery.Extend(func(query string) string {
		return fmt.Sprintf(`
		(documentation
			(type_annotation
				(%s)
				(comment)? @documentation
			)
		)`, query)
	}),
	annotation.TableArrayQuery.Extend(func(query string) string {
		return fmt.Sprintf(`
		(documentation
			(type_annotation
				(%s)
				(comment)? @documentation
			)
		)`, query)
	}),
}

func (tr *Transformer) getTableType(buffer *nvim.ScratchBuffer) (*symbol.Table, error) {
	for _, tableTypeAnnotationQuery := range tableTypeAnnotationQueries {
		match, err := buffer.TsQueryOne(tableTypeAnnotationQuery)

		if err != nil {
			return nil, err
		}

		if match == nil {
			continue
		}

		tableAnnotation := annotation.NewTable(*match)
		table := symbol.NewTable()

		keyType, err := tr.getType(tableAnnotation.Key())

		if err != nil {
			return nil, err
		}

		valueType, err := tr.getType(tableAnnotation.Value())

		if err != nil {
			return nil, err
		}

		table.Indexes = append(table.Indexes, *symbol.NewTableIndex())
		table.Indexes[len(table.Indexes)-1].Key = keyType
		table.Indexes[len(table.Indexes)-1].Value = valueType

		return table, nil
	}

	return nil, nil
}

func (tr *Transformer) getBuiltinType(source string) annotation.Type {
	switch source {
	case annotation.BuiltinVoid:
		return symbol.NewVoid()
	case annotation.BuiltinNil:
		return symbol.NewNil()
	case annotation.BuiltinAny:
		return symbol.NewAny()
	case annotation.BuiltinBoolean:
		return symbol.NewBoolean()
	case annotation.BuiltinString:
		return symbol.NewString()
	case annotation.BuiltinNumber:
		return symbol.NewNumber()
	case annotation.BuiltinInteger, annotation.BuiltinInt:
		return symbol.NewInteger()
	case annotation.BuiltinFunction:
		return symbol.NewBuiltinFunction()
	case annotation.BuiltinTable:
		return symbol.NewBuiltinTable()
	case annotation.BuiltinThread:
		return symbol.NewThread()
	case annotation.BuiltinUserdata:
		return symbol.NewUserdata()
	case annotation.BuiltinLightUserdata:
		return symbol.NewLightUserdata()
	}

	return nil
}

func (tr *Transformer) getReferenceType(name string) (*symbol.Reference, error) {
	err := tr.context.Extract("type", name)

	if err != nil {
		return nil, err
	}

	return symbol.NewReference(name), nil
}

var arrayTypeAnnotationQuery = annotation.ArrayQuery.Extend(func(query string) string {
	return fmt.Sprintf(`
	(documentation
		(type_annotation
			(%s)
		)
	)
`, query)
})

func (tr *Transformer) getArrayType(buffer *nvim.ScratchBuffer, _ string) (*symbol.Array, error) {
	match, err := buffer.TsQueryOne(arrayTypeAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if match == nil {
		return nil, nil
	}

	arrayAnnotation := annotation.NewArray(*match)
	itemsType, err := tr.getType(arrayAnnotation.ItemsType())

	if err != nil {
		return nil, err
	}

	return symbol.NewArray(itemsType), nil
}

var literalTableAnnotationQueries = map[string]treesitter.Query{
	"empty": annotation.LiteralTableEmptyQuery.Extend(func(query string) string {
		return fmt.Sprintf(`
		(documentation
			(type_annotation
				(%s)
				(comment)? @table.documentation
			)
		)`, query)
	}),
	"described": annotation.LiteralTableTypeQuery.Extend(func(query string) string {
		return fmt.Sprintf(`
		(documentation
			(type_annotation
				(%s)
				(comment)? @table.documentation
			)
		)`, query)
	}),
}

func (tr *Transformer) getLiteralTableType(buffer *nvim.ScratchBuffer) (*symbol.Table, error) {
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
			tableFieldType, err := tr.getType(fieldType.Value())

			if err != nil {
				return nil, err
			}

			tableField.Type = tableFieldType
			tableSymbol.Fields = append(tableSymbol.Fields, *tableField)
		}

		for _, indexType := range literalTableType.Indexes() {
			tableIndex := symbol.NewTableIndex()
			tableIndexKey, err := tr.getType(indexType.Key())

			if err != nil {
				return nil, err
			}

			tableIndexValue, err := tr.getType(indexType.Value())

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

var unionTypeAnnotationQuery = annotation.UnionQuery.Extend(func(query string) string {
	return fmt.Sprintf(`
	(documentation
		(type_annotation
			(%s)
		)
	)`, query)
})

func (tr *Transformer) getUnionType(buffer *nvim.ScratchBuffer, source string) (*symbol.Union, error) {
	match, err := buffer.TsQueryOne(unionTypeAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if match == nil {
		return nil, nil
	}

	unionAnnotation := annotation.NewUnion(*match)

	unionTypes := []annotation.Type{}

	for _, typ := range unionAnnotation.Types() {
		unionType, err := tr.getType(typ)

		if err != nil {
			return nil, err
		}

		// Flattening nested unions
		switch t := unionType.(type) {
		case *symbol.Union:
			for _, symbolUnionType := range t.Types {
				unionTypes = append(unionTypes, symbolUnionType)
			}
		default:
			unionTypes = append(unionTypes, unionType)
		}
	}

	if len(unionTypes) < 2 {
		return nil, fmt.Errorf("Could not retrieve all types in the union type '%s'", source)
	}

	return symbol.NewUnion(unionTypes), nil
}

var parenthesizedTypeAnnotationQuery = annotation.ParenthesizedQuery.Extend(func(query string) string {
	return fmt.Sprintf(`
	(documentation
		(type_annotation
			(%s)
		)
	)
`, query)
})

func (tr *Transformer) getParenthesizedType(buffer *nvim.ScratchBuffer, _ string) (annotation.Type, error) {
	match, err := buffer.TsQueryOne(parenthesizedTypeAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if match == nil {
		return nil, nil
	}

	parenthesizedAnnotation := annotation.NewParenthesized(*match)
	return tr.getType(parenthesizedAnnotation.Type())
}

var literalNumberTypeAnnotationQuery = annotation.LiteralNumberQuery.Extend(func(query string) string {
	return fmt.Sprintf(`
		(documentation
			(type_annotation
				(%s)
			)
		)`, query)
})

func (tr *Transformer) getLiteralNumberType(buffer *nvim.ScratchBuffer, _ string) (*symbol.NumericLiteral, error) {
	match, err := buffer.TsQueryOne(literalNumberTypeAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if match == nil {
		return nil, nil
	}

	literalNumberAnnotation := annotation.NewLiteralNumber(*match)
	return symbol.NewNumericLiteral(literalNumberAnnotation.Value()), nil
}

var literalBooleanTypeAnnotationQuery = annotation.LiteralBooleanQuery.Extend(func(query string) string {
	return fmt.Sprintf(`
		(documentation
			(type_annotation
				(%s)
			)
		)`, query)
})

func (tr *Transformer) getLiteralBooleanType(buffer *nvim.ScratchBuffer) (*symbol.BooleanLiteral, error) {
	match, err := buffer.TsQueryOne(literalBooleanTypeAnnotationQuery)

	if err != nil {
		return nil, err
	}

	if match == nil {
		return nil, nil
	}

	literalBooleanAnnotation := annotation.NewLiteralBoolean(*match)
	return symbol.NewBooleanLiteral(literalBooleanAnnotation.Value()), nil
}

var literalStringTypeQuery = annotation.LiteralStringQuery.Extend(func(query string) string {
	return fmt.Sprintf(`
	(documentation
		(type_annotation
			(%s)
		)
	)`, query)
})

func (tr *Transformer) getLiteralStringType(buffer *nvim.ScratchBuffer, _ string) (*symbol.StringLiteral, error) {
	match, err := buffer.TsQueryOne(literalStringTypeQuery)

	if err != nil {
		return nil, err
	}

	if match == nil {
		return nil, nil
	}

	literalStringAnnotation := annotation.NewLiteralString(*match)
	return symbol.NewStringLiteral(literalStringAnnotation.Value()), nil
}

func (tr *Transformer) getType(typ string) (annotation.Type, error) {
	tr.context.Logger().Verbosef("Transforming text type <%s>", typ)

	source := strings.TrimSpace(typ)
	builtinType := tr.getBuiltinType(strings.TrimSpace(source))

	if builtinType != nil {
		return builtinType, nil
	}

	buffer, err := tr.context.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	typeAnnotation := fmt.Sprintf("---@type %s", source)
	err = buffer.SetLines([]string{typeAnnotation})

	if err != nil {
		return nil, err
	}

	functionType, err := tr.getFunctionType(buffer)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as function: %w", source, err)
	}

	if functionType != nil {
		return functionType, nil
	}

	arrayType, err := tr.getArrayType(buffer, source)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as array: %w", source, err)
	}

	if arrayType != nil {
		return arrayType, nil
	}

	tableType, err := tr.getTableType(buffer)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as table: %w", source, err)
	}

	if tableType != nil {
		return tableType, nil
	}

	literalTableType, err := tr.getLiteralTableType(buffer)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as literal table: %w", source, err)
	}

	if literalTableType != nil {
		return literalTableType, nil
	}

	optionalType, err := tr.getOptionalType(buffer, source)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as optional: %w", source, err)
	}

	if optionalType != nil {
		return optionalType, nil
	}

	unionType, err := tr.getUnionType(buffer, source)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as union: %w", source, err)
	}

	if unionType != nil {
		return unionType, nil
	}

	parenthesizedType, err := tr.getParenthesizedType(buffer, source)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as parenthesized: %w", source, err)
	}

	if parenthesizedType != nil {
		return parenthesizedType, nil
	}

	numericLiteralType, err := tr.getLiteralNumberType(buffer, source)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as numeric literal: %w", source, err)
	}

	if numericLiteralType != nil {
		return numericLiteralType, nil
	}

	booleanLiteralType, err := tr.getLiteralBooleanType(buffer)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as boolean literal: %w", source, err)
	}

	if booleanLiteralType != nil {
		return booleanLiteralType, nil
	}

	stringLiteralType, err := tr.getLiteralStringType(buffer, source)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as string literal: %w", source, err)
	}

	if stringLiteralType != nil {
		return stringLiteralType, nil
	}

	referenceType, err := tr.getReferenceType(source)

	if err != nil {
		return nil, fmt.Errorf("Error lexing type <%s> as reference: %w", source, err)
	}

	if referenceType != nil {
		return referenceType, nil
	}

	tr.context.Logger().Warn(fmt.Sprintf("Unknown type '%s' received", source))
	return symbol.NewUnknown(), nil
}
