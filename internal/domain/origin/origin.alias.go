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
}

func (ao *AliasOrigin) Definition() []string {
	root, _ := ao.definition.Match.Find("alias")
	return strings.Split(root.Node.Text, "\n")
}

func NewAliasOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *AliasOrigin {
	return &AliasOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
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
}

func (aeo *AliasEnumeratorOrigin) Definition() []string {
	root, _ := aeo.definition.Match.Find("alias.enumerator")
	return strings.Split(root.Node.Text, "\n")
}

func NewAliasEnumeratorOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *AliasEnumeratorOrigin {
	return &AliasEnumeratorOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}
