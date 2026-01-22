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

func (c *Crawler) Crawl() (symbol.Origin, error) {
	var locations *[]nvim.Location
	var err error

	switch c.target.Kind() {
	case symbol.TargetKindValue:
		locations, err = c.findDefinitionLocations(c.target.Identifier())
	case symbol.TargetKindType:
		locations, err = c.findTypeDefinitionLocations(c.target.Name(), c.target.ParentName())
	}

	if err != nil {
		return nil, err
	}

	if locations == nil {
		c.target.Logger().Warnf("No locations found for '%s' %s symbol", c.target.Identifier(), c.target.Kind())
		return nil, nil
	}

	origin, err := c.findOrigin(*locations)

	if err != nil {
		return nil, err
	}

	if origin == nil {
		c.target.Logger().Warnf("No origin found for '%s' %s symbol", c.target.Identifier(), c.target.Kind())
		return nil, nil
	}

	return origin, nil
}

func NewCrawler(target target) *Crawler {
	return &Crawler{
		target: target,
	}
}
