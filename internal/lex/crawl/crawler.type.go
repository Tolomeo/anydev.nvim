package crawl

import (
	"fmt"
	"regexp"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) sourceType(name string, source *symbol.Source) error {
	typeOrigin, err := c.sourceTypeOrigin(name)

	if err != nil {
		return err
	}

	source.Origin = typeOrigin

	return nil
}

func (c *Crawler) sourceTypeOrigin(path string) (*symbol.Origin, error) {
	locations, err := c.findTypeDefinitionLocations(path)

	for _, loc := range *locations {
		fmt.Printf("\n\nLocation: %+v\n", loc)
	}

	switch {
	case err != nil:
		return nil, err
	case locations == nil:
		c.context.Logger.Warn(fmt.Sprintf("No locations found for '%s' type", path))
		return nil, nil
	}

	typeOrigin, err := c.findTypeOrigin(path, *locations)

	switch {
	case err != nil:
		return nil, err
	case typeOrigin == nil:
		c.context.Logger.Warn(fmt.Sprintf("No origin found for '%s' symbol", path))
		return nil, nil
	}

	hasDocumentation, err := c.sourceTypeOriginDocumentation(path, typeOrigin)

	switch {
	case err != nil:
		return nil, err
	case !hasDocumentation:
		return nil, fmt.Errorf("Error reading documentation for annotation type '%s'", path)
	}

	return typeOrigin, nil
}

func (c *Crawler) sourceTypeOriginDocumentation(name string, typeOrigin *symbol.Origin) (bool, error) {
	_, err := c.context.Nvim.Open(typeOrigin.Url())

	if err != nil {
		return false, err
	}

	commentBlockLines, err := c.context.Nvim.GetCommentBlockAt(typeOrigin.Line(), typeOrigin.Character())

	switch {
	case err != nil:
		return false, err
	case commentBlockLines == nil:
		return false, nil
	}

	typeOrigin.Documentation = *commentBlockLines

	return true, nil
}

var findAliasOriginQuery = func(aliasName string, lineRange *treesitter.LineRange) treesitter.Query {
	return treesitter.Query{
		Language: "luadoc",
		Query: fmt.Sprintf(`
			(alias_annotation) @alias
			(#match? @alias "\\@alias %s")
		`, regexp.QuoteMeta(aliasName)),
		Range: lineRange,
	}
}

func (c *Crawler) findTypeOrigin(path string, locations []nvim.Location) (*symbol.Origin, error) {
	for _, location := range locations {
		_, err := c.context.Nvim.Open(location.Url)

		if err != nil {
			return nil, err
		}

		lineRange := treesitter.LineRange{
			Start: location.TargetRange.Start.Line,
			End:   location.TargetRange.End.Line,
		}
		match, err := c.context.Nvim.TsQueryOne(findAliasOriginQuery(path, &lineRange))

		switch {
		case err != nil:
			return nil, err
		case match == nil:
			continue
		}

		fmt.Printf("\nFoundLocation: %+v\n\nMatch: %+v\n\nMatchRange: %+v\n\n", location, match, match.Range())

		block, _ := c.context.Nvim.GetTsCommentBlockAt(uint(match.Range().Start.Line), uint(match.Range().Start.Character))

		fmt.Printf("\nCommentBlock: %+v\n\n", block)

		documentation, err := c.context.Nvim.GetCommentBlockAt(uint(match.Range().Start.Line), uint(match.Range().Start.Character))

		fmt.Printf("\nDocumentation: %+v\n", documentation)

		typeOrigin := &symbol.Origin{
			Location:      location,
			Documentation: *documentation,
		}

		fmt.Printf("\nType origin:\n%+v\n", location)

		err = c.followValueOrigin(path, &typeOrigin)

		if err != nil {
			return nil, err
		}

		return typeOrigin, nil
	}

	return nil, nil
}

func (c *Crawler) findTypeDefinitionLocations(path string) (*[]nvim.Location, error) {
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
