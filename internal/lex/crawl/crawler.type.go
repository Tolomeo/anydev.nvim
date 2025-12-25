package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) sourceType(name string, source *symbol.Source) error {
	typeOrigin, err := c.sourceTypeOrigin(name)

	if err != nil {
		return err
	}

	source.Origin = typeOrigin

	return nil
}

func (c *Crawler) sourceTypeOrigin(path string) (*symbol.Origin, error) {
	locations, err := c.findTypeDefinitionLocations(path)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.context.Logger.Warn(fmt.Sprintf("No locations found for '%s' type", path))
		return nil, nil
	}

	typeOrigin, err := c.findTypeOrigin(path, *locations)

	switch {
	case err != nil:
		return nil, err
	case typeOrigin == nil:
		c.context.Logger.Warn(fmt.Sprintf("No origin found for '%s' symbol", path))
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

func (c *Crawler) sourceTypeOriginDocumentation(name string, typeOrigin *symbol.Origin) (bool, error) {
	_, err := c.context.Nvim.Open(typeOrigin.Url())

	if err != nil {
		return false, err
	}

	commentBlockLines, err := c.readCommentBlock(typeOrigin.Line(), typeOrigin.Character())

	switch {
	case err != nil:
		return false, err
	case commentBlockLines == nil:
		return false, nil
	}

	typeOrigin.Documentation = *commentBlockLines

	return true, nil
}

func (c *Crawler) findTypeOrigin(path string, locations []nvim.Location) (*symbol.Origin, error) {
	for _, location := range locations {
		_, err := c.context.Nvim.Open(location.Url)

		if err != nil {
			return nil, err
		}

		line, character :=
			uint(location.TargetRange.Start.Line),
			uint(location.TargetRange.Start.Character)
		node, err := c.context.Nvim.GetTSNodeAt([]string{treesitter.CLASS_ANNOTATION, treesitter.ALIAS_ANNOTATION}, line, character)

		switch {
		case err != nil:
			return nil, err
		case node == nil:
			continue
		}

		typeOrigin := &symbol.Origin{
			Location: location,
			Node:     *node,
		}

		fmt.Printf("\nType origin:\n%+v\n", location)

		err = c.followValueOrigin(path, &typeOrigin)

		if err != nil {
			return nil, err
		}

		return typeOrigin, nil
	}

	return nil, nil
}

func (c *Crawler) findTypeDefinitionLocations(path string) (*[]nvim.Location, error) {
	buffer, err := c.context.Nvim.Buffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	typeAnnotation := fmt.Sprintf("---@type %s", path)
	err = buffer.SetLines([]string{typeAnnotation})

	line, character := uint(0), uint(len(typeAnnotation))

	locations, err := c.context.Nvim.GetDefinitionLocations(line, character)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		return nil, nil
	}

	return locations, nil
}
