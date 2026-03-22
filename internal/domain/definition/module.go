package definition

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var ModuleExportQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(chunk
		(return_statement
			(expression_list
				(identifier) @module.value
			)
		) @module.export
		.
	)`,
}

type ModuleExport struct {
	text   string
	range_ treesitter.Range
}

func (me *ModuleExport) Range() treesitter.Range {
	return me.range_
}

func NewModuleExport(match nvim.TsQueryMatch) *ModuleExport {
	moduleExport := ModuleExport{}

	for _, capture := range match {
		switch capture.Id {
		case "module.export":
			moduleExport.text = capture.Node.Text
			moduleExport.range_ = capture.Node.Range
		}
	}

	return &moduleExport
}
