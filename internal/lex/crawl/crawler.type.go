package crawl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

// TODO enum
var typeOriginQueries = map[string]func(string) treesitter.Query{
	treesitter.ALIAS_ANNOTATION: func(name string) treesitter.Query {
		return treesitter.Query{
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(alias_annotation) @alias
				(#match? @alias "\\@alias *%s($|[^a-zA-Z0-9_])")
			)`, regexp.QuoteMeta(name)),
		}
	},
	treesitter.CLASS_ANNOTATION: func(name string) treesitter.Query {
		return treesitter.Query{
			Language: "luadoc",
			Query: fmt.Sprintf(`(
				(class_annotation) @class_annotation
				(#match? @class_annotation "\\@class *%s($|[^a-zA-Z0-9_])")
			)`, regexp.QuoteMeta(name)),
		}
	},
}

func (c *Crawler) findTypeOrigin(source symbol.Source, locations []nvim.Location) (*symbol.TypeOrigin, error) {
	for _, location := range locations {
		buffer, err := c.context.Nvim.OpenBuffer(location.Url)

		if err != nil {
			return nil, err
		}

		defer buffer.Close()

		for searchNode, searchQuery := range typeOriginQueries {
			tsRange := location.TargetRange.AsTreesitter()
			lineRange := tsRange.LineRange()

			targetNodes := []string{searchNode}
			line, character :=
				uint(location.TargetRange.Start.Line),
				uint(location.TargetRange.Start.Character)
			definition, err := buffer.GetTSNodeAt(targetNodes, line, character)

			switch {
			case err != nil:
				return nil, err
			case definition == nil:
				continue
			}

			query := searchQuery(source.Identifier())
			query.Range = &lineRange
			match, err := buffer.TsQueryOne(query)

			switch {
			case err != nil:
				return nil, err
			case match == nil:
				continue
			}

			return &symbol.TypeOrigin{
				Location:   location,
				Definition: *definition,
				// Documentation: *documentation,
			}, nil
		}
	}

	return nil, nil
}

func (c *Crawler) sourceTypeOriginDocumentation(source symbol.Source) (*treesitter.TsNode, error) {
	buffer, err := c.context.Nvim.OpenBuffer(source.GetOrigin().Url())

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	documentationBlock, err := buffer.GetTsCommentBlockAt(source.GetOrigin().Line(), source.GetOrigin().Character())

	switch {
	case err != nil:
		return nil, err
	case documentationBlock == nil:
		return nil, nil
	}

	return documentationBlock, nil
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
