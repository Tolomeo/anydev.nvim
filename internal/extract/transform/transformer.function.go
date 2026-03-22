package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (tr *Transformer) transformFunctionOrigin(functionOrigin *origin.FunctionOrigin) (*symbol.Function, error) {
	function := symbol.NewFunction()
	function.Name = functionOrigin.Name()

	for _, arg := range functionOrigin.Args() {
		function.Arguments = append(function.Arguments, *symbol.NewFunctionArgument(arg))
	}

	annotations, err := tr.getFunctionAtAnnotations(functionOrigin.Annotations())

	if err != nil {
		return nil, fmt.Errorf("Error lexing function %s: %w", function.Name, err)
	}

	if annotations == nil {
		return function, nil
	}

	for _, genericsAnnotation := range annotations.AtGenerics {
		for _, genericAnnotation := range genericsAnnotation.Generics() {
			generic := symbol.NewFunctionGeneric(genericAnnotation.Name(), nil)

			if genericAnnotation.Type() != nil {
				genericType, err := tr.transformType(*genericAnnotation.Type())

				if err != nil {
					return nil, err
				}

				generic.Type = genericType
			}

			function.Generics = append(function.Generics, *generic)
		}
	}

	for argIndex := range function.Arguments {
		name := function.Arguments[argIndex].Name

		paramAnnotation, hasParamAnnotation := annotations.AtParams[name]

		if !hasParamAnnotation {
			tr.context.Logger().Warnf(fmt.Sprintf("Using '%s' for undocumented argument type '%s'", function.Arguments[argIndex].Type, function.Arguments[argIndex].Name))
			continue
		}

		argSymbol, err := tr.getType(paramAnnotation.Type())

		if err != nil {
			return nil, err
		}

		// extracting type references only if they are not generics
		for _, typeReference := range tr.extractTypeReferences(argSymbol) {
			if _, isGenericTypeReference := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
				return typeReference.Value == generic.Name
			}); !isGenericTypeReference {
				err := tr.context.Extract("type", typeReference.Value)

				if err != nil {
					return nil, err
				}
			}
		}

		function.Arguments[argIndex].Type = argSymbol
		function.Arguments[argIndex].Optional = paramAnnotation.Optional()
	}

	for _, returnAnnotation := range annotations.AtReturns {
		functionReturn := symbol.NewFunctionReturn()
		functionReturn.Name = returnAnnotation.Name()

		returnSymbol, err := tr.getType(returnAnnotation.Type())

		if err != nil {
			return nil, err
		}

		// extracting type references only if they are not generics
		for _, typeReference := range tr.extractTypeReferences(returnSymbol) {
			if _, isGenericTypeReference := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
				return typeReference.Value == generic.Name
			}); !isGenericTypeReference {
				err := tr.context.Extract("type", typeReference.Value)

				if err != nil {
					return nil, err
				}
			}
		}

		functionReturn.Type = returnSymbol
		function.Returns = append(function.Returns, *functionReturn)
	}

	for _, overloadAnnotation := range annotations.AtOverloads {
		overloadType, err := tr.transformType(overloadAnnotation.Type())

		if err != nil {
			return nil, fmt.Errorf("Error lexing function overload annotation type: %w", err)
		}

		overloadFunctionType, isFunctionType := overloadType.(*symbol.Function)

		if !isFunctionType {
			return nil, fmt.Errorf("Error lexing function overload annotation type: lexed type '%+v' is not a function", overloadFunctionType)
		}

		functionOverload := symbol.NewFunctionOverload()
		functionOverload.Generics = overloadFunctionType.Generics
		functionOverload.Arguments = overloadFunctionType.Arguments
		functionOverload.Returns = overloadFunctionType.Returns
		function.Overloads = append(function.Overloads, *functionOverload)
	}

	return function, nil
}
