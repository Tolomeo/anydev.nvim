package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

type variableOriginAtAnnotations struct {
	AtType *annotation.AtType
}

func (tr *Transformer) getVariableOriginAnnotations(docblock []string) (variableOriginAtAnnotations, error) {
	annotations := variableOriginAtAnnotations{
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

func (tr *Transformer) getVariableOriginType(variableOrigin *origin.VariableOrigin) (symbol.Type, error) {
	annotations, err := tr.getVariableOriginAnnotations(variableOrigin.Annotations())

	if err != nil {
		return nil, err
	}

	if annotations.AtType == nil {
		return tr.context.Follow(variableOrigin.Name())
	}

	if len(annotations.AtType.Types()) < 1 {
		return nil, fmt.Errorf("Error lexing @type annotations for meta type '%s': no type annotations found", tr.context.Target().Name())
	}

	lexedType, err := tr.getType(annotations.AtType.Types()[0])

	if err != nil {
		return nil, err
	}

	return lexedType, nil
}
