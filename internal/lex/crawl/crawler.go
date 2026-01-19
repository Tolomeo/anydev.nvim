package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type context interface {
	Nvim() *nvim.Nvim 
	Logger() *log.Logger
	Target() symbol.Target
}

type Crawler struct {
	context context
}

func (c *Crawler) SourceValue() (symbol.Source, error) {
	target := c.context.Target()
	locations, err := c.findDefinitionLocations(target.Identifier())

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.context.Logger().Warn(fmt.Sprintf("No locations found for '%s' symbol", target.Identifier()))
		return nil, nil
	}

	source := symbol.NewValueSource(target.Identifier())
	origin, err := c.findOrigin(*locations)

	switch {
	case err != nil:
		return nil, err
	case origin == nil:
		c.context.Logger().Warn(fmt.Sprintf("No origin found for '%s' symbol", target.Identifier()))
		return nil, nil
	default:
		source.SetOrigin(origin)
	}

	/* err = c.followValueOrigin(&source)

	if err != nil {
		return nil, err
	} */

	documentation, err := c.sourceDocumentation(origin)

	switch {
	case err != nil:
		return nil, err
	case documentation == nil:
		c.context.Logger().Warn(fmt.Sprintf("No documentation found for '%s' symbol", target.Identifier()))
		return source, nil
	}

	source.GetOrigin().SetDocumentation(*documentation)
	return source, nil
}

func (c *Crawler) SourceType() (symbol.Source, error) {
	target := c.context.Target()
	typeName, parentTypeName := target.Name(), target.ParentName()

	source := symbol.NewTypeSource(typeName, parentTypeName)
	locations, err := c.findTypeDefinitionLocations(target.Name(), target.ParentName())

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.context.Logger().Warn(fmt.Sprintf("No locations found for '%s' type", target.Identifier()))
		return nil, nil
	}

	origin, err := c.findOrigin(*locations)

	switch {
	case err != nil:
		return nil, err
	case origin == nil:
		c.context.Logger().Warn(fmt.Sprintf("No origin found for '%s' symbol", target.Identifier()))
		return nil, nil
	default:
		source.SetOrigin(origin)
	}

	documentation, err := c.sourceDocumentation(origin)

	switch {
	case err != nil:
		return nil, err
	case documentation == nil:
		c.context.Logger().Warn(fmt.Sprintf("No documentation found for '%s' symbol", source.Identifier()))
		return source, nil
	default:
		source.GetOrigin().SetDocumentation(*documentation)
	}

	return source, nil
}

func NewCrawler(context context) *Crawler {
	return &Crawler{
		context: context,
	}
}
