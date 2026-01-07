package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

func (l *Lexer) lexClassType(source *symbol.TypeSource) (*symbol.Table, error) {
	// Replacing all dots in the alias name with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	name := source.Identifier()
	origin := source.GetOrigin()

	patchedName := strings.ReplaceAll(name, ".", "_")
	definitionText := origin.DefinitionText()
	patchedDefinitionText := strings.Replace(definitionText, name, patchedName, 1)

	documentationText := origin.DocumentationText()
	patchedDocumentationText := strings.Replace(documentationText, definitionText, patchedDefinitionText, 1)
	patchedDocumentationLines := strings.Split(patchedDocumentationText, "\n")

	lexedAnnotations, err := l.lexAnnotations(patchedDocumentationLines)

	if err != nil {
		return nil, err
	}

	_, foundClassAnnotation := lexedAnnotations.Classes[patchedName]

	if !foundClassAnnotation {
		return nil, nil
	}

	class := symbol.NewTable()
	class.Name = name
	// TODO: the documentation is gathered by the annotations lexer
	class.Documentation = origin.DocumentationLines()

	// fmt.Println(name)
	classFields, err := l.context.Nvim.GetTypeCompletion(name)
	// fmt.Println(classFields)

	if err != nil {
		return nil, err
	}

	for _, fieldName := range classFields {
		/* if found := slices.ContainsFunc(class.Fields, func(field symbol.TableField) bool {
			return field.Name == fieldName
		}); found {
			l.context.Logger.Info(fmt.Sprintf("Skipping '%s' field '%s': already lexed", name, fieldName))
			continue
		} */

		if fieldName != "run" {
			continue
		}

		err := l.context.Push(fieldName, func(path string) error {
			fmt.Printf("\nLexing: %s\n", l.context.Current())

			classField := symbol.TableField{Name: fieldName}
			source, err := l.crawler.SourceTypeMember(name, fieldName)

			switch {
			case err != nil:
				return err
			case source == nil:
				l.context.Logger.Warn(fmt.Sprintf("Using unknown for '%s' field '%s', with no origin", name, fieldName))
				classField.Value = symbol.NewUnknown()
				return nil
			}

			annotations, err := l.lexAnnotations(origin.DocumentationLines())

			if err != nil {
				return err
			}

			classField.Private = annotations.Private
			classField.Protected = annotations.Protected

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

	patchedName := strings.ReplaceAll(name, ".", "_")
	definitionText := origin.DefinitionText()
	patchedDefinitionText := strings.Replace(definitionText, name, patchedName, 1)

	documentationText := origin.DocumentationText()
	patchedDocumentationText :=
		strings.Replace(documentationText, definitionText, patchedDefinitionText, 1)
	patchedDocumentationLines := strings.Split(patchedDocumentationText, "\n")

	lexedAnnotations, err := l.lexAnnotations(patchedDocumentationLines)

	if err != nil {
		return nil, err
	}

	alias, found := lexedAnnotations.Aliases[patchedName]

	if !found {
		return nil, nil
	}

	// TODO: attach documentation
	lexedAlias, err := l.lexType(alias.Type)

	if err != nil {
		return nil, err
	}

	return lexedAlias, nil
}
