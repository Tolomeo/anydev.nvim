package crawl

import (
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
	err := c.sourceType(name, &source)

	if err != nil {
		return nil, err
	}

	return &source, nil
}

func NewCrawler(context *context.Context) *Crawler {
	return &Crawler{
		context: context,
	}
}
