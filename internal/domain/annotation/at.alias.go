package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AtAliasQuery = treesitter.Query{
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
	) @alias`, AnyTypeQuery),
}

type AtAlias struct {
	match string
	name  string
	type_ string
}

func (ao *AtAlias) Match() string {
	return ao.match
}

func (ao *AtAlias) Name() string {
	return ao.name
}

func (ao *AtAlias) Type() string {
	return ao.type_
}

func NewAtAlias(match nvim.TsQueryMatch) *AtAlias {
	atAlias := AtAlias{}

	for _, capture := range match {
		switch capture.Id {
		case "alias":
			atAlias.match = capture.Node.Text
		case "alias.name":
			atAlias.name = capture.Node.Text
		case "alias.type":
			atAlias.type_ = capture.Node.Text
		}
	}

	return &atAlias
}
