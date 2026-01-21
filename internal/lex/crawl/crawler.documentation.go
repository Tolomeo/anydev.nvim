package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) sourceDocumentation(origin *locationOrigin) error {
	switch origin.Definition.Type {
	case treesitter.ASSIGNMENT_STATEMENT,
		treesitter.VARIABLE_DECLARATION,
		treesitter.FUNCTION_DECLARATION:
		return c.sourceDefinitionDocumentation(origin)
	case treesitter.ALIAS_ANNOTATION,
		treesitter.CLASS_ANNOTATION,
		treesitter.FIELD_ANNOTATION:
		return c.sourceTypeDefinitionDocumentation(origin)
	}

	return fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", c.target.Identifier(), origin)
}

func (c *Crawler) sourceDefinitionDocumentation(origin *locationOrigin) error {
	buffer, err := c.target.Nvim().OpenBuffer(origin.Location.Url)

	if err != nil {
		return err
	}

	defer buffer.Close()

	line, character := uint(origin.Definition.Range.Start.Line)-1, uint(origin.Definition.Range.Start.Character)
	documentationBlock, err := buffer.GetTsCommentBlockAt(line, character)

	switch {
	case err != nil:
		return err
	case documentationBlock == nil:
		c.target.Logger().Warnf("No documentation found for origin <%v>", origin)
		return nil
	}

	origin.Documentation = *documentationBlock

	return nil
}

func (c *Crawler) sourceTypeDefinitionDocumentation(origin *locationOrigin) error {
	buffer, err := c.target.Nvim().OpenBuffer(origin.Location.Url)

	if err != nil {
		return err
	}

	defer buffer.Close()

	line, character := uint(origin.Definition.Range.Start.Line), uint(origin.Definition.Range.Start.Character)
	documentationBlock, err := buffer.GetTsCommentBlockAt(line, character)

	switch {
	case err != nil:
		return err
	case documentationBlock == nil:
		c.target.Logger().Warnf("No documentation found for origin <%v>", origin)
		return nil
	}

	origin.Documentation = *documentationBlock

	return nil
}
