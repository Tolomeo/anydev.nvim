package lex

import (
	// "fmt"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
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
	aliasType, hasType := origin.Captures().Find("alias.type")

	if !hasType {
		return nil, fmt.Errorf("No type capture found for alias type '%s'", l.context.Target().Identifier())
	}

	// TODO: attach documentation
	lexedAlias, err := l.lexTypeAnnotation(TypeAnnotation{aliasType.Node.Text})

	if err != nil {
		return nil, err
	}

	return lexedAlias, nil
}

func (l *Lexer) lexAliasEnumeratorType(origin *symbol.AliasEnumeratorOrigin) (symbol.Type, error) {
	l.context.Logger().Debugf("Alias enumerator: %+v", origin)

	typeCaptures, hasTypes := origin.Captures().FindAll("alias.type")

	if !hasTypes {
		return nil, fmt.Errorf("No type members found for alias enumerator '%s'", l.context.Target().Identifier())
	}

	enumeratorTypes, _ := slicesx.MapFunc(typeCaptures, func(capture treesitter.Capture) (string, error) {
		return capture.Node.Text, nil
	})

	enumeratorType, err := l.lexTypeAnnotation(TypeAnnotation{strings.Join(enumeratorTypes, "|")})

	if err != nil {
		return nil, err
	}

	l.context.Logger().Debugf("%+v", origin)

	return enumeratorType, nil

}
