package origin

type OriginChain []Origin

func (o OriginChain) Last() Origin {
	return o[len(o)-1]
}

func (o OriginChain) First() Origin {
	return o[0]
}

func (o OriginChain) Append(o2 OriginChain) OriginChain {
	o = append(o, o2...)
	return o
}

func NewOriginChain(origins ...Origin) OriginChain {
	originChain := OriginChain{}

	for _, o := range origins {
		originChain = append(originChain, o)
	}

	return originChain
}
