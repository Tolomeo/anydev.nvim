package origin

import (
	"slices"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type OriginChain []Origin

func (o OriginChain) Last() Origin {
	return o[len(o)-1]
}

func (o OriginChain) First() Origin {
	return o[0]
}

func (o OriginChain) Append(origins ...Origin) OriginChain {
	o = append(o, origins...)
	return o
}

func (o OriginChain) Concat(o2 OriginChain) OriginChain {
	o = append(o, o2...)
	return o
}

func (o OriginChain) ContainsLocation(location nvim.Location) bool {
	locationUrl, locationLine, locationCharacter := location.Url, location.StartLine(), location.StartCharacter()

	return slices.ContainsFunc(o, func(originChainOrigin Origin) bool {
		return originChainOrigin.Url() == locationUrl &&
			originChainOrigin.Line() == locationLine &&
			originChainOrigin.Character() == locationCharacter
	})
}

func NewOriginChain(origins ...Origin) OriginChain {
	originChain := OriginChain{}

	for _, o := range origins {
		originChain = append(originChain, o)
	}

	return originChain
}
