package origin

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type FieldAnnotationOrigin struct {
	origin
	node annotation.AtField
}

func (fo *FieldAnnotationOrigin) Definition() []string {
	return strings.Split(fo.node.Root().Text, "\n")
}

func (fo *FieldAnnotationOrigin) Name() *string {
	if fo.node.Name() == nil {
		return nil
	}

	return &fo.node.Name().Text
}

func (fo *FieldAnnotationOrigin) Index() *string {
	if fo.node.Index() == nil {
		return nil
	}

	return &fo.node.Index().Text
}

func (fo *FieldAnnotationOrigin) Private() bool {
	return fo.node.Private() != nil
}

func (fo *FieldAnnotationOrigin) Protected() bool {
	return fo.node.Protected() != nil
}

func (fo *FieldAnnotationOrigin) Package() bool {
	return fo.node.Package() != nil
}

func (fo *FieldAnnotationOrigin) Optional() bool {
	return fo.node.Optional() != nil
}

func (fo *FieldAnnotationOrigin) Type() string {
	if fo.Optional() {
		return fmt.Sprintf("(%s)?", fo.node.Type().Text)
	}

	return fo.node.Type().Text
}

func NewFieldAnnotationOrigin(location nvim.Location, node annotation.AtField) *FieldAnnotationOrigin {
	fieldOrigin := FieldAnnotationOrigin{
		origin: origin{
			location: location,
		},
		node: node,
	}

	return &fieldOrigin
}
