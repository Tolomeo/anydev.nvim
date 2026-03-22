package transform

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) transformFunctionCallOrigin(functionCallOrigin *origin.FunctionCallOrigin) (symbol.Type, error) {
	buffer, err := tr.context.Nvim().OpenTemporary()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(functionCallOrigin.Annotations())

	if err != nil {
		return nil, err
	}

	annotations, err := tr.getTypeAtAnnotations(buffer)

	if err != nil {
		return nil, err
	}

	if annotations != nil && annotations.AtType != nil {
		// NB: we don't check here for the presence of multiple types
		originType, err := tr.transformType(annotations.AtType.Types()[0])

		if err != nil {
			return nil, err
		}

		return originType, nil
	}

	switch functionCallOrigin.FunctionName() {
	case "require":
		tr.context.Logger().Verbosef("Extracting 'require' function call as a module reference")
		moduleNameArgument := functionCallOrigin.FunctionArguments()[0]
		moduleName := strings.Trim(moduleNameArgument, "'\"")
		err := tr.context.Extract("module", strings.Trim(moduleNameArgument, "'\""))
		if err != nil {
			return nil, err
		}
		return symbol.NewModuleReference(moduleName), nil

	case "vim._defer_deprecated_module":
		tr.context.Logger().Verbosef("Extracting 'vim._defer_deprecated_module' function call as a module reference")
		moduleNameArgument := functionCallOrigin.FunctionArguments()[1]
		moduleName := strings.Trim(moduleNameArgument, "'\"")
		err := tr.context.Extract("module", moduleName)
		if err != nil {
			return nil, err
		}
		return symbol.NewModuleReference(moduleName), nil

	case "vim._defer_require", "setmetatable", "create_option_accessor":
		tr.context.Logger().Verbosef("Extracting '%s' function call as a table", functionCallOrigin.Name())
		annotations, err := tr.getTableAnnotations(functionCallOrigin.Annotations())
		if err != nil {
			return nil, err
		}
		name := functionCallOrigin.Name()
		if annotations.AtClass != nil {
			return tr.getClassTableSymbol(name, *annotations.AtClass)
		}
		// NOTE: enumerators not accounted for in here
		/* if annotations.AtEnum != nil {
			return tr.getEnumeratorTableSymbol(name, *annotations.AtEnum)
		} */
		return tr.getTableSymbol(name)
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
