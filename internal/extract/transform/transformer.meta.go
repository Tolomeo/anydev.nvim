package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
)

type virtualOriginAtAnnotations struct {
	AtType *AtTypeAnnotation
}

func (tr *Transformer) getVirtualOriginAnnotations(docblock []string) (virtualOriginAtAnnotations, error) {
	annotations := virtualOriginAtAnnotations{
		AtType: nil,
	}

	buffer, err := tr.context.Nvim().NewBuffer()

	if err != nil {
		return annotations, err
	}

	defer buffer.Close()

	err = buffer.SetLines(docblock)

	if err != nil {
		return annotations, err
	}

	atType, err := tr.getAtTypeAnnotations(buffer)

	if err != nil {
		return annotations, fmt.Errorf("Error lexing type annotation: %w", err)
	}

	annotations.AtType = atType
	return annotations, nil
}

func (tr *Transformer) getVirtualOriginType(metaOrigin *origin.VirtualOrigin) (annotation.Type, error) {
	unknown := symbol.NewUnknown()
	unknown.Documentation = metaOrigin.Annotations()

	annotations, err := tr.getVirtualOriginAnnotations(metaOrigin.Annotations())

	if err != nil {
		return nil, err
	}

	if annotations.AtType == nil {
		tr.context.Logger().Warn(fmt.Sprintf("Unknown meta type '%s' received", tr.context.Target().Name()))
		return unknown, nil
	}

	if len(annotations.AtType.Types) < 1 {
		return nil, fmt.Errorf("Error lexing @type annotations for meta type '%s': no type annotations found", tr.context.Target().Name())
	}

	lexedType, err := tr.getType(annotations.AtType.Types[0])

	if err != nil {
		return nil, err
	}

	return lexedType, nil
}
