package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/context"
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

type Crawler struct {
	context *context.Context
}

func (c *Crawler) SourceValue(path string) (*symbol.ValueSource, error) {
	source := symbol.ValueSource{Path: path}
	err := c.sourceValueOrigin(&source)

	if err != nil {
		return nil, err
	}

	return &source, nil
}

func (c *Crawler) SourceType(name string) (*symbol.TypeSource, error) {
	source := symbol.TypeSource{Path: name}
	err := c.sourceType(&source)

	if err != nil {
		return nil, err
	}

	return &source, nil
}

func (c *Crawler) SourceTypeMember(name string, member string) (*symbol.ValueSource, error) {
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

	memberOrigin, err := c.findValueOrigin(*locations)

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
		source.GetOrigin().SetDocumentation(documentation)

	}

	return &source, nil
}

func NewCrawler(context *context.Context) *Crawler {
	return &Crawler{
		context: context,
	}
}
