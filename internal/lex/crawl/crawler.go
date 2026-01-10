package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/context"
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

type Crawler struct {
	context *context.Context
}

func (c *Crawler) SourceValue(path string) (symbol.Source, error) {
	source := symbol.ValueSource{Path: path}
	locations, err := c.findValueDefinitionLocations(source.Identifier())

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.context.Logger.Warn(fmt.Sprintf("No locations found for '%s' symbol", source.Identifier()))
		return nil, nil
	}

	pathOrigin, err := c.findValueOrigin(&source, *locations)

	switch {
	case err != nil:
		return nil, err
	case pathOrigin == nil:
		c.context.Logger.Warn(fmt.Sprintf("No origin found for '%s' symbol", source.Identifier()))
		return nil, nil
	default:
		source.SetOrigin(pathOrigin)
	}

	/* err = c.followValueOrigin(&source)

	if err != nil {
		return nil, err
	} */

	documentation, err := c.sourceValueOriginDocumentation(&source)

	switch {
	case err != nil:
		return nil, err
	case documentation == nil:
		c.context.Logger.Warn(fmt.Sprintf("No documentation found for '%s' symbol", source.Identifier()))
		return &source, nil
	}

	source.GetOrigin().SetDocumentation(*documentation)
	return &source, nil
}

func (c *Crawler) SourceType(name string) (symbol.Source, error) {
	source := symbol.TypeSource{Path: name}
	/* err := c.sourceType(&source)

	if err != nil {
		return nil, err
	}

	return &source, nil */

	locations, err := c.getTypeDefinitionLocations(source.Path)

	/* for _, loc := range *locations {
		fmt.Printf("\n\nLocation: %+v\n", loc)
	} */

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.context.Logger.Warn(fmt.Sprintf("No locations found for '%s' type", source.Path))
		return nil, nil
	}

	origin, err := c.findTypeOrigin(&source, *locations)

	switch {
	case err != nil:
		return nil, err
	case origin == nil:
		c.context.Logger.Warn(fmt.Sprintf("No origin found for '%s' symbol", source.Path))
		return nil, nil
	}

	source.SetOrigin(origin)

	documentation, err := c.sourceTypeOriginDocumentation(&source)

	switch {
	case err != nil:
		return nil, err
	case documentation == nil:
		c.context.Logger.Warn(fmt.Sprintf("No documentation found for '%s' symbol", source.Identifier()))
		return &source, nil
	}

	source.GetOrigin().SetDocumentation(*documentation)
	return &source, nil
}

func (c *Crawler) SourceTypeMember(name string, member string) (symbol.Source, error) {
	source := symbol.ValueSource{Path: fmt.Sprintf("%s.%s", name, member)}

	locations, err := c.getTypeDefinitionLocations(name, member)

	/* for _, loc := range *locations {
		fmt.Printf("\n\nLocation: %+v\n", loc)
	} */

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.context.Logger.Warn(fmt.Sprintf("No locations found for '%s' type", source.Path))
		return nil, nil
	}

	memberOrigin, err := c.findValueOrigin(&source, *locations)

	switch {
	case err != nil:
		return nil, err
	case memberOrigin == nil:
		return nil, nil
	default:
		source.SetOrigin(memberOrigin)
	}

	documentation, err := c.sourceValueOriginDocumentation(&source)

	switch {
	case err != nil:
		return nil, err
	case documentation == nil:
		c.context.Logger.Warn(fmt.Sprintf("No documentation found for '%s' symbol", source.Identifier()))
		return nil, nil
	default:
		source.GetOrigin().SetDocumentation(*documentation)
	}

	return &source, nil
}

func NewCrawler(context *context.Context) *Crawler {
	return &Crawler{
		context: context,
	}
}
