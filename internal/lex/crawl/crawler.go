package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/context"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/anyx"
)

type logger interface {
	Info(message string)
	Warn(message string)
	Error(message string)
}

type Crawler struct {
	context *context.Context
}

func (c *Crawler) Source(path string) (*Source, error) {
	source := Source{path: path}
	err := c.sourceValue(path, &source)

	if err != nil {
		return nil, err
	}

	return &source, nil
}

func (c *Crawler) SourceType(name string) (*Source, error) {
	source := Source{path: name}
	err := c.sourceType(name, &source)

	if err != nil {
		return nil, err
	}

	return &source, nil
}

func (c *Crawler) readCommentBlock(line uint, character uint) (*[]string, error) {
	lines, err := c.context.Nvim.GetBufferLines(int(line)-1, int(line))

	switch {
	case err != nil:
		return nil, err
	case len(lines) < 1:
		return nil, nil
	}

	// clamping the received character to be inside the line
	character = max(0, min(character, uint(len(lines[0])-1)))

	node, err := c.context.Nvim.GetTSNodeAt([]string{treesitter.COMMENT}, line, character)

	switch {
	case err != nil:
		return nil, err
	case node == nil:
		return nil, nil
	}

	luaCode := `
		local args = {...}
		local startLine = args[1]
		local endLine = args[2]

		local function is_comment(ln)
			local rest = ln:match("^%s*(.*)")
			return rest:sub(1,2) == "--"
		end

		local previous_line = vim.api.nvim_buf_get_lines(0, startLine -1, startLine, false)[1]

		while previous_line and is_comment(previous_line) do
			startLine = startLine - 1
			previous_line = vim.api.nvim_buf_get_lines(0, startLine -1, startLine, false)[1]
		end

		local next_line = vim.api.nvim_buf_get_lines(0, endLine + 1, endLine + 2, false)[1]

		while next_line and is_comment(next_line) do
			endLine = endLine + 1
			next_line = vim.api.nvim_buf_get_lines(0, endLine + 1, endLine + 2, false)[1]
		end

		return vim.api.nvim_buf_get_lines(0, startLine, endLine + 1, true)
	`

	result, err := c.context.Nvim.ExecLua(luaCode, []any{node.Range.Start.Line, node.Range.End.Line})

	if err != nil {
		return nil, err
	}

	bufferLines, err := anyx.ToSliceOf[string](result)

	if err != nil {
		return nil, fmt.Errorf("Error reading buffer lines return value: %w", err)
	}

	return &bufferLines, nil

}

func NewCrawler(context *context.Context) *Crawler {
	return &Crawler{
		context: context,
	}
}
