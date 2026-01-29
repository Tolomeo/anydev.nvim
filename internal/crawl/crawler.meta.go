package crawl

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
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

func (c *Crawler) getMeta(symbolOrigin origin.Origin) (*symbol.Meta, error) {
	meta := symbol.NewMeta()

	/* switch symbolOriginType := symbolOrigin.(type) {
	case *origin.FieldOrigin:
		meta.SetPrivate(symbolOriginType.Private())
		meta.SetProtected(symbolOriginType.Protected())
		meta.SetPackage(symbolOriginType.Package())
		// meta.SetDeprecated(symbolOriginType.Deprecated())
		return meta, nil
	} */

	buffer, err := c.context.Nvim().NewBuffer()

	if err != nil {
		return meta, err
	}

	defer buffer.Close()

	err = buffer.SetLines(symbolOrigin.Documentation())

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
				meta.SetPrivate(true)
			case "protected":
				meta.SetProtected(true)
			case "package":
				meta.SetPackage(true)
			case "deprecated":
				meta.SetDeprecated(true)
			}
		}
	}

	return meta, nil
}
