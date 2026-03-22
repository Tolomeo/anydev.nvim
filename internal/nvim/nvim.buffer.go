package nvim

import (
	"fmt"
	"path"
	"regexp"
	"slices"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/msgpackrpc"
	"github.com/Tolomeo/anydev.nvim/internal/utils/anyx"
)

func (n *Nvim) edit(file string) (string, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_command",
		Params: []any{fmt.Sprintf("edit %s", file)},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return file, fmt.Errorf("Error opening %s: %v\n", file, err)
	}

	return file, nil
}

/* func (n *Nvim) write() error {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_command",
		Params: []any{"write"},
	}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to write buffer: %v\n", err)
	}

	return nil
} */

/* func (n *Nvim) getBufferName() (string, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_buf_get_name",
		Params: []any{0},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return "", fmt.Errorf("Error reading buffer name: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return "", fmt.Errorf("Error executing lua: %v\n", err)
	}

	return result.(string), nil
} */

func (n *Nvim) setBufferLines(lines []string) error {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_buf_set_lines",
		Params: []any{0, 0, -1, true, lines},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error sending nvim_buf_set_lines rpc message: %v\n", err)
	}

	_, err = response.Result()

	if err != nil {
		return fmt.Errorf("Error setting buffer lines: %v\n", err)
	}

	return nil
}

/* func (n *Nvim) getBufferText(startRow int, startCol int, endRow int, endCol int) ([]string, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_buf_get_text",
		Params: []any{0, startRow, startCol, endRow, endCol, struct{}{}},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer text: %w\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer text: %w\n", err)
	}

	bufferText, err := anyx.ToSliceOf[string](result)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer text return value: %w", err)
	}

	return bufferText, nil
} */

func (n *Nvim) getBufferLines(start int, end int) ([]string, error) {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_buf_get_lines",
		Params: []any{0, start, end, false},
	}
	response, err := n.rpc.Send(request)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer name: %v\n", err)
	}

	result, err := response.Result()

	if err != nil {
		return []string{}, fmt.Errorf("Error executing lua: %v\n", err)
	}

	bufferLines, err := anyx.ToSliceOf[string](result)

	if err != nil {
		return []string{}, fmt.Errorf("Error reading buffer lines return value: %w", err)
	}

	return bufferLines, nil
}

func (n *Nvim) deleteBuffer() error {
	request := msgpackrpc.RequestMessage{
		Method: "nvim_buf_delete",
		Params: []any{0, struct{ force bool }{force: true}}}
	_, err := n.rpc.Send(request)

	if err != nil {
		return fmt.Errorf("Error trying to delete buffer: %v\n", err)
	}

	return nil
}

func (nvim *Nvim) OpenFile(name string) (*FileBuffer, error) {
	if bufIndex := slices.IndexFunc(nvim.files[:], func(buf *FileBuffer) bool {
		return buf.name == name
	}); bufIndex != -1 {
		nvim.logger.Sillyf("Reopening %s buffer", name)
		return nvim.files[bufIndex], nil
	}

	nvim.logger.Sillyf("Opening %s buffer", name)

	buffer := &FileBuffer{
		buffer: buffer{
			nvim: nvim,
			name: name,
		},
	}

	nvim.files = append([]*FileBuffer{buffer}, nvim.files...)

	if len(nvim.files) >= 30 {
		for _, buf := range nvim.files[len(nvim.files)-10:] {
			_, err := nvim.edit(buf.name)

			if err != nil {
				return nil, err
			}

			err = nvim.deleteBuffer()

			if err != nil {
				return nil, err
			}

		}

		nvim.files = nvim.files[:len(nvim.files)-10]
	}

	_, err := nvim.edit(name)

	if err != nil {
		return nil, err
	}

	return buffer, nil
}

var scratchBufferCounter = 0
var scratchBufferName = regexp.MustCompile(`anydev\.\d+\.lua$`)

func (nvim *Nvim) OpenTemporary() (*TemporaryBuffer, error) {
	scratchBufferCounter += 1
	name := path.Join(nvim.config, fmt.Sprintf("anydev.%d.lua", scratchBufferCounter))

	nvim.logger.Sillyf("Opening %s scratch buffer", name)

	buffer := &TemporaryBuffer{
		buffer: buffer{
			nvim: nvim,
			name: name,
		},
	}

	_, err := nvim.edit(name)

	if err != nil {
		return nil, err
	}

	return buffer, nil
}
