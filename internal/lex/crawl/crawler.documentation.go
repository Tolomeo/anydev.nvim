package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) sourceDocumentation(source symbol.Source) (*treesitter.TsNode, error) {
	switch source.GetOrigin().Type() {
	case treesitter.ASSIGNMENT_STATEMENT,
		treesitter.VARIABLE_DECLARATION,
		treesitter.FUNCTION_DECLARATION:
		return c.sourceDefinitionDocumentation(source)
	case treesitter.ALIAS_ANNOTATION,
		treesitter.CLASS_ANNOTATION,
		treesitter.FIELD_ANNOTATION:
		return c.sourceTypeDefinitionDocumentation(source)
	}

	return nil, fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", source.Identifier(), source)
}

func (c *Crawler) sourceDefinitionDocumentation(source symbol.Source) (*treesitter.TsNode, error) {
	buffer, err := c.context.Nvim().OpenBuffer(source.GetOrigin().Url())

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

func (c *Crawler) sourceTypeDefinitionDocumentation(source symbol.Source) (*treesitter.TsNode, error) {
	buffer, err := c.context.Nvim().OpenBuffer(source.GetOrigin().Url())

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
