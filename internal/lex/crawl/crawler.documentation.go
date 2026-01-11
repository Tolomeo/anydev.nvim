package crawl

import (
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) sourceDefinitionDocumentation(source symbol.Source) (*treesitter.TsNode, error) {
	buffer, err := c.context.Nvim.OpenBuffer(source.GetOrigin().Url())

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
