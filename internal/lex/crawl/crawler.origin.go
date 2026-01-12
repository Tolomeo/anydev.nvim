package crawl

import (
	"fmt"
	"regexp"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

var findOriginQueries = map[string]func(symbol.Source) treesitter.Query{
	treesitter.ASSIGNMENT_STATEMENT: nil,
	treesitter.VARIABLE_DECLARATION: nil,
	treesitter.FUNCTION_DECLARATION: nil,
	treesitter.ALIAS_ANNOTATION: func(source symbol.Source) treesitter.Query {
		return treesitter.Query{
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(alias_annotation) @alias
				(#match? @alias "\\@alias *%s($|[^a-zA-Z0-9_])")
			)`, regexp.QuoteMeta(source.Identifier())),
		}
	},
	treesitter.CLASS_ANNOTATION: func(source symbol.Source) treesitter.Query {
		return treesitter.Query{
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(class_annotation) @class_annotation
				(#match? @class_annotation "\\@class *%s($|[^a-zA-Z0-9_])")
			)`, regexp.QuoteMeta(source.Identifier())),
		}
	},
	treesitter.FIELD_ANNOTATION: func(source symbol.Source) treesitter.Query {
		return treesitter.Query{
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(field_annotation) @field_annotation
				(#match? @field_annotation "\\@field *%s($|[^a-zA-Z0-9_])")
				)`, regexp.QuoteMeta(source.Name())),
		}
	},
}

func (c *Crawler) findOrigin(locations []nvim.Location, source symbol.Source) (symbol.Origin, error) {
	for _, location := range locations {
		buffer, err := c.context.Nvim.OpenBuffer(location.Url)

		if err != nil {
			return nil, err
		}

		defer buffer.Close()

		for nodeType, nodeQuery := range findOriginQueries {
			tsRange := location.TargetRange.AsTreesitter()
			lineRange := tsRange.LineRange()

			targetNode := []string{nodeType}
			line, character :=
				uint(location.TargetRange.Start.Line),
				uint(location.TargetRange.Start.Character)
			node, err := buffer.GetTSNodeAt(targetNode, line, character)

			switch {
			case err != nil:
				return nil, err
			case node == nil:
				continue
			}

			if nodeQuery != nil {
				query := nodeQuery(source)
				fmt.Println(query.Query)
				query.Range = &lineRange
				match, err := buffer.TsQueryOne(query)

				switch {
				case err != nil:
					return nil, err
				case match == nil:
					continue
				}
			}

			return c.followOrigin(source, &symbol.ValueOrigin{
				Location:   location,
				Definition: *node,
			})
		}
	}

	return nil, nil
}

func (c *Crawler) followOrigin(source symbol.Source, origin symbol.Origin) (symbol.Origin, error) {
	switch origin.Type() {
	case treesitter.ASSIGNMENT_STATEMENT:
		requiredOrigin, err := c.followModuleRequireAssignment(source, origin)

		switch {
		case err != nil:
			return nil, err
		case requiredOrigin != nil:
			return requiredOrigin, nil
		}

		variableOrigin, err := c.followVariableAssignment(source, origin)

		switch {
		case err != nil:
			return nil, err
		case variableOrigin != nil:
			return variableOrigin, nil
		}
	}

	return origin, nil
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

func (c *Crawler) followModuleRequireAssignment(source symbol.Source, origin symbol.Origin) (symbol.Origin, error) {
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

	moduleLocations, err := c.findModuleDefinitionLocations(moduleNameCapture.Node.Text)

	if err != nil {
		return nil, err
	}

	moduleOrigin, err := c.findOrigin(*moduleLocations, source)

	switch {
	case err != nil:
		return nil, err
	case moduleOrigin == nil:
		return nil, fmt.Errorf("Error following require statement '%s'", origin.DefinitionLines())
	}

	return moduleOrigin, nil
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

func (c *Crawler) followVariableAssignment(source symbol.Source, origin symbol.Origin) (symbol.Origin, error) {
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

	rightValueLocations, err := c.findDefinitionLocations(rightValue.Node.Text)

	if err != nil {
		return nil, err
	}

	rightValueOrigin, err := c.findOrigin(*rightValueLocations, source)

	if err != nil {
		return nil, err
	}

	return rightValueOrigin, nil
}
