package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var EnumAnnotationQuery = annotation.AtEnumQuery

var EnumAnnotationMembersQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(table_constructor
		(field
			name: (identifier)
			value: [
				(number)
				(string)
			]
		)
	) @enum.members`,
}

var EnumAnnotationMemberQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(field
		name: (identifier) @enum.member.name
		value: [
			(number) @enum.member.value.number
			(string) @enum.member.value.string
		] @enum.member.value
	) @enum.member`,
}

type enumOriginMember struct {
	name  string
	value string
	type_ string
}

func (eom enumOriginMember) Name() string {
	return eom.name
}

func (eom enumOriginMember) Value() string {
	return eom.value
}

func (eom enumOriginMember) Type() string {
	return eom.type_
}

type EnumeratorAnnotationOrigin struct {
	origin
	definition string
	name       string
	members    []enumOriginMember
}

func (eo *EnumeratorAnnotationOrigin) Definition() []string {
	return strings.Split(eo.definition, "\n")
}

func (eo *EnumeratorAnnotationOrigin) Name() string {
	return eo.name
}

func (eo *EnumeratorAnnotationOrigin) Members() []enumOriginMember {
	return eo.members
}

func NewEnumAnnotationOrigin(location nvim.Location, enumMatch nvim.TsQueryMatch, memberMatches ...nvim.TsQueryMatch) *EnumeratorAnnotationOrigin {
	enumOrigin := EnumeratorAnnotationOrigin{
		origin: origin{
			location:    location,
			captures:    enumMatch,
			annotations: []string{},
		},
		members: []enumOriginMember{},
	}

	for _, capture := range enumMatch {
		switch capture.Id {
		case "enum.name":
			enumOrigin.name = capture.Node.Text
		}
	}

	for _, memberMatches := range memberMatches {
		for _, capture := range memberMatches {
			switch capture.Id {
			case "enum.members":
				enumOrigin.definition = capture.Node.Text
			case "enum.member.name":
				enumOrigin.members = append(enumOrigin.members, enumOriginMember{name: capture.Node.Text})
			case "enum.member.value.number":
				enumOrigin.members[len(enumOrigin.members)-1].value = capture.Node.Text
				enumOrigin.members[len(enumOrigin.members)-1].type_ = "number"
			case "enum.member.value.string":
				enumOrigin.members[len(enumOrigin.members)-1].value = capture.Node.Text
				enumOrigin.members[len(enumOrigin.members)-1].type_ = "string"
			}
		}
	}

	return &enumOrigin
}
