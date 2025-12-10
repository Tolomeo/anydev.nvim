package crawl

import (
	"errors"
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/ts"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type logger interface {
	Info(message string)
	Warn(message string)
	Error(message string)
}

type Crawler struct {
	config *CrawlerConfig
}

func (c *Crawler) scratch(lines []string) error {
	buffer := path.Join(c.config.nvim.Options().Config().Dir(), "anydev.crawler.lua")

	_, err := c.config.nvim.Open(buffer)

	if err != nil {
		return err
	}

	err = c.config.nvim.SetBufferLines(lines)

	if err != nil {
		return err
	}

	return nil
}

func (c *Crawler) CrawlRuntime(path string) (Source, error) {
	source, err := c.getRuntimeSource(path)

	if err != nil {
		return nil, err
	}

	return source, nil
}

func (c *Crawler) getRuntimeSource(path string) (Source, error) {
	luaType, err := c.config.nvim.GetLuaTypeName(path)

	if err != nil {
		return nil, err
	}

	switch luaType {
	case "table":
		tableOrigin, err := c.getOrigin(path, ts.ASSIGNMENT_STATEMENT)

		if err != nil {
			return nil, err
		}

		tableSource := TableSource{
			path:   path,
			origin: tableOrigin,
		}

		fields, err := c.config.nvim.GetCompletion(path)

		if err != nil {
			return nil, err
		}

		for _, field := range fields {
			child, err := c.CrawlRuntime(path + "." + field)

			if err != nil {
				return nil, err
			}

			tableSource.fields = append(tableSource.fields, &child)
		}

		return &tableSource, nil

	case "function":
		functionOrigin, err := c.getOrigin(path, ts.FUNCTION_DECLARATION, ts.ASSIGNMENT_STATEMENT)

		if err != nil {
			return nil, err
		}

		return &FunctionSource{
			path:   path,
			origin: functionOrigin,
		}, nil

	case "boolean", "number", "string", "userdata", "thread", "nil":
		variableOrigin, err := c.getOrigin(path, ts.ASSIGNMENT_STATEMENT)

		if err != nil {
			return nil, err
		}

		return &VariableSource{
			path:   path,
			origin: variableOrigin,
		}, nil
	}

	return nil, fmt.Errorf("Unrecognized type '%s' received for path '%s'", luaType, path)
}

func (c *Crawler) getOrigin(path string, fieldType string, fieldTypes ...string) (*origin, error) {
	var orig = origin{}

	locations, err := c.getLocation(path)

	switch {
	case errors.Is(err, nvim.ErrDefinitionLocationNotFound):
		c.config.logger.Warn(fmt.Sprintf("Location not found for '%s' symbol path", path))
		return nil, nil
	case err != nil:
		return nil, err
	}

	locationIndex, err := slicesx.IndexFunc(locations, func(location nvim.Location) (bool, error) {
		_, err = c.config.nvim.Open(location.Url)

		if err != nil {
			return false, err
		}

		position := nvim.CursorPosition{
			Line:      uint(location.TargetRange.Start.Line),
			Character: uint(location.TargetRange.Start.Character),
		}
		definitionLines, err := c.config.nvim.ReadTSNodeAt(position, fieldType, fieldTypes...)

		switch {
		case errors.Is(err, nvim.ErrTsNodeNotFound):
			return false, nil
		case err != nil:
			return false, err
		}

		orig.SetLocation(location)
		orig.SetDefinition(definitionLines)

		return true, nil
	})

	switch {
	case err != nil:
		return nil, err
	case locationIndex == -1:
		c.config.logger.Warn(fmt.Sprintf("Origin not found for '%s' symbol path", path))
		return nil, nil
	}

	functionLocation := locations[locationIndex]

	documentationPosition := nvim.CursorPosition{
		Line:      uint(max(0, functionLocation.TargetRange.Start.Line-1)),
		Character: uint(functionLocation.TargetRange.Start.Character),
	}
	documentationBufferLines, err := c.config.nvim.ReadCommentBlockAt(documentationPosition)

	switch {
	case errors.Is(err, nvim.ErrTsNodeNotFound):
		c.config.logger.Warn(fmt.Sprintf("Documentation not found for '%s' symbol path", path))
		return &orig, nil
	case err != nil:
		return &orig, err
	default:
		orig.SetDocumentation(documentationBufferLines)
	}

	return &orig, nil
}

func (c *Crawler) getLocation(path string) ([]nvim.Location, error) {
	lines := []string{"local ref = " + path}

	err := c.scratch(lines)

	if err != nil {
		return []nvim.Location{}, err
	}

	line, character := uint(0), uint(len(lines[0]))

	locations, err := c.config.nvim.GetDefinitionLocation(line, character)

	if err != nil {
		return []nvim.Location{}, err
	}

	return locations, nil
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
