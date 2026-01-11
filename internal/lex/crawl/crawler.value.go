package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (c *Crawler) followValueVariableAssignment(source symbol.Source, origin *symbol.ValueOrigin) (*symbol.ValueOrigin, error) {
	buffer, err := c.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.DefinitionLines())

	if err != nil {
		return nil, err
	}

	var captures *nvim.TsQueryMatch = nil

	for _, variableAssignmentQuery := range variableAssignmentQueries {
		assignmentCaptures, err := buffer.TsQueryOne(variableAssignmentQuery)

		switch {
		case err != nil:
			return nil, err
		case assignmentCaptures == nil:
			continue
		}

		captures = assignmentCaptures
	}

	if captures == nil {
		return nil, nil
	}

	rightValue, found := slicesx.FindFunc(*captures, func(capture treesitter.Capture) bool {
		return capture.Id == "assignment.right"
	})

	if !found {
		return nil, fmt.Errorf("Error retrieving read variable name from variable to variable assignment in '%s'", origin.DefinitionLines())
	}

	rightValueLocations, err := c.findValueDefinitionLocations(rightValue.Node.Text)

	if err != nil {
		return nil, err
	}

	rightValueOrigin, err := c.findValueOrigin(source, *rightValueLocations)

	if err != nil {
		return nil, err
	}

	return rightValueOrigin, nil
}

func (c *Crawler) followRequireValueAssignment(source symbol.Source, origin *symbol.ValueOrigin) (*symbol.ValueOrigin, error) {
	buffer, err := c.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.DefinitionLines())

	if err != nil {
		return nil, err
	}

	captures, err := buffer.TsQueryOne(requireAssignmentQuery)

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	moduleNameCapture, found := slicesx.FindFunc(*captures, func(capture treesitter.Capture) bool {
		return capture.Id == "require.module"
	})

	if !found {
		return nil, fmt.Errorf("Error retrieving required module name from require statement in '%s'", origin.DefinitionLines())
	}

	moduleLocations, err := c.findModuleValueLocations(moduleNameCapture.Node.Text)

	if err != nil {
		return nil, err
	}

	moduleOrigin, err := c.findValueOrigin(source, *moduleLocations)

	switch {
	case err != nil:
		return nil, err
	case moduleOrigin == nil:
		return nil, fmt.Errorf("Error following require statement '%s'", origin.DefinitionLines())
	}

	return moduleOrigin, nil
}

func (c *Crawler) followValueOrigin(source symbol.Source, origin *symbol.ValueOrigin) (*symbol.ValueOrigin, error) {
	switch origin.Type() {
	case treesitter.ASSIGNMENT_STATEMENT:
		requiredOrigin, err := c.followRequireValueAssignment(source, origin)

		switch {
		case err != nil:
			return nil, err
		case requiredOrigin != nil:
			return requiredOrigin, nil
		}

		variableOrigin, err := c.followValueVariableAssignment(source, origin)

		switch {
		case err != nil:
			return nil, err
		case variableOrigin != nil:
			return variableOrigin, nil
		}
	}

	return origin, nil
}

var valueOriginQueries = map[string]func(name string) treesitter.Query{
	treesitter.ASSIGNMENT_STATEMENT: nil,
	treesitter.VARIABLE_DECLARATION: nil,
	treesitter.FUNCTION_DECLARATION: nil,
}

func (c *Crawler) findValueOrigin(source symbol.Source, locations []nvim.Location) (*symbol.ValueOrigin, error) {
	for _, location := range locations {
		buffer, err := c.context.Nvim.OpenBuffer(location.Url)

		if err != nil {
			return nil, err
		}

		defer buffer.Close()

		for searchNode, searchQuery := range valueOriginQueries {
			tsRange := location.TargetRange.AsTreesitter()
			lineRange := tsRange.LineRange()

			targetNodes := []string{searchNode}
			line, character :=
				uint(location.TargetRange.Start.Line),
				uint(location.TargetRange.Start.Character)
			node, err := buffer.GetTSNodeAt(targetNodes, line, character)

			switch {
			case err != nil:
				return nil, err
			case node == nil:
				continue
			}

			if searchQuery != nil {
				query := (searchQuery)(source.Identifier())
				query.Range = &lineRange
				match, err := buffer.TsQueryOne(query)

				switch {
				case err != nil:
					return nil, err
				case match == nil:
					continue
				}
			}

			return c.followValueOrigin(source, &symbol.ValueOrigin{
				Location:   location,
				Definition: *node,
			})
		}
	}

	return nil, nil
}

func (c *Crawler) sourceValueOriginDocumentation(source symbol.Source) (*treesitter.TsNode, error) {
	buffer, err := c.context.Nvim.OpenBuffer(source.GetOrigin().Url())

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	documentationBlock, err := buffer.GetTsCommentBlockAt(source.GetOrigin().Line()-1, source.GetOrigin().Character())

	switch {
	case err != nil:
		return nil, err
	case documentationBlock == nil:
		return nil, nil
	}

	return documentationBlock, nil
}

func (c *Crawler) findModuleValueLocations(moduleName string) (*[]nvim.Location, error) {
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

func (c *Crawler) findValueDefinitionLocations(path string) (*[]nvim.Location, error) {
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
