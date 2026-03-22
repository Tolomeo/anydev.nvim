package treesitter

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

func (r *Range) LineRange() LineRange {
	return LineRange{
		Start: r.Start.Line,
		End:   r.End.Line,
	}
}
