package crawl

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/languageserver"
)

func (c *Crawler) GetDocumentation() (symbol.Documentation, error) {

	var markupContent *languageserver.MarkupContent
	var err error

	switch c.context.Target().Kind() {
	case target.TargetKindValue:
		markupContent, err = c.getDefinitionDocumentation(c.context.Target().Identifier())
	case target.TargetKindType:
		markupContent, err = c.getTypeDefinitionDocumentation(c.context.Target().Name(), c.context.Target().ParentName())
	}

	if err != nil {
		return symbol.Documentation{}, err
	}

	if markupContent == nil {
		return symbol.Documentation{}, nil
	}

	return strings.Split(markupContent.Value, "\n"), nil
}

func (c *Crawler) getDefinitionDocumentation(identifier string) (*languageserver.MarkupContent, error) {
	buffer, err := c.context.Nvim().OpenScratchBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	lines := []string{}
	line, character := uint(0), uint(len(identifier))-1
	identifierParts := strings.Split(identifier, ".")

	if len(identifierParts) > 1 {
		tail := identifierParts[len(identifierParts)-1]
		_, isKeyword := definition.Keywords[tail]

		if isKeyword {
			reference := fmt.Sprintf("%s_", identifier)
			character = uint(len(reference)) - 1
			head := strings.Join(identifierParts[:len(identifierParts)-1], ".")
			lines = append(lines, fmt.Sprintf("%s = %s['%s']", reference, head, tail))
		} else {
			lines = append(lines, fmt.Sprintf("%s", identifier))
		}
	}

	err = buffer.SetLines(lines)

	if err != nil {
		return nil, err
	}

	hover, err := buffer.GetHover(line, character)

	if err != nil {
		return nil, err
	}

	return hover, nil
}

func (c *Crawler) getTypeDefinitionDocumentation(typeName string, parentTypeName string) (*languageserver.MarkupContent, error) {
	buffer, err := c.context.Nvim().OpenScratchBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	lines := []string{}

	if parentTypeName == "" {
		lines = append(lines, fmt.Sprintf("---@type %s", typeName))
	} else {
		lines = append(lines, fmt.Sprintf("---@type %s", parentTypeName), "local ref", fmt.Sprintf("ref.%s", typeName))
	}

	err = buffer.SetLines(lines)

	if err != nil {
		return nil, err
	}

	lastLineIndex := len(lines) - 1
	lastLine := lines[lastLineIndex]
	line, character := uint(lastLineIndex), uint(len(lastLine))
	hover, err := buffer.GetHover(line, character)

	if err != nil {
		return nil, err
	}

	return hover, nil
}
