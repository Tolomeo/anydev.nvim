package crawl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) sourceType(source *symbol.TypeSource) error {
	locations, err := c.getTypeDefinitionLocations(source.Path)

	/* for _, loc := range *locations {
		fmt.Printf("\n\nLocation: %+v\n", loc)
	} */

	switch {
	case err != nil:
		return err
	case locations == nil:
		c.context.Logger.Warn(fmt.Sprintf("No locations found for '%s' type", source.Path))
		return nil
	}

	origin, err := c.findTypeOrigin(source, *locations)

	switch {
	case err != nil:
		return err
	case origin == nil:
		c.context.Logger.Warn(fmt.Sprintf("No origin found for '%s' symbol", source.Path))
		return nil
	}

	source.SetOrigin(origin)
	return nil
}

// TODO enum
var customTypeQueries = map[string]func(name string, lineRange *treesitter.LineRange) treesitter.Query{
	treesitter.ALIAS_ANNOTATION: func(aliasName string, lineRange *treesitter.LineRange) treesitter.Query {
		return treesitter.Query{
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(alias_annotation) @alias
				(#match? @alias "\\@alias %s")
			)`, regexp.QuoteMeta(aliasName)),
			Range: lineRange,
		}
	},
	treesitter.CLASS_ANNOTATION: func(aliasName string, lineRange *treesitter.LineRange) treesitter.Query {
		return treesitter.Query{
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(class_annotation) @class_annotation
				(#match? @class_annotation "\\@class %s")
			)`, regexp.QuoteMeta(aliasName)),
			Range: lineRange,
		}
	},
}

func (c *Crawler) findTypeOrigin(source *symbol.TypeSource, locations []nvim.Location) (*symbol.TypeOrigin, error) {
	for _, location := range locations {
		buffer, err := c.context.Nvim.OpenBuffer(location.Url)

		if err != nil {
			return nil, err
		}

		defer buffer.Close()

		for nodeType, matchNameQuery := range customTypeQueries {
			tsRange := location.TargetRange.AsTreesitter()
			lineRange := tsRange.LineRange()

			definition, err := buffer.GetTSNodeAt([]string{nodeType}, uint(location.TargetRange.Start.Line), uint(location.TargetRange.Start.Character))

			switch {
			case err != nil:
				return nil, err
			case definition == nil:
				continue
			}

			match, err := buffer.TsQueryOne(matchNameQuery(source.Path, &lineRange))

			switch {
			case err != nil:
				return nil, err
			case match == nil:
				continue
			}

			// fmt.Printf("\nFoundLocation: %+v\n\nMatch: %+v\n\nMatchRange: %+v\n\n", location, match, match.Range())

			documentation, err := buffer.GetTsCommentBlockAt(uint(match.Range().Start.Line), uint(match.Range().Start.Character))

			switch {
			case err != nil:
				return nil, err
			case documentation == nil:
				continue
			}

			// fmt.Printf("\nDocumentation: %+v\n", documentation)

			typeOrigin := &symbol.TypeOrigin{
				Location:      location,
				Definition:    *definition,
				Documentation: *documentation,
			}

			// fmt.Printf("\nType origin:\n%+v\n", location)

			return typeOrigin, nil
		}
	}

	return nil, nil
}

func (c *Crawler) getTypeDefinitionLocations(typeName string, typeField ...string) (*[]nvim.Location, error) {
	buffer, err := c.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	typeAnnotation := fmt.Sprintf("---@type %s", typeName)
	ref := "local ref"
	refAccess := strings.Join(append([]string{"ref"}, typeField...), ".")
	err = buffer.SetLines([]string{typeAnnotation, ref, refAccess})

	if err != nil {
		return nil, err
	}

	line, character := uint(2), uint(len(refAccess))
	locations, err := buffer.GetTypeDefinitionLocations(line, character)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		return nil, nil
	}

	return locations, nil
}
