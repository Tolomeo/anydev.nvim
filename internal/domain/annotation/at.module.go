package annotation

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AtModuleQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(module_annotation
	 "@module"
	 .
	 (string) @module.name
	) @module`,
}

type AtModule struct {
	root treesitter.TsNode
	name treesitter.TsNode
}

func (am *AtModule) Root() treesitter.TsNode {
	return am.root
}

func (am *AtModule) Name() treesitter.TsNode {
	return am.name
}

func NewAtModule(match nvim.TsQueryMatch) *AtModule {
	atModule := AtModule{}

	for _, capture := range match {
		switch capture.Id {
		case "module":
			atModule.root = capture.Node
		case "module.name":
			atModule.name = capture.Node
		}
	}

	return &atModule
}
