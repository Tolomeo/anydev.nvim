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
var simpleAliasQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(alias_annotation
		(identifier) @alias.name
		%s @alias.type
		(comment)? @alias.documentation
		(#not-eq? @alias.type "")
	)`, anyTypeQuery),
}

func (l *Lexer) lexSimpleAlias(source *symbol.TypeSource) (symbol.Symbol, error) {
	match, err := l.context.Nvim.TsQueryOne(simpleAliasQuery)

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

var enumAliasQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(continuation
		%s @enumAlias.type
	)
`, anyTypeQuery)}

func (l *Lexer) lexEnumAlias(source *symbol.TypeSource) (*symbol.Union, error) {
	matches, err := l.context.Nvim.TsQueryAll(enumAliasQuery)

	fmt.Println(matches, err)

	switch {
	case err != nil:
		return nil, err
	case matches == nil:
		return nil, nil
	}

	enumMembers := []symbol.Symbol{}

	for _, matchCaptures := range *matches {
		for _, capture := range matchCaptures {
			fmt.Println("capture")
			fmt.Println(source.Path)
			fmt.Println(capture)

			switch capture.Id {
			case "enumAlias.type":
				lexedType, err := l.lexType(capture.Node.Text)

				fmt.Println("Lexed union type")
				fmt.Println(lexedType, err)

				if err != nil {
					return nil, err
				}

				enumMembers = append(enumMembers, lexedType)
			}
		}
	}

	fmt.Println("after loop")
	fmt.Println(source.Path)
	// fmt.Println(source.Path)
	/* for _, member := range enumMembers {
		fmt.Println(member)

	} */

	if len(enumMembers) < 1 {
		return nil, fmt.Errorf("Could not retrieve enum members from enum alias '%s'", source.Path)
	}

	return newUnionType(enumMembers...), nil
}

func (l *Lexer) lexAlias(source *symbol.TypeSource) (symbol.Symbol, error) {
	// fmt.Println(source.Path)

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

	fmt.Println(source.Path)
	if source.Path == "vim.validate.Validator" {
		lines, _ := buffer.ReadLines()
		fmt.Println(lines)
		fmt.Println(simpleAliasType)
	}

	switch {
	case err != nil:
		return nil, err
	case simpleAliasType != nil:
		return simpleAliasType, nil
	}

	enumAliasType, err := l.lexEnumAlias(source)

	if source.Path == "vim.validate.Validator" {
		fmt.Println(enumAliasType)
	}
	// fmt.Println(enumAliasType)

	switch {
	case err != nil:
		return nil, err
	case enumAliasType != nil:
		return enumAliasType, nil
	}

	return nil, fmt.Errorf("Could not lex alias type '%s'", source.Path)
}
