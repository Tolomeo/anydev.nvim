package origin

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AliasAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`(
		(alias_annotation
			"@alias"
			.
			(identifier) @alias.name
			.
			(%s) @alias.type
			.
			(comment)? @alias.documentation
			.
		)
		(#not-eq? @alias.type "")
	) @alias`, anyTypeQuery),
}

type AliasOrigin struct {
	origin
	definition string
	name       string
	type_      string
}

func (ao *AliasOrigin) Definition() []string {
	return strings.Split(ao.definition, "\n")
}

func (ao *AliasOrigin) Name() string {
	return ao.name
}

func (ao *AliasOrigin) Type() string {
	return ao.type_
}

func NewAliasOrigin(location nvim.Location, captures nvim.TsQueryMatch, documentation []string) *AliasOrigin {
	aliasOrigin := AliasOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: documentation,
		},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "alias":
			aliasOrigin.definition = capture.Node.Text
		case "alias.name":
			aliasOrigin.name = capture.Node.Text
		case "alias.type":
			aliasOrigin.type_ = capture.Node.Text
		}
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
	) @alias.enumerator`, anyTypeQuery),
}

var AliasEnumeratorMemberAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(continuation
		(%s) @alias.type
	)
`, anyTypeQuery),
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

func NewAliasEnumeratorOrigin(location nvim.Location, captures nvim.TsQueryMatch, documentation []string) *AliasEnumeratorOrigin {
	aliasEnumeratorOrigin := AliasEnumeratorOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: documentation,
		},
		types: []string{},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "alias.enumerator":
			aliasEnumeratorOrigin.definition = capture.Node.Text
		case "alias.name":
			aliasEnumeratorOrigin.name = capture.Node.Text
		case "alias.type":
			aliasEnumeratorOrigin.types = append(aliasEnumeratorOrigin.types, capture.Node.Text)
		}
	}

	return &aliasEnumeratorOrigin
}
