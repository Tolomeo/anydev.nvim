package crawl

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type logger interface {
	Info(message string)
	Warn(message string)
	Error(message string)
}

type Crawler struct {
	config *CrawlerConfig
}

func (c *Crawler) Source(path string) (*Source, error) {
	source := Source{path: path}
	err := c.sourceRuntime(path, &source)

	if err != nil {
		return nil, err
	}

	return &source, nil
}

func (c *Crawler) SourceType(path string) (*Source, error) {
	source := Source{path: path}
	err := c.sourceType(path, &source)

	if err != nil {
		return nil, err
	}

	return &source, nil
}

type CrawlerConfig struct {
	nvim   *nvim.Nvim
	logger logger
}

func NewCrawlerConfig(logger logger, nvim *nvim.Nvim) *CrawlerConfig {
	return &CrawlerConfig{
		logger: logger,
		nvim:   nvim,
	}
}

func NewCrawler(config *CrawlerConfig) *Crawler {
	return &Crawler{
		config: config,
	}
}
