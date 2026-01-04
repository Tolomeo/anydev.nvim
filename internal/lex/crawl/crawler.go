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

	fmt.Println(name, member)
	fmt.Println(locations)
	fmt.Println()
	/* pathOrigin, err := c.findValueOrigin(*locations)

	switch {
	case err != nil:
		return err
	case pathOrigin == nil:
		c.context.Logger.Warn(fmt.Sprintf("No origin found for '%s' symbol", source.Path))
		return nil
	default:
		source.Origin = pathOrigin
	}

	source.Origin = origin */
	return nil, nil
}

func NewCrawler(context *context.Context) *Crawler {
	return &Crawler{
		context: context,
	}
}
