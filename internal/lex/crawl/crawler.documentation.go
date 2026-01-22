package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) getCommentBlock(definition treesitter.TsNode, location nvim.Location) (*treesitter.TsNode, error) {
	buffer, err := c.context.Nvim().OpenBuffer(location.Url)

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	var documentationBlock *treesitter.TsNode

	switch definition.Type {
	case treesitter.ASSIGNMENT_STATEMENT, treesitter.VARIABLE_DECLARATION, treesitter.FUNCTION_DECLARATION:
		documentationBlock, err = buffer.GetTsCommentBlockAt(location.StartLine()-1, location.StartCharacter())
	case treesitter.ALIAS_ANNOTATION, treesitter.CLASS_ANNOTATION, treesitter.FIELD_ANNOTATION:
		documentationBlock, err = buffer.GetTsCommentBlockAt(location.StartLine(), location.StartCharacter())
	default:
		return nil, fmt.Errorf("Unknown origin type received for source '%s' with value <%+v>", c.context.Target().Identifier(), definition)
	}

	if err != nil {
		return nil, err
	}

	return documentationBlock, nil
}
