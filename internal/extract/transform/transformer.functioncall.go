package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

type functionCallOriginAtAnnotations struct {
	AtType *AtTypeAnnotation
}

func (tr *Transformer) getFunctionCallOriginAnnotations(docblock []string) (functionCallOriginAtAnnotations, error) {
	annotations := functionCallOriginAtAnnotations{
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

func (tr *Transformer) getFunctionCallOriginType(functionCallOrigin *origin.FunctionCallOrigin) (annotation.Type, error) {
	annotations, err := tr.getVariableOriginAnnotations(functionCallOrigin.Annotations())

	if err != nil {
		return nil, err
	}

	// NB: we don't check here for the presence of multiple types
	if annotations.AtType != nil {
		originType, err := tr.getType(annotations.AtType.Types[0])

		if err != nil {
			return nil, err
		}

		return originType, nil
	}

	return symbol.NewUnknown(), nil
}
