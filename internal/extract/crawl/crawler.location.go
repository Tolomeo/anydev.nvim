package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

func (c *Crawler) findModuleDefinitionLocations(moduleUrl string) (*[]nvim.Location, error) {
	moduleBuffer, err := c.context.Nvim().OpenBuffer(moduleUrl)

	if err != nil {
		return nil, err
	}

	defer moduleBuffer.Close()

	match, err := moduleBuffer.TsQueryOne(definition.ModuleExportQuery)

	if err != nil {
		return nil, err
	}

	moduleExportDefinition := definition.NewModuleExport(*match)
	line, character := uint(moduleExportDefinition.Range().End.Line), uint(moduleExportDefinition.Range().End.Character)
	locations, err := moduleBuffer.GetDefinitionLocations(line, character)

	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (c *Crawler) findIdentifierDefinitionLocations(identifier string) (*[]nvim.Location, error) {
	buffer, err := c.context.Nvim().NewBuffer()

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

	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (c *Crawler) findTypeIdentifierDefinitionLocations(typeName string, parentTypeName string) (*[]nvim.Location, error) {
	buffer, err := c.context.Nvim().NewBuffer()

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

	if err != nil {
		return nil, err
	}

	if locations == nil {
		return nil, nil
	}

	return locations, nil
}

func (c *Crawler) findDefinitionLocationsAt(url string, line uint, character uint) (*[]nvim.Location, error) {
	buffer, err := c.context.Nvim().OpenBuffer(url)

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	locations, err := buffer.GetDefinitionLocations(line, character)

	if err != nil {
		return nil, err
	}

	if locations == nil {
		return nil, nil
	}

	return locations, nil
}
