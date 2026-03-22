package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getValueOriginType(valueOrigin *origin.ValueOrigin) (symbol.Type, error) {
	buffer, err := tr.context.Nvim().OpenTemporary()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(valueOrigin.Annotations())

	if err != nil {
		return nil, err
	}

	annotations, err := tr.getTypeAtAnnotations(buffer)

	if err != nil {
		return nil, err
	}

	if annotations == nil {
		valueType, err := tr.transformType(valueOrigin.Type())

		if err != nil {
			return nil, err
		}

		return valueType, nil
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

	return nil, fmt.Errorf("Unreachable: failed to transform value origin <%+v> with annotations <%+v>", valueOrigin, annotations)
}
