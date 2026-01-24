package nvim

import (
	"fmt"
	"iter"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

type Buffer struct {
	nvim *Nvim
	name string
}

func (b *Buffer) ReadLines() ([]string, error) {
	b.nvim.logger.Verbosef("Reading %s buffer lines", b.name)
	_, err := b.nvim.open(b.name)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer '%s': %w", b.name, err)
	}

	lines, err := b.nvim.getBufferLines(0, -1)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer '%s': %w", b.name, err)
	}

	return lines, nil
}

func (b *Buffer) SetLines(lines []string) error {
	b.nvim.logger.Verbosef("Writing %s buffer lines", b.name)

	_, err := b.nvim.open(b.name)

	if err != nil {
		return fmt.Errorf("Error writing to buffer '%s': %w", b.name, err)
	}

	if len(lines) < 1 {
		return b.nvim.setBufferLines([]string{""})
	}

	return b.nvim.setBufferLines(lines)
}

func (b *Buffer) GetTSNodeAt(nodeTypes []string, line uint, character uint) (*treesitter.TsNode, error) {
	b.nvim.logger.Verbosef("Getting treesitter node <%s> in %s buffer", nodeTypes, b.name)

	_, err := b.nvim.open(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.getTSNodeAt(nodeTypes, line, character)
}

func (b *Buffer) GetTsCommentBlockAt(line uint, character uint) (*treesitter.TsNode, error) {
	b.nvim.logger.Verbosef("Getting comment block node in %s buffer", b.name)

	_, err := b.nvim.open(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.getTsCommentBlockAt(line, character)
}

func (b *Buffer) TsQueryOne(query treesitter.Query) (*TsQueryMatch, error) {
	b.nvim.logger.Verbosef("Executing query one in %s buffer", b.name)

	_, err := b.nvim.open(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.tsQueryOne(query)
}

func (b *Buffer) TsQueryAll(query treesitter.Query) (*[]TsQueryMatch, error) {
	b.nvim.logger.Verbosef("Executing query all in %s buffer", b.name)

	_, err := b.nvim.open(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.tsQueryAll(query)
}

func (b *Buffer) SafeTsQueryAll(query treesitter.Query) (*[]SafeTsQueryResult, error) {
	b.nvim.logger.Verbosef("Executing safe query all in %s buffer", b.name)

	_, err := b.nvim.open(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.safeTsQueryAll(query)
}

func (b *Buffer) GetDefinitionLocations(line uint, character uint) (*[]Location, error) {
	b.nvim.logger.Verbosef("Getting definition locations in %s buffer", b.name)

	_, err := b.nvim.open(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.getDefinitionLocations(line, character)
}

func (b *Buffer) GetTypeDefinitionLocations(line uint, character uint) (*[]Location, error) {
	b.nvim.logger.Verbosef("Getting type definition locations in %s buffer", b.name)

	_, err := b.nvim.open(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.getTypeDefinitionLocations(line, character)
}

func (b *Buffer) QueryTsNodeAt(tsNodeQueryMap TsNodeQueryMap, line uint, character uint) (*TsNodeQueryMatch, error) {
	b.nvim.logger.Verbosef("Querying ts node in %s buffer", b.name)

	_, err := b.nvim.open(b.name)

	if err != nil {
		return nil, err
	}

	return b.nvim.queryTsNodeAt(tsNodeQueryMap, line, character)
}

func (b *Buffer) NextLineIterator(startLine uint) (iter.Seq2[*ScratchBuffer, error], error) {
	_, err := b.nvim.open(b.name)

	if err != nil {
		return nil, err
	}

	lines, err := b.nvim.getBufferLines(0, -1)

	// fmt.Println(lines)

	if err != nil {
		return nil, err
	}

	linebuffer, err := b.nvim.NewBuffer()

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

func (b *Buffer) Close() error {
	b.nvim.logger.Verbosef("Closing %s buffer", b.name)
	/* _, err := b.nvim.open(b.name)

	if err != nil {
		return err
	}

	err = b.nvim.deleteBuffer()

	if err != nil {
		return err
	} */

	return nil
}

func (nvim *Nvim) OpenBuffer(name string) (*Buffer, error) {
	nvim.logger.Verbosef("Opening %s buffer", name)

	buffer := &Buffer{
		nvim: nvim,
		name: name,
	}

	_, err := nvim.open(name)

	if err != nil {
		return nil, err
	}

	return buffer, nil
}

var newBufferCounter = 0

type ScratchBuffer struct {
	Buffer
}

func (s *ScratchBuffer) Close() error {
	s.nvim.logger.Verbosef("Closing %s buffer", s.name)

	_, err := s.nvim.open(s.name)

	if err != nil {
		return err
	}

	err = s.nvim.deleteBuffer()

	if err != nil {
		return err
	}

	return nil
}

func (nvim *Nvim) NewBuffer() (*ScratchBuffer, error) {
	newBufferCounter += 1
	name := path.Join(nvim.Options().Config().Dir(), fmt.Sprintf("anydev.%d.lua", newBufferCounter))

	nvim.logger.Verbosef("Opening %s scratch buffer", name)

	buffer := &ScratchBuffer{
		Buffer: Buffer{
			nvim: nvim,
			name: name,
		},
	}

	_, err := nvim.open(name)

	if err != nil {
		return nil, err
	}

	return buffer, nil
}
