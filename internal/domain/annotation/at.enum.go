package annotation

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AtEnumMembersQuery = treesitter.Query{
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

var AtEnumMemberQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(field
		name: (identifier) @enum.member.name
		value: [
			(number)
			(string)
		] @enum.member.value
	) @enum.member`,
}

type AtEnumMember struct {
	name  string
	value string
	type_ string
}

func (aem AtEnumMember) Name() string {
	return aem.name
}

func (aem AtEnumMember) Value() string {
	return aem.value
}

func (aem AtEnumMember) Type() string {
	return aem.type_
}

func NewAtEnumMember(captures nvim.TsQueryMatch) *AtEnumMember {
	atEnumMember := AtEnumMember{}

	for _, capture := range captures {
		switch capture.Id {
		case "enum.member.name":
			atEnumMember.name = capture.Node.Text
		case "enum.member.value":
			atEnumMember.value = capture.Node.Text
			atEnumMember.type_ = capture.Node.Type
		}
	}

	return &atEnumMember
}

var AtEnumQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(enum_annotation
		"@enum"
		.
		(identifier) @enum.name
	) @enum`,
}

type AtEnum struct {
	text string
	name string
}

func (ae *AtEnum) Text() string {
	return ae.text
}

func (ae *AtEnum) Name() string {
	return ae.name
}

func NewAtEnum(captures nvim.TsQueryMatch) *AtEnum {
	atEnum := AtEnum{}

	for _, capture := range captures {
		switch capture.Id {
		case "enum":
			atEnum.text = capture.Node.Text
		case "enum.name":
			atEnum.name = capture.Node.Text
		}
	}

	return &atEnum
}
