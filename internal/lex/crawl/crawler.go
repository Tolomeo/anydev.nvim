package crawl

import (
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

func (c *Crawler) SourceValue() (symbol.Origin, error) {
	locations, err := c.findDefinitionLocations(c.target.Identifier())

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.target.Logger().Warnf("No locations found for '%s' symbol", c.target.Identifier())
		return nil, nil
	}

	origin, err := c.findOrigin(*locations)

	switch {
	case err != nil:
		return nil, err
	case origin == nil:
		c.target.Logger().Warnf("No origin found for '%s' symbol", c.target.Identifier())
		return nil, nil
	}

	documentation, err := c.sourceDocumentation(origin)

	switch {
	case err != nil:
		return nil, err
	case documentation == nil:
		c.target.Logger().Warnf("No documentation found for '%s' symbol", c.target.Identifier())
		return origin, nil
	}

	origin.SetDocumentation(*documentation)
	return origin, nil
}

func (c *Crawler) SourceType() (symbol.Origin, error) {
	locations, err := c.findTypeDefinitionLocations(c.target.Name(), c.target.ParentName())

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.target.Logger().Warnf("No locations found for '%s' type", c.target.Identifier())
		return nil, nil
	}

	origin, err := c.findOrigin(*locations)

	switch {
	case err != nil:
		return nil, err
	case origin == nil:
		c.target.Logger().Warnf("No origin found for '%s' symbol", c.target.Identifier())
		return nil, nil
	}

	documentation, err := c.sourceDocumentation(origin)

	switch {
	case err != nil:
		return nil, err
	case documentation == nil:
		c.target.Logger().Warnf("No documentation found for '%s' symbol", c.target.Identifier())
		return origin, nil
	}

	origin.SetDocumentation(*documentation)

	return origin, nil
}

func NewCrawler(target target) *Crawler {
	return &Crawler{
		target: target,
	}
}
