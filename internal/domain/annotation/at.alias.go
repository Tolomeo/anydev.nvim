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
	match string
	type_ string
}

func (aae *AtAliasEnumeratorMember) Match() string {
	return aae.match
}

func (aae *AtAliasEnumeratorMember) Type() string {
	return aae.type_
}

func NewAtAliasEnumeratorMember(match nvim.TsQueryMatch) *AtAliasEnumeratorMember {
	atAliasEnumeratorMember := AtAliasEnumeratorMember{}

	for _, capture := range match {
		switch capture.Id {
		case "alias.enumerator.member":
			atAliasEnumeratorMember.match = capture.Node.Text
		case "alias.type":
			atAliasEnumeratorMember.type_ = capture.Node.Text
		}
	}

	return &atAliasEnumeratorMember
}

type AtAliasEnumerator struct {
	match string
	name  string
}

func (aae *AtAliasEnumerator) Match() string {
	return aae.match
}

func (aae *AtAliasEnumerator) Name() string {
	return aae.name
}

func NewAtAliasEnumerator(match nvim.TsQueryMatch) *AtAliasEnumerator {
	atAliasEnumerator := AtAliasEnumerator{}

	for _, capture := range match {
		switch capture.Id {
		case "alias.enumerator":
			atAliasEnumerator.match = capture.Node.Text
		case "alias.name":
			atAliasEnumerator.name = capture.Node.Text
		}
	}

	return &atAliasEnumerator
}
