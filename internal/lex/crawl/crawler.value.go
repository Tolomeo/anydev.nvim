package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

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


func (c *Crawler) followVariableAssignment(source symbol.Source) (*symbol.ValueOrigin, error) {
	buffer, err := c.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(source.GetOrigin().DefinitionLines())

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
		return nil, fmt.Errorf("Error retrieving read variable name from variable to variable assignment in '%s'", source.GetOrigin().DefinitionLines())
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

func (c *Crawler) followRequireValueAssignment(source symbol.Source) (*symbol.ValueOrigin, error) {
	buffer, err := c.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(source.GetOrigin().DefinitionLines())

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
		return nil, fmt.Errorf("Error retrieving required module name from require statement in '%s'", source.GetOrigin().DefinitionLines())
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
		return nil, fmt.Errorf("Error following require statement '%s'", source.GetOrigin().DefinitionLines())
	}

	return moduleOrigin, nil
}

func (c *Crawler) followValueOrigin(source symbol.Source) error {
	switch source.GetOrigin().Type() {
	case treesitter.ASSIGNMENT_STATEMENT:
		requiredOrigin, err := c.followRequireValueAssignment(source)

		switch {
		case err != nil:
			return err
		case requiredOrigin != nil:
			source.SetOrigin(requiredOrigin)
			return nil
		}

		variableOrigin, err := c.followVariableAssignment(source)

		switch {
		case err != nil:
			return err
		case variableOrigin != nil:
			source.SetOrigin(variableOrigin)
			return nil
		}

		return nil
	}

	return nil
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

			if searchQuery == nil {
				return &symbol.ValueOrigin{
					Location:   location,
					Definition: *node,
				}, nil
			}

			query := (searchQuery)(source.Identifier())
			query.Range = &lineRange
			match, err := buffer.TsQueryOne(query)

			switch {
			case err != nil:
				return nil, err
			case match == nil:
				continue
			}

			return &symbol.ValueOrigin{
				Location:   location,
				Definition: *node,
			}, nil
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
