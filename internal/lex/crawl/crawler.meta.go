package crawl

import (
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var metaAnnotationsQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(documentation 
		[
			(qualifier_annotation "@private") @private
			(qualifier_annotation "@protected") @protected
			(package_annotation) @package
			(deprecated_annotation) @deprecated
		]
	)`,
}

func (c *Crawler) getMeta(origin symbol.Origin) (symbol.Meta, error) {
	meta := symbol.Meta{}
	buffer, err := c.context.Nvim().NewBuffer()

	if err != nil {
		return meta, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.Documentation())

	if err != nil {
		return meta, nil
	}

	matches, err := buffer.TsQueryAll(metaAnnotationsQuery)

	if err != nil {
		return meta, nil
	}

	if matches == nil {
		return meta, nil
	}

	for _, match := range *matches {
		for _, capture := range match {
			switch capture.Id {
			case "private":
				meta.Private = true
			case "protected":
				meta.Protected = true
			case "package":
				meta.Package = true
			case "deprecated":
				meta.Deprecated = true
			}
		}
	}

	return meta, nil
}
