package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) sourceDocumentation(origin *symbol.Origin) error {
	switch origin.Type() {
	case treesitter.ASSIGNMENT_STATEMENT,
		treesitter.VARIABLE_DECLARATION,
		treesitter.FUNCTION_DECLARATION:
		return c.sourceDefinitionDocumentation(origin)
	case treesitter.ALIAS_ANNOTATION,
		treesitter.CLASS_ANNOTATION,
		treesitter.FIELD_ANNOTATION:
		return c.sourceTypeDefinitionDocumentation(origin)
	}

	return fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", c.context.Target().Identifier(), origin)
}

func (c *Crawler) sourceDefinitionDocumentation(origin *symbol.Origin) error {
	buffer, err := c.context.Nvim().OpenBuffer(origin.Url())

	if err != nil {
		return err
	}

	defer buffer.Close()

	line, character := origin.Line()-1, origin.Character()
	documentationBlock, err := buffer.GetTsCommentBlockAt(line, character)

	switch {
	case err != nil:
		return err
	case documentationBlock == nil:
		c.context.Logger().Warnf("No documentation found for origin <%v>", origin)
		return nil
	}

	origin.SetDocumentation(*documentationBlock)

	return nil
}

func (c *Crawler) sourceTypeDefinitionDocumentation(origin *symbol.Origin) error {
	buffer, err := c.context.Nvim().OpenBuffer(origin.Url())

	if err != nil {
		return err
	}

	defer buffer.Close()

	line, character := origin.Line(), origin.Character()
	documentationBlock, err := buffer.GetTsCommentBlockAt(line, character)

	switch {
	case err != nil:
		return err
	case documentationBlock == nil:
		c.context.Logger().Warnf("No documentation found for origin <%v>", origin)
		return nil
	}

	origin.SetDocumentation(*documentationBlock)

	return nil
}
