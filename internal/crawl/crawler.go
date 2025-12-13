package crawl

import (
	"errors"
	"fmt"
	"path"
	"regexp"

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

func (c *Crawler) SourceRuntime(path string) (*Source, error) {
	return c.sourceRuntime(path)
}

func (c *Crawler) sourceRuntime(path string) (*Source, error) {
	source := Source{path: path}

	pathOrigin, err := c.sourceOrigin(path)

	if err != nil {
		return nil, err
	}

	source.origin = pathOrigin

	fields, err := c.config.nvim.GetCompletion(path)

	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		child, err := c.sourceRuntime(fmt.Sprintf("%s.%s", path, field))

		if err != nil {
			return nil, err
		}

		source.fields = append(source.fields, child)
	}

	return &source, nil
}

var moduleRequireQuery string = `
	(assignment_statement
		(variable_list)
		(expression_list
			value: (function_call
				name: (identifier) @require.call
				arguments: (arguments
					(string
						content: (string_content) @require.module
					)
				)
			)
		) @require
		(#eq? @require.call "require")
	)
`

func (c *Crawler) resolveOrigin(path string, pathOrigin **origin) error {
	o := *pathOrigin

	switch o.node.Type {
	case ts.ASSIGNMENT_STATEMENT:
		captures, err := c.config.nvim.TsQuery(nvim.TsQueryConfig{Language: "lua", Query: moduleRequireQuery, Range: &ts.LineRange{
			Start: o.node.Range.Start.Line,
			End:   o.node.Range.End.Line,
		}})

		if errors.Is(nvim.ErrTsNodeNotFound, err) {
			return nil
		}

		if err != nil {
			return err
		}

		moduleNameCapture, found := slicesx.FindFunc(captures, func(capture ts.Capture) bool {
			return capture.Id == "require.module"
		})

		if !found {
			return fmt.Errorf("Error retrieving required module name from require statement in '%s'", o.node.Text)
		}

		moduleLocations, err := c.findRequireLocations(moduleNameCapture.Node.Text)

		if err != nil {
			return err
		}

		moduleOrigin, err := c.findOrigin(path, moduleLocations)

		if err != nil {
			return err
		}

		*pathOrigin = moduleOrigin

		return nil
	}

	return nil
}

func (c *Crawler) findOrigin(path string, locations []nvim.Location) (*origin, error) {
	for _, location := range locations {
		_, err := c.config.nvim.Open(location.Url)

		if err != nil {
			return nil, err
		}

		line, character :=
			uint(location.TargetRange.Start.Line),
			uint(location.TargetRange.Start.Character)
		node, err := c.config.nvim.GetTSNodeAt([]string{ts.ASSIGNMENT_STATEMENT, ts.VARIABLE_DECLARATION, ts.FUNCTION_DECLARATION}, line, character)

		if errors.Is(nvim.ErrTsNodeNotFound, err) {
			continue
		}

		if err != nil {
			return nil, err
		}

		pathOrigin := &origin{
			location: location,
			node:     node,
		}

		err = c.resolveOrigin(path, &pathOrigin)

		if err != nil {
			return nil, err
		}

		return pathOrigin, nil
	}

	return nil, nil
}

func (c *Crawler) sourceOriginDocumentation(_ string, pathOrigin *origin) (bool, error) {
	commentBlockLines, err := c.config.nvim.GetCommentBlockAt(pathOrigin.Line()-1, pathOrigin.Character())

	if errors.Is(nvim.ErrCommentBlockNotFound, err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	eCommentContent := regexp.MustCompile(`^[ \t]*-{2,3}(.*)$`)

	documentation := []string{}

	for _, sourceLine := range commentBlockLines {
		matches := eCommentContent.FindStringSubmatch(sourceLine)

		if len(matches) < 2 {
			documentation = append(documentation, "")
			continue
		}

		documentation = append(documentation, matches[1])
	}

	pathOrigin.documentation = documentation

	return true, nil

}

func (c *Crawler) sourceOrigin(path string) (*origin, error) {
	locations, err := c.findDefinitionLocations(path)

	if err != nil {
		return nil, err
	}

	pathOrigin, err := c.findOrigin(path, locations)

	if err != nil {
		return nil, err
	}

	if pathOrigin == nil {
		c.config.logger.Warn(fmt.Sprintf("Origin not found for '%s' symbol path", path))
		return nil, nil
	}

	hasDocumentation, err := c.sourceOriginDocumentation(path, pathOrigin)

	if err != nil {
		return nil, err
	}

	if !hasDocumentation {
		c.config.logger.Warn(fmt.Sprintf("Documentation not found for '%s' symbol path", path))
	}

	return pathOrigin, nil
}

func (c *Crawler) findRequireLocations(moduleName string) ([]nvim.Location, error) {
	lines := []string{fmt.Sprintf("local ref = require('%s')", moduleName)}

	err := c.scratch(lines)

	if err != nil {
		return []nvim.Location{}, err
	}

	line, character := uint(0), uint(len(lines[0])-2)

	locations, err := c.config.nvim.GetDefinitionLocation(line, character)

	if err != nil {
		return []nvim.Location{}, err
	}

	return locations, nil
}

func (c *Crawler) findDefinitionLocations(path string) ([]nvim.Location, error) {
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
