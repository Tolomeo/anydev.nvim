package transform

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) transformTableOrigin(tableOrigin *origin.TableOrigin) (symbol.Type, error) {
	annotations, err := tr.getTableAnnotations(tableOrigin.Annotations())

	if err != nil {
		return nil, err
	}

	name := tableOrigin.Name()

	if annotations.AtClass != nil {
		return tr.getClassTableSymbol(name, *annotations.AtClass)
	}

	if annotations.AtEnum != nil {
		return tr.getEnumeratorTableSymbol(name, *annotations.AtEnum)
	}

	return tr.getTableSymbol(name)
}

type tableAnnotations struct {
	AtClass *annotation.AtClass
	AtEnum  *annotation.AtEnum
}

func (tr *Transformer) getTableAnnotations(docblock []string) (tableAnnotations, error) {
	tAnnotations := tableAnnotations{}
	buffer, err := tr.context.Nvim().NewBuffer()

	if err != nil {
		return tAnnotations, err
	}

	defer buffer.Close()

	err = buffer.SetLines(docblock)

	if err != nil {
		return tAnnotations, err
	}

	atClassMatch, err := buffer.TsQueryOne(annotation.AtClassQuery)

	if err != nil {
		return tAnnotations, err
	}

	if atClassMatch != nil {
		tAnnotations.AtClass = annotation.NewAtClass(*atClassMatch)
	}

	atEnumMatch, err := buffer.TsQueryOne(annotation.AtEnumQuery)

	if err != nil {
		return tAnnotations, err
	}

	if atEnumMatch != nil {
		tAnnotations.AtEnum = annotation.NewAtEnum(*atEnumMatch)
	}

	return tAnnotations, nil
}

func (tr *Transformer) getTableSymbol(name string) (*symbol.Table, error) {
	table := symbol.NewTable()
	table.Name = name
	tableChildren, err := tr.context.Nvim().GetValueCompletion(tr.context.Target().Identifier())

	if err != nil {
		return nil, err
	}

	for _, child := range tableChildren {
		err := tr.context.ExtractChild(table, child)

		if err != nil {
			return nil, err
		}
	}

	return table, nil
}

func (tr *Transformer) getClassTableSymbol(name string, atClassAnnotation annotation.AtClass) (*symbol.TypeReference, error) {
	table := symbol.NewTable()
	table.Name = name

	return tr.getReferenceType(atClassAnnotation.Name().Text)
}

func (tr *Transformer) getEnumeratorTableSymbol(name string, atEnum annotation.AtEnum) (*symbol.Table, error) {
	err := tr.context.Extract("type", atEnum.Name().Text)

	if err != nil {
		return nil, err
	}

	table := symbol.NewTable()
	table.Name = name
	tableFields, err := tr.context.Nvim().GetValueCompletion(tr.context.Target().Identifier())

	if err != nil {
		return nil, err
	}

	for _, fieldName := range tableFields {
		field := symbol.NewTableField()
		field.Name = fieldName
		field.Metadata = *symbol.NewMetadata()
		field.Documentation = []string{}
		field.Type = symbol.NewTypeReference(atEnum.Name().Text)

		table.Fields = append(table.Fields, *field)
	}

	return table, nil
}
