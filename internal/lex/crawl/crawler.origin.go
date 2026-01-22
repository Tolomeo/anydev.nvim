package crawl

import (
	"fmt"
	"regexp"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (c *Crawler) getOriginQueryMap() nvim.TsNodeQueryMap {
	return nvim.TsNodeQueryMap{
		treesitter.ASSIGNMENT_STATEMENT: nil,
		treesitter.VARIABLE_DECLARATION: nil,
		treesitter.FUNCTION_DECLARATION: nil,
		treesitter.ALIAS_ANNOTATION: {
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(alias_annotation) @alias
				(#match? @alias "\\@alias *%s($|[^a-zA-Z0-9_])")
			)`, regexp.QuoteMeta(c.context.Target().Identifier())),
		},
		treesitter.CLASS_ANNOTATION: {
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(class_annotation) @class_annotation
				(#match? @class_annotation "\\@class *%s($|[^a-zA-Z0-9_])")
			)`, regexp.QuoteMeta(c.context.Target().Identifier())),
		},
		treesitter.FIELD_ANNOTATION: {
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(field_annotation) @field_annotation
				(#match? @field_annotation "\\@field *%s($|[^a-zA-Z0-9_])")
				)`, regexp.QuoteMeta(c.context.Target().Name())),
		},
	}
}

func (c *Crawler) findOrigin(locations []nvim.Location) (*symbol.Origins, error) {
	for _, location := range locations {
		buffer, err := c.context.Nvim().OpenBuffer(location.Url)

		if err != nil {
			return nil, err
		}

		defer buffer.Close()

		queryMap := c.getOriginQueryMap()
		line, character :=
			uint(location.TargetRange.Start.Line),
			uint(location.TargetRange.Start.Character)
		queryMatch, err := buffer.QueryTsNodeAt(queryMap, line, character)

		switch {
		case err != nil:
			return nil, err
		case queryMatch == nil:
			continue
		}

		origin := symbol.NewOrigin(location, queryMatch.Node)
		c.sourceDocumentation(origin)

		targetOrigin := symbol.NewOrigins(origin)

		switch targetOrigin.Last().Type() {
		case treesitter.ASSIGNMENT_STATEMENT:
			followed, err := c.followModuleRequireAssignment(targetOrigin)

			switch {
			case err != nil:
				return nil, err
			case followed:
				return targetOrigin, nil
			}

			followed, err = c.followVariableAssignment(targetOrigin)

			switch {
			case err != nil:
				return nil, err
			case followed:
				return targetOrigin, nil
			}
		}

		return targetOrigin, nil
	}

	return nil, nil
}

var requireAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
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
	)`,
}

func (c *Crawler) followModuleRequireAssignment(origin *symbol.Origins) (bool, error) {
	buffer, err := c.context.Nvim().NewBuffer()

	if err != nil {
		return false, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.Last().Definition())

	if err != nil {
		return false, err
	}

	captures, err := buffer.TsQueryOne(requireAssignmentQuery)

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	moduleNameCapture, found := slicesx.FindFunc(*captures, func(capture treesitter.Capture) bool {
		return capture.Id == "require.module"
	})

	if !found {
		return false, fmt.Errorf("Error retrieving required module name from require statement in '%s'", origin.Last().Definition())
	}

	moduleLocations, err := c.findModuleDefinitionLocations(moduleNameCapture.Node.Text)

	if err != nil {
		return false, err
	}

	moduleOrigin, err := c.findOrigin(*moduleLocations)

	switch {
	case err != nil:
		return false, err
	case moduleOrigin == nil:
		return false, fmt.Errorf("Error following require statement '%s'", origin.Last().Definition())
	}

	origin.Merge(moduleOrigin)

	return true, nil
}

var variableAssignmentQueries = map[string]treesitter.Query{
	"dotIndexAssignment": {
		Language: "lua",
		Query: `
		(assignment_statement
			(variable_list
				name: (_)
			) @assignment.left
			(expression_list
				value: (dot_index_expression) @assignment.right 
			)
		)`},
	"identifierAssignment": {
		Language: "lua",
		Query: `
		(assignment_statement
			(variable_list
				name: (_)
			) @assignment.left
			(expression_list
				value: (identifier) @assignment.right 
			)
		)`,
	},
}

func (c *Crawler) followVariableAssignment(origin *symbol.Origins) (bool, error) {
	buffer, err := c.context.Nvim().NewBuffer()

	if err != nil {
		return false, err
	}

	defer buffer.Close()

	err = buffer.SetLines(origin.Last().Definition())

	if err != nil {
		return false, err
	}

	var captures *nvim.TsQueryMatch = nil

	for _, variableAssignmentQuery := range variableAssignmentQueries {
		assignmentCaptures, err := buffer.TsQueryOne(variableAssignmentQuery)

		switch {
		case err != nil:
			return false, err
		case assignmentCaptures == nil:
			continue
		}

		captures = assignmentCaptures
	}

	if captures == nil {
		return false, nil
	}

	rightValue, found := slicesx.FindFunc(*captures, func(capture treesitter.Capture) bool {
		return capture.Id == "assignment.right"
	})

	if !found {
		return false, fmt.Errorf("Error retrieving read variable name from variable to variable assignment in '%s'", origin.Last().Definition())
	}

	rightValueLocations, err := c.findDefinitionLocations(rightValue.Node.Text)

	if err != nil {
		return false, err
	}

	rightValueOrigin, err := c.findOrigin(*rightValueLocations)

	if err != nil {
		return false, err
	}

	origin.Merge(rightValueOrigin)

	return true, nil
}
