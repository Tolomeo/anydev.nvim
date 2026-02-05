package origin

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AliasAnnotationQuery = annotation.AtAliasQuery

type AliasOrigin struct {
	origin
	aliasAnnotation annotation.AtAlias
}

func (ao *AliasOrigin) Definition() []string {
	return strings.Split(ao.aliasAnnotation.Match(), "\n")
}

func (ao *AliasOrigin) Name() string {
	return ao.aliasAnnotation.Name()
}

func (ao *AliasOrigin) Type() string {
	return ao.aliasAnnotation.Type()
}

func NewAliasOrigin(location nvim.Location, captures nvim.TsQueryMatch) *AliasOrigin {
	aliasOrigin := AliasOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: []string{},
		},
		aliasAnnotation: *annotation.NewAtAlias(captures),
	}

	return &aliasOrigin
}

// Luadoc matches an empty type node even when the type is not present
// That means that enum aliases have an empty type node defined
var AliasEnumeratorAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`(
		(alias_annotation
			"@alias"
			.
			(identifier) @alias.name
			.
			(%s) @alias.emptytype
			.
			(comment)? @alias.documentation
			.
		)
		(#eq? @alias.emptytype "")
	) @alias.enumerator`, annotation.AnyTypeQuery),
}

type AliasEnumeratorOrigin struct {
	origin
	definition string
	name       string
	types      []string
}

func (aeo *AliasEnumeratorOrigin) Definition() []string {
	return strings.Split(aeo.definition, "\n")
}

func (aeo *AliasEnumeratorOrigin) Name() string {
	return aeo.name
}

func (aeo *AliasEnumeratorOrigin) Types() []string {
	return aeo.types
}

func NewAliasEnumeratorOrigin(location nvim.Location, aliasCaptures nvim.TsQueryMatch, membersCaptures ...nvim.TsQueryMatch) *AliasEnumeratorOrigin {
	aliasEnumeratorOrigin := AliasEnumeratorOrigin{
		origin: origin{
			location:    location,
			captures:    aliasCaptures,
			annotations: []string{},
		},
		types: []string{},
	}

	for _, capture := range aliasCaptures {
		switch capture.Id {
		case "alias.enumerator":
			aliasEnumeratorOrigin.definition = capture.Node.Text
		case "alias.name":
			aliasEnumeratorOrigin.name = capture.Node.Text
		}
	}

	for _, memberCaptures := range membersCaptures {
		for _, capture := range memberCaptures {
			switch capture.Id {
			case "alias.type":
				aliasEnumeratorOrigin.types = append(aliasEnumeratorOrigin.types, capture.Node.Text)
			}
		}
	}

	return &aliasEnumeratorOrigin
}

var AliasEnumeratorMemberAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(continuation
		(%s) @alias.type
	)`, annotation.AnyTypeQuery),
}
