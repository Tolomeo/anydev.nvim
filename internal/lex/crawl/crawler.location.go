package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

func (c *Crawler) findModuleDefinitionLocations(moduleName string) (*[]nvim.Location, error) {
	buffer, err := c.target.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	lines := []string{fmt.Sprintf("local ref = require('%s')", moduleName)}
	err = buffer.SetLines(lines)

	if err != nil {
		return nil, err
	}

	line, character := uint(0), uint(len(lines[0])-2)

	locations, err := buffer.GetDefinitionLocations(line, character)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		return nil, nil
	}

	return locations, nil
}

func (c *Crawler) findDefinitionLocations(identifier string) (*[]nvim.Location, error) {
	buffer, err := c.target.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	assignment := "local ref = " + identifier
	err = buffer.SetLines([]string{assignment})

	if err != nil {
		return nil, err
	}

	line, character := uint(0), uint(len(assignment))

	locations, err := buffer.GetDefinitionLocations(line, character)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		return nil, nil
	}

	return locations, nil
}

func (c *Crawler) findTypeDefinitionLocations(typeName string, parentTypeName string) (*[]nvim.Location, error) {
	buffer, err := c.target.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	lines := []string{}

	if parentTypeName == "" {
		lines = append(lines, fmt.Sprintf("---@type %s", typeName))
	} else {
		lines = append(lines, fmt.Sprintf("---@type %s", parentTypeName), "local ref", fmt.Sprintf("ref.%s", typeName))
	}

	err = buffer.SetLines(lines)

	if err != nil {
		return nil, err
	}

	lastLineIndex := len(lines) - 1
	lastLine := lines[lastLineIndex]
	line, character := uint(lastLineIndex), uint(len(lastLine))
	locations, err := buffer.GetDefinitionLocations(line, character)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		return nil, nil
	}

	return locations, nil
}
