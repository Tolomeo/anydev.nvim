package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getVariableOriginType(variableOrigin *origin.VariableOrigin) (symbol.Type, error) {
	buffer, err := tr.context.Nvim().NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(variableOrigin.Annotations())

	if err != nil {
		return nil, err
	}

	annotations, err := tr.getTypeAtAnnotations(buffer)

	if err != nil {
		return nil, err
	}

	if annotations == nil {
		return tr.context.Follow(variableOrigin.Name())
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

	if annotations.AtClass != nil {
		className := annotations.AtClass.Name().Text

		table := symbol.NewTable()
		table.Name = className

		err := tr.context.Extract("type", className)

		if err != nil {
			return nil, err
		}

		return symbol.NewTypeReference(className), nil
	}

	if annotations.AtOverload != nil {
		functionType, err := tr.transformType(annotations.AtOverload.Type())

		if err != nil {
			return nil, err
		}

		return functionType, nil
	}

	return nil, fmt.Errorf("Unreachable: failed extracting variable origin type <%+v> from defined annotations <%+v>", variableOrigin, annotations)
}
