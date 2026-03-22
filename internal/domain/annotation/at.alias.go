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
	root  treesitter.TsNode
	name  treesitter.TsNode
	type_ treesitter.TsNode
}

func (ao *AtAlias) Root() treesitter.TsNode {
	return ao.root
}

func (ao *AtAlias) Name() treesitter.TsNode {
	return ao.name
}

func (ao *AtAlias) Type() treesitter.TsNode {
	return ao.type_
}

func NewAtAlias(match nvim.TsQueryMatch) *AtAlias {
	atAlias := AtAlias{}

	for _, capture := range match {
		switch capture.Id {
		case "alias":
			atAlias.root = capture.Node
		case "alias.name":
			atAlias.name = capture.Node
		case "alias.type":
			atAlias.type_ = capture.Node
		}
	}

	return &atAlias
}

// Luadoc matches an empty type node even when the type is not present
// That means that enum aliases have an empty type node defined
var AtAliasEnumeratorQuery = treesitter.Query{
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
	) @alias.enumerator`, AnyTypeQuery),
}

var AtAliasEnumeratorMemberQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(continuation
		(%s) @alias.type
	) @alias.enumerator.member`, AnyTypeQuery),
}

type AtAliasEnumeratorMember struct {
	root  treesitter.TsNode
	type_ treesitter.TsNode
}

func (aae *AtAliasEnumeratorMember) Root() treesitter.TsNode {
	return aae.root
}

func (aae *AtAliasEnumeratorMember) Type() treesitter.TsNode {
	return aae.type_
}

func NewAtAliasEnumeratorMember(match nvim.TsQueryMatch) *AtAliasEnumeratorMember {
	atAliasEnumeratorMember := AtAliasEnumeratorMember{}

	for _, capture := range match {
		switch capture.Id {
		case "alias.enumerator.member":
			atAliasEnumeratorMember.root = capture.Node
		case "alias.type":
			atAliasEnumeratorMember.type_ = capture.Node
		}
	}

	return &atAliasEnumeratorMember
}

type AtAliasEnumerator struct {
	root treesitter.TsNode
	name treesitter.TsNode
}

func (aae *AtAliasEnumerator) Root() treesitter.TsNode {
	return aae.root
}

func (aae *AtAliasEnumerator) Name() treesitter.TsNode {
	return aae.name
}

func NewAtAliasEnumerator(match nvim.TsQueryMatch) *AtAliasEnumerator {
	atAliasEnumerator := AtAliasEnumerator{}

	for _, capture := range match {
		switch capture.Id {
		case "alias.enumerator":
			atAliasEnumerator.root = capture.Node
		case "alias.name":
			atAliasEnumerator.name = capture.Node
		}
	}

	return &atAliasEnumerator
}
