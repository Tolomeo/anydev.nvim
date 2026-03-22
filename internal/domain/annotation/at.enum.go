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
	name  treesitter.TsNode
	value treesitter.TsNode
	type_ string
}

func (aem AtEnumMember) Name() treesitter.TsNode {
	return aem.name
}

func (aem AtEnumMember) Value() treesitter.TsNode {
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
			atEnumMember.name = capture.Node
		case "enum.member.value":
			atEnumMember.value = capture.Node
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
	root treesitter.TsNode
	name treesitter.TsNode
}

func (ae *AtEnum) Root() treesitter.TsNode {
	return ae.root
}

func (ae *AtEnum) Name() treesitter.TsNode {
	return ae.name
}

func NewAtEnum(captures nvim.TsQueryMatch) *AtEnum {
	atEnum := AtEnum{}

	for _, capture := range captures {
		switch capture.Id {
		case "enum":
			atEnum.root = capture.Node
		case "enum.name":
			atEnum.name = capture.Node
		}
	}

	return &atEnum
}
