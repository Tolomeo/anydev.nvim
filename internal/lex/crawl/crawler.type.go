package crawl

import (
	"fmt"
	"regexp"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) sourceType(name string, source *symbol.TypeSource) error {
	locations, err := c.getTypeDefinitionLocations(name)

	/* for _, loc := range *locations {
		fmt.Printf("\n\nLocation: %+v\n", loc)
	} */

	switch {
	case err != nil:
		return err
	case locations == nil:
		c.context.Logger.Warn(fmt.Sprintf("No locations found for '%s' type", name))
		return nil
	}

	origin, err := c.findTypeOrigin(name, *locations)

	switch {
	case err != nil:
		return err
	case origin == nil:
		c.context.Logger.Warn(fmt.Sprintf("No origin found for '%s' symbol", name))
		return nil
	}

	source.Origin = origin
	return nil
}

// TODO enum
var customTypeQueries = map[string]func(name string, lineRange *treesitter.LineRange) treesitter.Query{
	"alias": func(aliasName string, lineRange *treesitter.LineRange) treesitter.Query {
		return treesitter.Query{
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(alias_annotation) @alias
				(#match? @alias "\\@alias %s")
			)`, regexp.QuoteMeta(aliasName)),
			Range: lineRange,
		}
	},
	"class": func(aliasName string, lineRange *treesitter.LineRange) treesitter.Query {
		return treesitter.Query{
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(class_annotation) @alias
				(#match? @class "\\@class %s")
			)`, regexp.QuoteMeta(aliasName)),
			Range: lineRange,
		}
	},
}

func (c *Crawler) findTypeOrigin(path string, locations []nvim.Location) (*symbol.TypeOrigin, error) {
	for _, location := range locations {
		_, err := c.context.Nvim.Open(location.Url)

		if err != nil {
			return nil, err
		}

		lineRange := treesitter.LineRange{
			Start: location.TargetRange.Start.Line,
			End:   location.TargetRange.End.Line,
		}

		for _, query := range customTypeQueries {
			match, err := c.context.Nvim.TsQueryOne(query(path, &lineRange))

			switch {
			case err != nil:
				return nil, err
			case match == nil:
				continue
			}

			fmt.Printf("\nFoundLocation: %+v\n\nMatch: %+v\n\nMatchRange: %+v\n\n", location, match, match.Range())

			documentation, err := c.context.Nvim.GetTsCommentBlockAt(uint(match.Range().Start.Line), uint(match.Range().Start.Character))

			switch {
			case err != nil:
				return nil, err
			case documentation == nil:
				continue
			}

			fmt.Printf("\nDocumentation: %+v\n", documentation)

			typeOrigin := &symbol.TypeOrigin{
				Location:      location,
				Documentation: *documentation,
			}

			fmt.Printf("\nType origin:\n%+v\n", location)

			return typeOrigin, nil
		}
	}

	return nil, nil
}

func (c *Crawler) getTypeDefinitionLocations(path string) (*[]nvim.Location, error) {
	buffer, err := c.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Delete()

	typeAnnotation := fmt.Sprintf("---@type %s", path)
	err = buffer.SetLines([]string{typeAnnotation})

	line, character := uint(0), uint(len(typeAnnotation))

	locations, err := c.context.Nvim.GetDefinitionLocations(line, character)

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		return nil, nil
	}

	return locations, nil
}
