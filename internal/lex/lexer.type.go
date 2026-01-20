package lex

import (
	// "fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

func (l *Lexer) lexClassType() (*symbol.Table, error) {
	// Replacing all dots in the alias className with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	className := l.target.Identifier()
	classOrigin := l.target.Origin()

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
	class.Documentation = l.target.Origin().DocumentationLines()

	// fmt.Println(name)
	classFields, err := l.target.Nvim().GetTypeCompletion(className)
	// fmt.Println(classFields)

	if err != nil {
		return nil, err
	}

	for _, fieldName := range classFields {
		fieldValue, err := l.target.ExtractChild(fieldName)

		if err != nil {
			return nil, err
		}

		/* annotations, err := l.lexAtAnnotations(tableFieldSource.GetOrigin().DefinitionLines())
		tableField.Private = annotations.AtPrivate
		tableField.Protected = annotations.AtProtected
		tableField.Package = annotations.AtPackage
		tableField.Deprecated = annotations.AtDeprecated
		tableField.Protected = annotations.AtProtected */
		class.Fields = append(class.Fields, symbol.TableField{
			Name:  fieldName,
			Value: fieldValue,
		})
	}

	return class, nil
}

func (l *Lexer) lexFieldType() (symbol.Symbol, error) {
	annotations, err := l.lexAtAnnotations(l.target.Origin().DocumentationLines())

	if err != nil {
		return nil, err
	}

	field, found := annotations.AtFields[l.target.Name()]

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

func (l *Lexer) lexAliasType() (symbol.Symbol, error) {
	// fmt.Printf("\nsource: <%+v>\n", source.GetOrigin())
	// Replacing all dots in the alias name with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	name := l.target.Identifier()
	origin := l.target.Origin()

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
