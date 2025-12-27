package crawl

import (
	"github.com/Tolomeo/anydev.nvim/internal/context"
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

type Crawler struct {
	context *context.Context
}

func (c *Crawler) SourceValue(path string) (*symbol.Source, error) {
	source := symbol.Source{Path: path}
	err := c.sourceValue(path, &source)

	if err != nil {
		return nil, err
	}

	return &source, nil
}

func (c *Crawler) SourceType(name string) (*symbol.Source, error) {
	source := symbol.Source{Path: name}
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
