package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var ClassAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(class_annotation
		"@class"
		.
		"(exact)"? @class.exact
		.
		(identifier) @class.name
		.
		(":" 
			. (%s) @class.parent
			("," (%s) @class.parent)*
		)?
	) @class`, AnyTypeQuery, AnyTypeQuery),
}

type AtClass struct {
	match   string
	exact   bool
	name    string
	parents []string
}

func (c *AtClass) Match() string {
	return c.match
}

func (c *AtClass) Exact() bool{
	return c.exact
}

func (c *AtClass) Name() string {
	return c.name
}

func (c *AtClass) Parents() []string {
	return c.parents
}

func NewClass(captures nvim.TsQueryMatch) *AtClass {
	class := AtClass{
		parents: []string{},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "class":
			class.match = capture.Node.Text
		case "class.exact":
			class.exact = true
		case "class.name":
			class.name = capture.Node.Text
		case "class.parent":
			class.parents = append(class.parents, capture.Node.Text)
		}
	}

	return &class
}
