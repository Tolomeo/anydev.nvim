package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

func (l *Lexer) lexClassType(source symbol.Source) (*symbol.Table, error) {
	// Replacing all dots in the alias className with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	className := source.Identifier()
	classOrigin := source.GetOrigin()

	classPatchedName := strings.ReplaceAll(className, ".", "_")
	classDefinitionText := classOrigin.DefinitionText()
	classPatchedDefinitionText := strings.Replace(classDefinitionText, className, classPatchedName, 1)

	classDocumentationText := classOrigin.DocumentationText()
	classPatchedDocumentationText := strings.Replace(classDocumentationText, classDefinitionText, classPatchedDefinitionText, 1)
	classPatchedDocumentationLines := strings.Split(classPatchedDocumentationText, "\n")

	classAnnotations, err := l.lexAtAnnotations(classPatchedDocumentationLines)

	if err != nil {
		return nil, err
	}

	_, foundClassAnnotation := classAnnotations.Classes[classPatchedName]

	if !foundClassAnnotation {
		return nil, nil
	}

	class := symbol.NewTable()
	class.Name = className
	// TODO: the documentation is gathered by the annotations lexer
	class.Documentation = source.GetOrigin().DocumentationLines()

	// fmt.Println(name)
	classFields, err := l.context.Nvim.GetTypeCompletion(className)
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

		// fmt.Println(fieldName)

		/* if fieldName != "fs_access" {
			continue
		} */

		err := l.context.Push(fieldName, func(path string) error {
			fmt.Printf("\nLexing: %s\n", l.context.Current())

			field := symbol.TableField{Name: fieldName, Value: symbol.NewUnknown()}
			fieldSource, err := l.crawler.SourceTypeMember(className, fieldName)

			switch {
			case err != nil:
				return err
			case fieldSource == nil:
				l.context.Logger.Warn(fmt.Sprintf("Using unknown for '%s' field '%s', with no origin", className, fieldName))
				class.Fields = append(class.Fields, field)
				return nil
			}

			fieldAnnotations, err := l.lexAtAnnotations(fieldSource.GetOrigin().DocumentationLines())

			if err != nil {
				return err
			}

			field.Private = fieldAnnotations.Private
			field.Protected = fieldAnnotations.Protected

			fieldValue, err := l.lex(fieldSource)

			if err != nil {
				return err
			}

			field.Value = fieldValue
			class.Fields = append(class.Fields, field)
			return nil
		})

		if err != nil {
			return nil, err
		}

	}

	return class, nil
}

func (l *Lexer) lexAliasType(source symbol.Source) (symbol.Symbol, error) {
	// fmt.Printf("\nsource: <%+v>\n", source.GetOrigin())
	// Replacing all dots in the alias name with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	name := source.Identifier()
	origin := source.GetOrigin()

	patchedName := strings.ReplaceAll(name, ".", "_")
	definitionText := origin.DefinitionText()
	patchedDefinitionText := strings.Replace(definitionText, name, patchedName, 1)

	documentationText := origin.DocumentationText()
	patchedDocumentationText :=
		strings.Replace(documentationText, definitionText, patchedDefinitionText, 1)
	patchedDocumentationLines := strings.Split(patchedDocumentationText, "\n")

	lexedAnnotations, err := l.lexAtAnnotations(patchedDocumentationLines)

	if err != nil {
		return nil, err
	}

	alias, found := lexedAnnotations.Aliases[patchedName]

	if !found {
		return nil, nil
	}

	// TODO: attach documentation
	lexedAlias, err := l.lexTypeAnnotation(alias.Type)

	if err != nil {
		return nil, err
	}

	return lexedAlias, nil
}
