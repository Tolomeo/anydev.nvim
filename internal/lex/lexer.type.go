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

	_, foundClassAnnotation := classAnnotations.AtClasses[classPatchedName]

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

		err := l.context.Push(fieldName, func(path string) error {
			fmt.Printf("\nLexing: '%s' class: '%s' field\n", className, fieldName)

			field := symbol.TableField{Name: fieldName, Value: symbol.NewUnknown()}
			fieldSource, err := l.crawler.SourceType(fieldName, className)

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

			field.Private = fieldAnnotations.AtPrivate
			field.Protected = fieldAnnotations.AtProtected
			field.Package = fieldAnnotations.AtPackage
			field.Deprecated = fieldAnnotations.AtDeprecated
			field.Protected = fieldAnnotations.AtProtected
			fieldValue, err := l.lex(fieldSource)

			switch {
			case err != nil:
				return err
			case fieldValue == nil:
				l.context.Logger.Warn(fmt.Sprintf("No types found for type '%s' field %s", className, fieldName))
				field.Value = symbol.NewUnknown()
			default:
				field.Value = fieldValue
			}

			class.Fields = append(class.Fields, field)
			return nil
		})

		if err != nil {
			return nil, err
		}

	}

	return class, nil
}

func (l *Lexer) lexFieldType(source symbol.Source) (symbol.Symbol, error) {
	annotations, err := l.lexAtAnnotations(source.GetOrigin().DocumentationLines())

	if err != nil {
		return nil, err
	}

	field, found := annotations.AtFields[source.Name()]

	if !found {
		return nil, nil
	}

	// TODO: assign qualifiers taken from annotations
	lexedField, err := l.lexTypeAnnotation(field.Type)

	if err != nil {
		return nil, err
	}

	return lexedField, nil
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

	alias, found := lexedAnnotations.AtAliases[patchedName]

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
