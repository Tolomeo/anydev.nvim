package crawl

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/languageserver"
)

func (c *Crawler) getDefinitionDocumentation(identifier string) (*languageserver.MarkupContent, error) {
	buffer, err := c.context.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	lines := []string{}
	identifierParts := strings.Split(identifier, ".")

	if len(identifierParts) > 1 {
		tail := identifierParts[len(identifierParts)-1]
		_, isKeyword := definition.Keywords[tail]

		if isKeyword {
			head := strings.Join(identifierParts[:len(identifierParts)-1], ".")
			lines = append(lines, fmt.Sprintf("%s['%s']", head, tail))
		} else {
			lines = append(lines, identifier)
		}
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

func (c *Crawler) getTypeDefinitionDocumentation(typeName string, parentTypeName string) (*languageserver.MarkupContent, error) {
	buffer, err := c.context.Nvim().NewBuffer()

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
