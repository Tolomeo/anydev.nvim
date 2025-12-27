package languageserver

import "github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"

func (r *Range) AsTreesitter() *treesitter.Range {
	return &treesitter.Range{
		Start: treesitter.Position{
			Line:      r.Start.Line,
			Character: r.End.Character,
		},
		End: treesitter.Position{
			Line:      r.End.Line,
			Character: r.End.Character,
		},
	}
}

func (r *Range) Contains(rng Range) bool {
	if r.Start.Line == rng.Start.Line &&
		r.Start.Character > rng.Start.Character {
		return false
	}

	if r.End.Line == rng.End.Line &&
		r.End.Character < rng.End.Character {
		return false
	}

	return r.Start.Line <= rng.Start.Line &&
		r.End.Line >= rng.End.Line
}

func (r *Range) StartPosition() Position {
	return Position{
		Line: r.Start.Line,
		Character: r.Start.Character,
	}
}

func (r *Range) EndPosition() Position {
	return Position{
		Line: r.End.Line,
		Character: r.End.Character,
	}
}

func (r *Range) LineRange() LineRange {
	return LineRange{
		Start: r.Start.Line,
		End:   r.End.Line,
	}
}
