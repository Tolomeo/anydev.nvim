package origin

type OriginChain []Origin

func (o *OriginChain) Last() Origin {
	return (*o)[len(*o)-1]
}

func (o *OriginChain) First() Origin {
	return (*o)[0]
}

func (o *OriginChain) Merge(o2 *OriginChain) *OriginChain {
	*o = append(*o, *o2...)
	return o
}

func NewOrigins(origins ...Origin) *OriginChain {
	t := OriginChain{}

	for _, l := range origins {
		t = append(t, l)
	}

	return &t
}
