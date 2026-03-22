package transform

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getVirtualOriginType(virtualOrigin *origin.VirtualOrigin) (symbol.Type, error) {
	buffer, err := tr.context.Nvim().OpenTemporary()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(virtualOrigin.Annotations())

	if err != nil {
		return nil, err
	}

	annotations, err := tr.getTypeAtAnnotations(buffer)

	if err != nil {
		return nil, err
	}

	if annotations == nil {
		tr.context.Logger().Warn(fmt.Sprintf("Unknown meta type '%s' received", tr.context.Target().Name()))
		return symbol.NewUnknown(), nil
	}

	if annotations.AtType != nil {
		// NB: we don't check for the presence of multiple types here
		lexedType, err := tr.transformType(annotations.AtType.Types()[0])

		if err != nil {
			return nil, err
		}

		return lexedType, nil
	}

	if annotations.AtModule != nil {
		moduleName := strings.Trim(annotations.AtModule.Name().Text, "'\"")
		err := tr.context.Extract("module", moduleName)

		if err != nil {
			return nil, err
		}

		return symbol.NewModuleReference(moduleName), nil
	}

	return nil, fmt.Errorf("Unreachable: failed to transform virtual type receved <%+v> with type annotations <%+v>", virtualOrigin, annotations)
}
