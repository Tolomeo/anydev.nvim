package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

type functionCallOriginAtAnnotations struct {
	AtType *AtTypeAnnotation
}

func (tr *Transformer) getFunctionCallOriginAtAnnotations(docblock []string) (functionCallOriginAtAnnotations, error) {
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

func (tr *Transformer) getFunctionCallOriginType(functionCallOrigin *origin.FunctionCallOrigin) (symbol.Type, error) {
	annotations, err := tr.getFunctionCallOriginAtAnnotations(functionCallOrigin.Annotations())

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

	followedType, err := tr.context.Follow()

	// TODO: support generics?
	switch functionType := followedType.(type) {
	case *symbol.Function:
		if len(functionType.Returns) == 0 {
			tr.context.Logger().Warnf("Unknown return type of function call origin: <%+v>", functionType)
			return symbol.NewUnknown(), nil
		}

		if len(functionType.Returns) > 1 {
			tr.context.Logger().Errorf("Unsupported multiple return types of function call origin: <%+v>", functionType)
			return symbol.NewUnknown(), nil
		}

		return functionType.Returns[0].Type, nil
	}

	return nil, fmt.Errorf("Unsupported followed type result <%T> of function call origin", followedType)
}
