package annotation

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

/* var AtClassQuery = treesitter.Query{
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
} */

var AtClassQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(class_annotation
		"@class"
		.
		"(exact)"? @class.exact
		.
		(identifier) @class.name
	) @class`,
}

type AtClass struct {
	root  treesitter.TsNode
	exact *treesitter.TsNode
	name  treesitter.TsNode
	// parents []string
}

func (c *AtClass) Root() treesitter.TsNode {
	return c.root
}

func (c *AtClass) Exact() *treesitter.TsNode {
	return c.exact
}

func (c *AtClass) Name() treesitter.TsNode {
	return c.name
}

/* func (c *AtClass) Parents() []string {
	return c.parents
} */

func NewAtClass(captures nvim.TsQueryMatch) *AtClass {
	class := AtClass{
		// parents: []string{},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "class":
			class.root = capture.Node
		case "class.exact":
			class.exact = &capture.Node
		case "class.name":
			class.name = capture.Node
			/* case "class.parent":
			class.parents = append(class.parents, capture.Node.Text) */
		}
	}

	return &class
}
