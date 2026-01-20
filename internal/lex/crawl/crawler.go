package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type target interface {
	symbol.Target
	Nvim() *nvim.Nvim
	Logger() *log.Logger
}

type Crawler struct {
	target target
}

func (c *Crawler) SourceValue() (symbol.Source, error) {
	locations, err := c.findDefinitionLocations(c.target.Identifier())

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.target.Logger().Warn(fmt.Sprintf("No locations found for '%s' symbol", c.target.Identifier()))
		return nil, nil
	}

	source := symbol.NewValueSource(c.target.Identifier())
	origin, err := c.findOrigin(*locations)

	switch {
	case err != nil:
		return nil, err
	case origin == nil:
		c.target.Logger().Warn(fmt.Sprintf("No origin found for '%s' symbol", c.target.Identifier()))
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
		c.target.Logger().Warn(fmt.Sprintf("No documentation found for '%s' symbol", c.target.Identifier()))
		return source, nil
	}

	source.GetOrigin().SetDocumentation(*documentation)
	return source, nil
}

func (c *Crawler) SourceType() (symbol.Source, error) {
	typeName, parentTypeName := c.target.Name(), c.target.ParentName()

	source := symbol.NewTypeSource(typeName, parentTypeName)
	locations, err := c.findTypeDefinitionLocations(c.target.Name(), c.target.ParentName())

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.target.Logger().Warn(fmt.Sprintf("No locations found for '%s' type", c.target.Identifier()))
		return nil, nil
	}

	origin, err := c.findOrigin(*locations)

	switch {
	case err != nil:
		return nil, err
	case origin == nil:
		c.target.Logger().Warn(fmt.Sprintf("No origin found for '%s' symbol", c.target.Identifier()))
		return nil, nil
	default:
		source.SetOrigin(origin)
	}

	documentation, err := c.sourceDocumentation(origin)

	switch {
	case err != nil:
		return nil, err
	case documentation == nil:
		c.target.Logger().Warn(fmt.Sprintf("No documentation found for '%s' symbol", source.Identifier()))
		return source, nil
	default:
		source.GetOrigin().SetDocumentation(*documentation)
	}

	return source, nil
}

func NewCrawler(target target) *Crawler {
	return &Crawler{
		target: target,
	}
}
