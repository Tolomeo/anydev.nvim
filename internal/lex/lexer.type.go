package lex

import (
	// "fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

func (l *Lexer) lexClassType(origin *symbol.ClassOrigin) (*symbol.Table, error) {
	// Replacing all dots in the alias className with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	className := l.context.Target().Identifier()
	// classOrigin := l.context.Target().Origin()

	classPatchedName := strings.ReplaceAll(className, ".", "_")
	classDefinitionText := strings.Join(origin.Definition(), "\n")
	classPatchedDefinitionText := strings.Replace(classDefinitionText, className, classPatchedName, 1)

	classDocumentationText := strings.Join(origin.Documentation(), "\n")
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

	// fmt.Println(name)
	classFields, err := l.context.Nvim().GetTypeCompletion(className)
	// fmt.Println(classFields)

	if err != nil {
		return nil, err
	}

	for _, fieldName := range classFields {
		err := l.context.ExtractChild(class, fieldName)

		if err != nil {
			return nil, err
		}
	}

	return class, nil
}

func (l *Lexer) lexFieldType(origin *symbol.FieldOrigin) (symbol.Type, error) {
	annotations, err := l.lexAtAnnotations(origin.Documentation())

	if err != nil {
		return nil, err
	}

	field, found := annotations.AtFields[l.context.Target().Name()]

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

func (l *Lexer) lexAliasType(origin *symbol.AliasOrigin) (symbol.Type, error) {
	// fmt.Printf("\nsource: <%+v>\n", source.GetOrigin())
	// Replacing all dots in the alias name with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	name := l.context.Target().Identifier()
	// origin := l.context.Target().Origin()

	patchedName := strings.ReplaceAll(name, ".", "_")
	definitionText := strings.Join(origin.Definition(), "\n")
	patchedDefinitionText := strings.Replace(definitionText, name, patchedName, 1)

	documentationText := strings.Join(origin.Documentation(), "\n")
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
