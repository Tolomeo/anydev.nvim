package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) sourceDocumentation(origin *symbol.Origin) (*treesitter.TsNode, error) {
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

	return nil, fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", c.target.Identifier(), origin)
}

func (c *Crawler) sourceDefinitionDocumentation(origin *symbol.Origin) (*treesitter.TsNode, error) {
	buffer, err := c.target.Nvim().OpenBuffer(origin.Url())

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	documentationBlock, err := buffer.GetTsCommentBlockAt(origin.Line()-1, origin.Character())

	switch {
	case err != nil:
		return nil, err
	case documentationBlock == nil:
		return nil, nil
	}

	return documentationBlock, nil
}

func (c *Crawler) sourceTypeDefinitionDocumentation(origin *symbol.Origin) (*treesitter.TsNode, error) {
	buffer, err := c.target.Nvim().OpenBuffer(origin.Url())

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	documentationBlock, err := buffer.GetTsCommentBlockAt(origin.Line(), origin.Character())

	switch {
	case err != nil:
		return nil, err
	case documentationBlock == nil:
		return nil, nil
	}

	return documentationBlock, nil
}
