package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type EnumeratorAnnotationOrigin struct {
	origin
	node        annotation.AtEnum
	memberNodes []annotation.AtEnumMember
}

func (eo *EnumeratorAnnotationOrigin) Definition() []string {
	return strings.Split(eo.node.Root().Text, "\n")
}

func (eo *EnumeratorAnnotationOrigin) Name() string {
	return eo.node.Name().Text
}

func (eo *EnumeratorAnnotationOrigin) Members() []annotation.AtEnumMember {
	return eo.memberNodes
}

func NewEnumAnnotationOrigin(location nvim.Location, node annotation.AtEnum, memberNodes ...annotation.AtEnumMember) *EnumeratorAnnotationOrigin {
	enumOrigin := EnumeratorAnnotationOrigin{
		origin: origin{
			location: location,
		},
		node:        node,
		memberNodes: memberNodes,
	}

	return &enumOrigin
}
