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
}

type Crawler struct {
	context context
}

func (c *Crawler) SourceValue(path string) (symbol.Source, error) {
	source := symbol.NewValueSource(path)

	locations, err := c.findDefinitionLocations(source.Identifier())

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.context.Logger().Warn(fmt.Sprintf("No locations found for '%s' symbol", source.Identifier()))
		return nil, nil
	}

	valueOrigin, err := c.findOrigin(*locations, source)

	switch {
	case err != nil:
		return nil, err
	case valueOrigin == nil:
		c.context.Logger().Warn(fmt.Sprintf("No origin found for '%s' symbol", source.Identifier()))
		return nil, nil
	default:
		source.SetOrigin(valueOrigin)
	}

	/* err = c.followValueOrigin(&source)

	if err != nil {
		return nil, err
	} */

	documentation, err := c.sourceDocumentation(source)

	switch {
	case err != nil:
		return nil, err
	case documentation == nil:
		c.context.Logger().Warn(fmt.Sprintf("No documentation found for '%s' symbol", source.Identifier()))
		return source, nil
	}

	source.GetOrigin().SetDocumentation(*documentation)
	return source, nil
}

func (c *Crawler) SourceType(typeName string, parentTypeName string) (symbol.Source, error) {
	source := symbol.NewTypeSource(typeName, parentTypeName)

	locations, err := c.findTypeDefinitionLocations(source)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.context.Logger().Warn(fmt.Sprintf("No locations found for '%s' type", source.Identifier()))
		return nil, nil
	}

	origin, err := c.findOrigin(*locations, source)

	switch {
	case err != nil:
		return nil, err
	case origin == nil:
		c.context.Logger().Warn(fmt.Sprintf("No origin found for '%s' symbol", source.Identifier()))
		return nil, nil
	default:
		source.SetOrigin(origin)
	}

	documentation, err := c.sourceDocumentation(source)

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
