package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	ts "github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (c *Crawler) sourceRuntime(path string, source *Source) error {
	pathOrigin, err := c.sourceOrigin(path)

	if err != nil {
		return err
	}

	source.origin = pathOrigin

	return nil
}

var variableAssignmentQueries = map[string]string{
	"dotIndexAssignment": `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: (dot_index_expression) @assignment.right 
		)
	)`,
	"identifierAssignment": `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: (identifier) @assignment.right 
		)
	)`,
}

func (c *Crawler) followVariableAssignment(path string, o *origin) (*origin, error) {
	buffer, err := c.config.nvim.Buffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	err = buffer.SetLines(o.Definition())

	if err != nil {
		return nil, err
	}

	captures, hasCaptures, err := slicesx.MapFindFunc(mapx.Values(variableAssignmentQueries), func(variableAssignmentQuery string) (*nvim.TsQueryMatch, bool, error) {
		assignmentCaptures, err := c.config.nvim.TsQueryOne(nvim.TsQueryConfig{Language: "lua", Query: variableAssignmentQuery})

		switch {
		case err != nil:
			return nil, false, err
		case assignmentCaptures == nil:
			return nil, false, nil
		}

		return assignmentCaptures, true, nil
	})

	switch {
	case err != nil:
		return nil, err
	case !hasCaptures:
		return nil, nil
	}

	rightValue, found := slicesx.FindFunc(*captures, func(capture ts.Capture) bool {
		return capture.Id == "assignment.right"
	})

	if !found {
		return nil, fmt.Errorf("Error retrieving read variable name from variable to variable assignment in '%s'", o.node.Text)
	}

	rightValueLocations, err := c.findDefinitionLocations(rightValue.Node.Text)

	if err != nil {
		return nil, err
	}

	rightValueOrigin, err := c.findOrigin(path, *rightValueLocations)

	if err != nil {
		return nil, err
	}

	return rightValueOrigin, nil
}

var requireAssignmentQuery string = `
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

func (c *Crawler) followRequireAssignment(path string, o *origin) (*origin, error) {
	buffer, err := c.config.nvim.Buffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	err = buffer.SetLines(o.Definition())

	if err != nil {
		return nil, err
	}

	captures, err := c.config.nvim.TsQueryOne(nvim.TsQueryConfig{Language: "lua", Query: requireAssignmentQuery})

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	moduleNameCapture, found := slicesx.FindFunc(*captures, func(capture ts.Capture) bool {
		return capture.Id == "require.module"
	})

	if !found {
		return nil, fmt.Errorf("Error retrieving required module name from require statement in '%s'", o.node.Text)
	}

	moduleLocations, err := c.findModuleLocations(moduleNameCapture.Node.Text)

	if err != nil {
		return nil, err
	}

	moduleOrigin, err := c.findOrigin(path, *moduleLocations)

	switch {
	case err != nil:
		return nil, err
	case moduleOrigin == nil:
		return nil, fmt.Errorf("Error following require statement '%s'", o.Definition())
	}

	return moduleOrigin, nil
}

func (c *Crawler) follow(path string, pathOrigin **origin) error {
	o := *pathOrigin

	switch o.node.Type {
	case ts.ASSIGNMENT_STATEMENT:
		requiredOrigin, err := c.followRequireAssignment(path, o)

		switch {
		case err != nil:
			return err
		case requiredOrigin != nil:
			*pathOrigin = requiredOrigin
			return nil
		}

		variableOrigin, err := c.followVariableAssignment(path, o)

		switch {
		case err != nil:
			return err
		case variableOrigin != nil:
			*pathOrigin = variableOrigin
			return nil
		}

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

		switch {
		case err != nil:
			return nil, err
		case node == nil:
			continue
		}

		pathOrigin := &origin{
			location: location,
			node:     *node,
		}

		err = c.follow(path, &pathOrigin)

		if err != nil {
			return nil, err
		}

		return pathOrigin, nil
	}

	return nil, nil
}

func (c *Crawler) sourceOriginDocumentation(path string, pathOrigin *origin) (bool, error) {
	_, err := c.config.nvim.Open(pathOrigin.Url())

	if err != nil {
		return false, err
	}

	commentBlockLines, err := c.config.nvim.GetCommentBlockAt(pathOrigin.Line()-1, pathOrigin.Character())

	switch {
	case err != nil:
		return false, err
	case commentBlockLines == nil:
		return false, nil
	}

	pathOrigin.documentation = *commentBlockLines

	return true, nil

}

func (c *Crawler) sourceOrigin(path string) (*origin, error) {
	locations, err := c.findDefinitionLocations(path)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.config.logger.Warn(fmt.Sprintf("No locations found for '%s' symbol", path))
		return nil, nil
	}

	pathOrigin, err := c.findOrigin(path, *locations)

	switch {
	case err != nil:
		return nil, err
	case pathOrigin == nil:
		c.config.logger.Warn(fmt.Sprintf("No origin found for '%s' symbol", path))
		return nil, nil
	}

	hasDocumentation, err := c.sourceOriginDocumentation(path, pathOrigin)

	switch {
	case err != nil:
		return nil, err
	case !hasDocumentation:
		c.config.logger.Warn(fmt.Sprintf("No documentation found for '%s' symbol", path))
	}

	return pathOrigin, nil
}

func (c *Crawler) findModuleLocations(moduleName string) (*[]nvim.Location, error) {
	buffer, err := c.config.nvim.Buffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	lines := []string{fmt.Sprintf("local ref = require('%s')", moduleName)}
	err = buffer.SetLines(lines)

	if err != nil {
		return nil, err
	}

	line, character := uint(0), uint(len(lines[0])-2)

	locations, err := c.config.nvim.GetDefinitionLocation(line, character)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		return nil, nil
	}

	return locations, nil
}

func (c *Crawler) findDefinitionLocations(path string) (*[]nvim.Location, error) {
	buffer, err := c.config.nvim.Buffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	assignment := "local ref = " + path
	err = buffer.SetLines([]string{assignment})

	if err != nil {
		return nil, err
	}

	line, character := uint(0), uint(len(assignment))

	locations, err := c.config.nvim.GetDefinitionLocation(line, character)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		return nil, nil
	}

	return locations, nil
}
