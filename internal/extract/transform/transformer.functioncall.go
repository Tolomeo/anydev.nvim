package transform

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

type functionCallOriginAtAnnotations struct {
	AtType *annotation.AtType
	functionAtAnnotations
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
		originType, err := tr.getType(annotations.AtType.Types()[0])

		if err != nil {
			return nil, err
		}

		return originType, nil
	}

	switch functionCallOrigin.FunctionName() {
	case "vim._defer_deprecated_module":
		tr.context.Logger().Verbosef("Extracting 'vim._defer_deprecated_module' function call as a module require")
		moduleNameArgument := functionCallOrigin.FunctionArguments()[1]
		err := tr.context.Extract("module", strings.Trim(moduleNameArgument, "'\""))
		if err != nil {
			return nil, err
		}
		return symbol.NewModuleReference(moduleNameArgument), nil

	case "vim._defer_require", "setmetatable", "create_option_accessor":
		tr.context.Logger().Verbosef("Extracting '%s' function call as a table", functionCallOrigin.Name())
		table := symbol.NewTable()
		table.Name = functionCallOrigin.Name()
		tableChildren, err := tr.context.Nvim().GetValueCompletion(tr.context.Target().Identifier())

		if err != nil {
			return nil, err
		}

		for _, child := range tableChildren {
			err := tr.context.ExtractChild(table, child)

			if err != nil {
				return nil, err
			}
		}

		return table, nil
	}

	followedType, err := tr.context.Follow(functionCallOrigin.FunctionName())

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
