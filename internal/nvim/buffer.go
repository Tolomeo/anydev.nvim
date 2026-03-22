package nvim

import (
	"fmt"
	"iter"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/languageserver"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type Buffer interface {
	Name() string
	ReadLines() ([]string, error)
	SetLines([]string) error
	GetTSNodeAt(nodeTypes []string, line uint, character uint) (*treesitter.TsNode, error)
	TsQueryOne(query treesitter.Query) (*TsQueryMatch, error)
	SafeTsQueryOne(query treesitter.Query) (*SafeTsQueryResult, error)
	TsQueryAll(query treesitter.Query) (*TsQueryMatches, error)
	SafeTsQueryAll(query treesitter.Query) (*[]SafeTsQueryResult, error)
	GetDefinitionLocations(line uint, character uint) (*[]Location, error)
	GetTypeDefinitionLocations(line uint, character uint) (*[]Location, error)
	GetHover(line uint, character uint) (*languageserver.MarkupContent, error)
	QueryTsNodeAt(tsNodeQueryMap TsNodeQueryMap, line uint, character uint) (*TsNodeQueryMatch, error)
	NextLineIterator(startLine uint) (iter.Seq2[*ScratchBuffer, error], error)
	Close() error
}

type buffer struct {
	nvim *Nvim
	name string
}

func (b *buffer) Name() string {
	return b.name
}

func (b *buffer) ReadLines() ([]string, error) {
	b.nvim.logger.Sillyf("Reading %s buffer lines", b.name)
	_, err := b.nvim.edit(b.name)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer '%s': %w", b.name, err)
	}

	lines, err := b.nvim.getBufferLines(0, -1)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer '%s': %w", b.name, err)
	}

	return lines, nil
}

func (b *buffer) SetLines(lines []string) error {
	b.nvim.logger.Sillyf("Writing %s buffer lines", b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return fmt.Errorf("Error writing to buffer '%s': %w", b.name, err)
	}

	if len(lines) < 1 {
		return b.nvim.setBufferLines([]string{""})
	}

	return b.nvim.setBufferLines(lines)
}

func (b *buffer) GetTSNodeAt(nodeTypes []string, line uint, character uint) (*treesitter.TsNode, error) {
	b.nvim.logger.Sillyf("Getting treesitter node <%s> in %s buffer", nodeTypes, b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.getTSNodeAt(nodeTypes, line, character)
}

func (b *buffer) GetTsCommentBlockAt(line uint, character uint) (*treesitter.TsNode, error) {
	b.nvim.logger.Sillyf("Getting comment block node in %s buffer", b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.getTsCommentBlockAt(line, character)
}

func (b *buffer) GetTsNodeAnnotations(node treesitter.TsNode) ([]string, error) {
	b.nvim.logger.Sillyf("Getting node <%v> annotations in %s buffer", node, b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.getNodeAnnotations(node)
}

func (b *buffer) TsQueryOne(query treesitter.Query) (*TsQueryMatch, error) {
	b.nvim.logger.Sillyf("Executing query one in %s buffer", b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.tsQueryOne(query)
}

func (b *buffer) TsQueryAll(query treesitter.Query) (*TsQueryMatches, error) {
	b.nvim.logger.Sillyf("Executing query all in %s buffer", b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.tsQueryAll(query)
}

func (b *buffer) SafeTsQueryOne(query treesitter.Query) (*SafeTsQueryResult, error) {
	b.nvim.logger.Sillyf("Executing safe query one in %s buffer", b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.safeTsQueryOne(query)
}

func (b *buffer) SafeTsQueryAll(query treesitter.Query) (*[]SafeTsQueryResult, error) {
	b.nvim.logger.Sillyf("Executing safe query all in %s buffer", b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.safeTsQueryAll(query)
}

func (b *buffer) GetDefinitionLocations(line uint, character uint) (*[]Location, error) {
	b.nvim.logger.Sillyf("Getting definition locations in %s buffer", b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.getDefinitionLocations(line, character)
}

func (b *buffer) GetTypeDefinitionLocations(line uint, character uint) (*[]Location, error) {
	b.nvim.logger.Sillyf("Getting type definition locations in %s buffer", b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.getTypeDefinitionLocations(line, character)
}

func (b *buffer) GetHover(line uint, character uint) (*languageserver.MarkupContent, error) {
	b.nvim.logger.Sillyf("Getting lsp hover at %s:%d:%d", b.name, line, character)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.getLspHover(line, character)
}

func (b *buffer) QueryTsNodeAt(tsNodeQueryMap TsNodeQueryMap, line uint, character uint) (*TsNodeQueryMatch, error) {
	b.nvim.logger.Sillyf("Querying ts node in %s buffer", b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.queryTsNodeAt(tsNodeQueryMap, line, character)
}

func (b *buffer) NextLineIterator(startLine uint) (iter.Seq2[*ScratchBuffer, error], error) {
	_, err := b.nvim.edit(b.name)

	if err != nil {
		return nil, err
	}

	lines, err := b.nvim.getBufferLines(0, -1)

	// fmt.Println(lines)

	if err != nil {
		return nil, err
	}

	linebuffer, err := b.nvim.OpenScratchBuffer()

	if err != nil {
		return nil, err
	}

	return func(yield func(*ScratchBuffer, error) bool) {
		for i := startLine; i < uint(len(lines)); i++ {

			err := linebuffer.SetLines([]string{lines[i]})

			if err != nil && !yield(linebuffer, err) {
				linebuffer.Close()
				return
			}

			if !yield(linebuffer, nil) {
				linebuffer.Close()
				return
			}
		}

		linebuffer.Close()
	}, nil
}

func (b *buffer) Close() error {
	b.nvim.logger.Sillyf("Closing %s buffer", b.name)

	_, err := b.nvim.edit(b.name)

	if err != nil {
		return err
	}

	err = b.nvim.deleteBuffer()

	if err != nil {
		return err
	}

	return nil
}

type ScratchBuffer struct {
	buffer
}

type FileBuffer struct {
	buffer
}

func (f *FileBuffer) SetLines() error {
	return fmt.Errorf("Error writing to a file buffer: denied")
}

func (f *FileBuffer) Close() error {
	f.nvim.logger.Sillyf("Delegating '%s' buffer close to nvim queue", f.name)
	return nil
}
