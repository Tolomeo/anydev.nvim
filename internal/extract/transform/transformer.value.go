package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

type valueOriginAtAnnotations struct {
	AtType *annotation.AtType
}

func (tr *Transformer) getValueOriginAnnotations(docblock []string) (valueOriginAtAnnotations, error) {
	annotations := valueOriginAtAnnotations{
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

func (tr *Transformer) getValueOriginType(valueOrigin *origin.ValueOrigin) (symbol.Type, error) {
	annotations, err := tr.getValueOriginAnnotations(valueOrigin.Annotations())

	if err != nil {
		return nil, err
	}

	if annotations.AtType != nil {
		if len(annotations.AtType.Types()) < 1 {
			return nil, fmt.Errorf("Error lexing @type annotations for meta type '%s': no type annotations found", tr.context.Target().Name())
		}

		lexedType, err := tr.transformType(annotations.AtType.Types()[0])

		if err != nil {
			return nil, err
		}

		return lexedType, nil
	}

	valueType, err := tr.transformType(valueOrigin.Type())

	if err != nil {
		return nil, err
	}

	return valueType, nil
}
