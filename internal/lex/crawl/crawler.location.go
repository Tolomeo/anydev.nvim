package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

func (c *Crawler) findModuleLocations(moduleName string) (*[]nvim.Location, error) {
	buffer, err := c.context.Nvim.NewBuffer()

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

func (c *Crawler) findDefinitionLocations(path string) (*[]nvim.Location, error) {
	buffer, err := c.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	assignment := "local ref = " + path
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
