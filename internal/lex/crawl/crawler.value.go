package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

/* func (c *Crawler) sourceValue(path string, source *symbol.ValueSource) error {
	origin, err := c.sourceValueOrigin(path)

	if err != nil {
		return err
	}

	source.Origin = origin

	return nil
} */

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

func (c *Crawler) followVariableAssignment(source *symbol.ValueSource) (*symbol.ValueOrigin, error) {
	buffer, err := c.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(source.DefinitionLines())

	if err != nil {
		return nil, err
	}

	captures, hasCaptures, err := slicesx.MapFindFunc(mapx.Values(variableAssignmentQueries), func(variableAssignmentQuery string) (*nvim.TsQueryMatch, bool, error) {
		assignmentCaptures, err := buffer.TsQueryOne(treesitter.Query{Language: "lua", Query: variableAssignmentQuery})

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

	rightValue, found := slicesx.FindFunc(*captures, func(capture treesitter.Capture) bool {
		return capture.Id == "assignment.right"
	})

	if !found {
		return nil, fmt.Errorf("Error retrieving read variable name from variable to variable assignment in '%s'", source.DefinitionLines())
	}

	rightValueLocations, err := c.findValueDefinitionLocations(rightValue.Node.Text)

	if err != nil {
		return nil, err
	}

	rightValueOrigin, err := c.findValueOrigin(*rightValueLocations)

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

func (c *Crawler) followRequireValueAssignment(source *symbol.ValueSource) (*symbol.ValueOrigin, error) {
	buffer, err := c.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(source.DefinitionLines())

	if err != nil {
		return nil, err
	}

	captures, err := buffer.TsQueryOne(treesitter.Query{Language: "lua", Query: requireAssignmentQuery})

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
		return nil, fmt.Errorf("Error retrieving required module name from require statement in '%s'", source.DefinitionLines())
	}

	moduleLocations, err := c.findModuleValueLocations(moduleNameCapture.Node.Text)

	if err != nil {
		return nil, err
	}

	moduleOrigin, err := c.findValueOrigin(*moduleLocations)

	switch {
	case err != nil:
		return nil, err
	case moduleOrigin == nil:
		return nil, fmt.Errorf("Error following require statement '%s'", source.DefinitionLines())
	}

	return moduleOrigin, nil
}

func (c *Crawler) followValueOrigin(source *symbol.ValueSource) error {
	switch source.Type() {
	case treesitter.ASSIGNMENT_STATEMENT:
		requiredOrigin, err := c.followRequireValueAssignment(source)

		switch {
		case err != nil:
			return err
		case requiredOrigin != nil:
			source.Origin = requiredOrigin
			return nil
		}

		variableOrigin, err := c.followVariableAssignment(source)

		switch {
		case err != nil:
			return err
		case variableOrigin != nil:
			source.Origin = variableOrigin
			return nil
		}

		return nil
	}

	return nil
}

func (c *Crawler) findValueOrigin(locations []nvim.Location) (*symbol.ValueOrigin, error) {
	for _, location := range locations {
		buffer, err := c.context.Nvim.OpenBuffer(location.Url)

		if err != nil {
			return nil, err
		}

		defer buffer.Close()

		line, character :=
			uint(location.TargetRange.Start.Line),
			uint(location.TargetRange.Start.Character)
		node, err := buffer.GetTSNodeAt([]string{treesitter.ASSIGNMENT_STATEMENT, treesitter.VARIABLE_DECLARATION, treesitter.FUNCTION_DECLARATION}, line, character)

		switch {
		case err != nil:
			return nil, err
		case node == nil:
			continue
		}

		pathOrigin := &symbol.ValueOrigin{
			Location:   location,
			Definition: *node,
		}

		return pathOrigin, nil
	}

	return nil, nil
}

func (c *Crawler) sourceValueOrigin(source *symbol.ValueSource) error {
	locations, err := c.findValueDefinitionLocations(source.Path)

	switch {
	case err != nil:
		return err
	case locations == nil:
		c.context.Logger.Warn(fmt.Sprintf("No locations found for '%s' symbol", source.Path))
		return nil
	}

	pathOrigin, err := c.findValueOrigin(*locations)

	switch {
	case err != nil:
		return err
	case pathOrigin == nil:
		c.context.Logger.Warn(fmt.Sprintf("No origin found for '%s' symbol", source.Path))
		return nil
	default:
		source.Origin = pathOrigin
	}

	err = c.followValueOrigin(source)

	if err != nil {
		return err
	}

	documentation, err := c.sourceValueOriginDocumentation(source)

	switch {
	case err != nil:
		return err
	case documentation == nil:
		c.context.Logger.Warn(fmt.Sprintf("No documentation found for '%s' symbol", source.Path))
		return nil
	}

	pathOrigin.Documentation = documentation
	return nil
}

func (c *Crawler) sourceValueOriginDocumentation(source *symbol.ValueSource) (*treesitter.TsNode, error) {
	buffer, err := c.context.Nvim.OpenBuffer(source.Url())

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	documentationBlock, err := buffer.GetTsCommentBlockAt(source.Line()-1, source.Character())

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
