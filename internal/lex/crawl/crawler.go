package crawl

import (
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type ctx interface {
	Target() *symbol.Target
	Nvim() *nvim.Nvim
	Logger() *log.Logger
}

type Crawler struct {
	context ctx
}

func (c *Crawler) GetOrigins() (*symbol.Origins, error) {
	var locations *[]nvim.Location
	var err error

	switch c.context.Target().Kind() {
	case symbol.TargetKindValue:
		locations, err = c.findDefinitionLocations(c.context.Target().Identifier())
	case symbol.TargetKindType:
		locations, err = c.findTypeDefinitionLocations(c.context.Target().Name(), c.context.Target().ParentName())
	}

	if err != nil {
		return nil, err
	}

	if locations == nil {
		c.context.Logger().Warnf("No locations found for '%s' %s symbol", c.context.Target().Identifier(), c.context.Target().Kind())
		return nil, nil
	}

	origins, err := c.getOrigins(*locations)

	if err != nil {
		return nil, err
	}

	if origins == nil {
		c.context.Logger().Warnf("No origin found for '%s' %s symbol", c.context.Target().Identifier(), c.context.Target().Kind())
		return nil, nil
	}

	return origins, nil
}

func (c *Crawler) GetMeta() (symbol.Meta, error) {
	origin := c.context.Target().Origins().First()

	return c.getMeta(origin)
}

func NewCrawler(context ctx) *Crawler {
	return &Crawler{
		context: context,
	}
}
