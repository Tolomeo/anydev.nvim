package nvim

import (
	"slices"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type TsQueryMatch []treesitter.Capture

func (tm TsQueryMatch) Find(captureId string) (treesitter.Capture, bool) {
	return slicesx.FindFunc(tm, func(capture treesitter.Capture) bool {
		return captureId == capture.Id
	})
}

func (tm TsQueryMatch) FindAll(captureId string) ([]treesitter.Capture, bool) {
	filtered, _ := slicesx.FilterFunc(tm, func(capture treesitter.Capture) (bool, error) {
		return captureId == capture.Id, nil
	})

	return filtered, len(filtered) > 0
}

func (tm TsQueryMatch) Omit(captureId string) TsQueryMatch {
	filtered, _ := slicesx.FilterFunc(tm, func(capture treesitter.Capture) (bool, error) {
		return captureId != capture.Id, nil
	})

	return filtered
}

func (tm TsQueryMatch) Append(captures ...treesitter.Capture) TsQueryMatch {
	return append(tm, captures...)
}

func (m *TsQueryMatch) LineRange() *treesitter.LineRange {
	startLines, _ := slicesx.MapFunc(*m, func(capture treesitter.Capture) (float64, error) {
		return capture.Node.Range.Start.Line, nil
	})
	endLines, _ := slicesx.MapFunc(*m, func(capture treesitter.Capture) (float64, error) {
		return capture.Node.Range.End.Line, nil
	})

	return &treesitter.LineRange{
		Start: slices.Min(startLines),
		End:   slices.Max(endLines),
	}
}

func (m *TsQueryMatch) BiggestCapture() treesitter.Capture {
	result := (*m)[0]
	resultSize := len(result.Node.Text)

	for _, capture := range (*m)[1:] {
		captureSize := len(capture.Node.Text)

		if captureSize > resultSize {
			result = capture
			resultSize = captureSize
		}
	}

	return result
}

// TODO: create a custom slicesx.MaxFunc util
func (m *TsQueryMatch) Range() *treesitter.Range {
	lineRange := m.LineRange()

	startRangeCaptures, _ := slicesx.FilterFunc(*m, func(capture treesitter.Capture) (bool, error) {
		return capture.Node.Range.Start.Line == lineRange.Start, nil
	})

	startCharacters, _ := slicesx.MapFunc(startRangeCaptures, func(capture treesitter.Capture) (float64, error) {
		return capture.Node.Range.Start.Character, nil
	})

	endRangeCaptures, _ := slicesx.FilterFunc(*m, func(capture treesitter.Capture) (bool, error) {
		return capture.Node.Range.End.Line == lineRange.End, nil
	})

	endCharacters, _ := slicesx.MapFunc(endRangeCaptures, func(capture treesitter.Capture) (float64, error) {
		return capture.Node.Range.End.Character, nil
	})

	return &treesitter.Range{
		Start: treesitter.Position{
			Line:      lineRange.Start,
			Character: slices.Min(startCharacters),
		},
		End: treesitter.Position{
			Line:      lineRange.End,
			Character: slices.Max(endCharacters),
		},
	}
}

type TsQueryMatches []TsQueryMatch

func (ms *TsQueryMatches) BiggestMatch() TsQueryMatch {
	result := (*ms)[0]
	resultSize := len(result.BiggestCapture().Node.Text)

	for _, match := range (*ms)[1:] {
		matchSize := len(match.BiggestCapture().Node.Text)

		if matchSize > resultSize {
			result = match
			resultSize = matchSize
		}
	}

	return result
}
