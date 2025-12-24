package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	ts "github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) sourceType(name string, source *Source) error {
	typeOrigin, err := c.sourceTypeOrigin(name)

	if err != nil {
		return err
	}

	source.origin = typeOrigin

	return nil
}

func (c *Crawler) sourceTypeOrigin(path string) (*origin, error) {
	locations, err := c.findTypeDefinitionLocations(path)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.config.logger.Warn(fmt.Sprintf("No locations found for '%s' type", path))
		return nil, nil
	}

	typeOrigin, err := c.findTypeOrigin(path, *locations)

	switch {
	case err != nil:
		return nil, err
	case typeOrigin == nil:
		c.config.logger.Warn(fmt.Sprintf("No origin found for '%s' symbol", path))
		return nil, nil
	}

	hasDocumentation, err := c.sourceTypeOriginDocumentation(path, typeOrigin)

	switch {
	case err != nil:
		return nil, err
	case !hasDocumentation:
		return nil, fmt.Errorf("Error reading documentation for annotation type '%s'", path)
	}

	return typeOrigin, nil
}

func (c *Crawler) sourceTypeOriginDocumentation(name string, typeOrigin *origin) (bool, error) {
	_, err := c.config.nvim.Open(typeOrigin.Url())

	if err != nil {
		return false, err
	}

	commentBlockLines, err := c.config.nvim.GetCommentBlockAt(typeOrigin.Line(), typeOrigin.Character())

	switch {
	case err != nil:
		return false, err
	case commentBlockLines == nil:
		return false, nil
	}

	typeOrigin.documentation = *commentBlockLines

	return true, nil
}

func (c *Crawler) findTypeOrigin(path string, locations []nvim.Location) (*origin, error) {
	for _, location := range locations {
		_, err := c.config.nvim.Open(location.Url)

		if err != nil {
			return nil, err
		}

		line, character :=
			uint(location.TargetRange.Start.Line),
			uint(location.TargetRange.Start.Character)
		node, err := c.config.nvim.GetTSNodeAt([]string{ts.CLASS_ANNOTATION, ts.ALIAS_ANNOTATION}, line, character)

		// err = c.config.nvim.GetTSNodeAtTest([]string{ts.CLASS_ANNOTATION, ts.ALIAS_ANNOTATION}, line, character)

		// fmt.Printf("\nError:\n%+v\n", err)

		switch {
		case err != nil:
			return nil, err
		case node == nil:
			continue
		}

		typeOrigin := &origin{
			location: location,
			node:     *node,
		}

		fmt.Printf("\nType origin:\n%+v\n", location)

		err = c.follow(path, &typeOrigin)

		if err != nil {
			return nil, err
		}

		return typeOrigin, nil
	}

	return nil, nil
}

func (c *Crawler) findTypeDefinitionLocations(path string) (*[]nvim.Location, error) {
	buffer, err := c.config.nvim.Buffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	typeAnnotation := fmt.Sprintf("---@type %s", path)
	err = buffer.SetLines([]string{typeAnnotation})

	line, character := uint(0), uint(len(typeAnnotation))

	locations, err := c.config.nvim.GetDefinitionLocation(line, character)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		return nil, nil
	}

	return locations, nil
}
