package origin

type Origins []Origin

func (o *Origins) Last() Origin {
	return (*o)[len(*o)-1]
}

func (o *Origins) First() Origin {
	return (*o)[0]
}

func (o *Origins) Merge(o2 *Origins) *Origins {
	*o = append(*o, *o2...)
	return o
}

func NewOrigins(origins ...Origin) *Origins {
	t := Origins{}

	for _, l := range origins {
		t = append(t, l)
	}

	return &t
}
