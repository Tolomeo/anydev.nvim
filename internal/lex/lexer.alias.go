package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

// Luadoc matches an empty type node even when the type is not present
// So those false positives are excluded with the not-eq predicate
var simpleAliasQuery = fmt.Sprintf(`
	(alias_annotation
		(identifier) @alias.name
		%s @alias.type
		(comment)? @alias.documentation
		(#not-eq? @alias.type "")
	)
`, anyTypeQuery)

func (l *Lexer) lexSimpleAlias(source *symbol.TypeSource) (symbol.Symbol, error) {
	match, err := l.context.Nvim.TsQueryOne(treesitter.Query{Language: "luadoc", Query: simpleAliasQuery})

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	typeCapture, typeCaptureFound := slicesx.FindFunc(*match, func(capture treesitter.Capture) bool {
		return capture.Id == "alias.type"
	})

	if !typeCaptureFound {
		return nil, fmt.Errorf("Error retrieving type value from alias annotation '%s'", source.Origin.DocumentationText())
	}

	lexedAliasType, err := l.lexType(typeCapture.Node.Text)

	if err != nil {
		return nil, err
	}

	return lexedAliasType, nil
}

var enumAliasQuery = fmt.Sprintf(`
	[
		(alias_annotation
			(identifier) @enumAlias.name
		)
		(continuation
			%s @enumAlias.type
		)
	]
`, anyTypeQuery)

func (l *Lexer) lexEnumAlias(_ *symbol.TypeSource) (symbol.Symbol, error) {
	matches, err := l.context.Nvim.TsQueryAll(treesitter.Query{Language: "luadoc", Query: enumAliasQuery})

	switch {
	case err != nil:
		return nil, err
	case matches == nil:
		return nil, nil
	}

	enumMembers := []string{}

	for _, match := range *matches {
		for _, capture := range match {
			switch capture.Id {
			case "enumAlias.type":
				enumMembers = append(enumMembers, capture.Node.Text)
			}
		}
	}

	fmt.Println(enumMembers)

	return nil, nil

}

func (l *Lexer) lexAlias(source *symbol.TypeSource) (symbol.Symbol, error) {
	fmt.Println(source.Path)

	if source.Origin.Type() != "alias_annotation" {
		return nil, nil
	}

	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	// Replacing all dots in the alias name with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	origin := source.Origin
	name := source.Path
	definitionText := origin.DefinitionText()
	patchedDefinitionText := strings.Replace(definitionText, name, strings.ReplaceAll(name, ".", "_"), 1)
	documentationText := origin.DocumentationText()
	patchedDocumentationLines := strings.Split(
		strings.Replace(documentationText, definitionText, patchedDefinitionText, 1),
		"\n",
	)

	buffer.SetLines(patchedDocumentationLines)

	simpleAliasType, err := l.lexSimpleAlias(source)

	switch {
	case err != nil:
		return nil, err
	case simpleAliasType != nil:
		return simpleAliasType, nil
	}

	enumAliasType, err := l.lexEnumAlias(source)

	switch {
	case err != nil:
		return nil, err
	case enumAliasType != nil:
		return enumAliasType, nil
	}

	return newBuiltinType(symbol.BuiltinValueVoid), nil
}
