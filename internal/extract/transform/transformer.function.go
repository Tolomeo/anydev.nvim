package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type functionAtAnnotations struct {
	AtGenerics  []annotation.AtGenerics
	AtParams    map[string]annotation.AtParam
	AtReturns   []annotation.AtReturn
	AtOverloads []annotation.AtOverload
}

func (tr *Transformer) getFunctionAtAnnotations(docblock []string) (functionAtAnnotations, error) {
	annotations := functionAtAnnotations{
		AtGenerics:  []annotation.AtGenerics{},
		AtParams:    map[string]annotation.AtParam{},
		AtReturns:   []annotation.AtReturn{},
		AtOverloads: []annotation.AtOverload{},
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

	atGenericMatches, err := buffer.SafeTsQueryAll(annotation.AtGenericsQuery)

	if err != nil {
		return annotations, err
	}

	if atGenericMatches != nil {
		for _, match := range *atGenericMatches {
			if match.HasError {
				tr.context.Logger().Errorf("Skipping @generic annotation <%v> because it contains syntax errors", match.Captures)
				continue
			}

			annotations.AtGenerics = append(annotations.AtGenerics, *annotation.NewGenerics(match.Captures))
		}
	}

	atParamMatches, err := buffer.SafeTsQueryAll(annotation.AtParamQuery)

	if err != nil {
		return annotations, err
	}

	if atParamMatches != nil {
		for _, match := range *atParamMatches {
			if match.HasError {
				tr.context.Logger().Errorf("Skipping @param annotation <%v> because it contains syntax errors", match.Captures)
				continue
			}

			atParamAnnotation := annotation.NewAtParam(match.Captures)
			annotations.AtParams[atParamAnnotation.Name()] = *atParamAnnotation
		}
	}

	atReturnMatches, err := buffer.TsQueryAll(annotation.AtReturnQuery)

	if err != nil {
		return annotations, err
	}

	if atReturnMatches != nil {
		for _, match := range *atReturnMatches {
			annotations.AtReturns = append(annotations.AtReturns, *annotation.NewAtReturn(match))
		}
	}

	atOverloadMatches, err := buffer.SafeTsQueryAll(annotation.AtOverloadQuery)

	if err != nil {
		return annotations, err
	}

	if atOverloadMatches != nil {
		for _, match := range *atOverloadMatches {
			if match.HasError {
				tr.context.Logger().Errorf("Skipping @overload annotation <%v> because it contains syntax errors", match.Captures)
				continue
			}

			annotations.AtOverloads = append(annotations.AtOverloads, *annotation.NewOverload(match.Captures))
		}
	}

	return annotations, nil
}

func (tr *Transformer) getFunctionOriginType(functionOrigin *origin.FunctionOrigin) (*symbol.Function, error) {
	function := symbol.NewFunction()
	function.Name = functionOrigin.Name()

	for _, arg := range functionOrigin.Args() {
		function.Arguments = append(function.Arguments, *symbol.NewFunctionArgument(arg))
	}

	annotations, err := tr.getFunctionAtAnnotations(functionOrigin.Annotations())

	if err != nil {
		return nil, fmt.Errorf("Error lexing function %s: %w", function.Name, err)
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

		parsedArgType, err := tr.transformType(paramAnnotation.Type())

		if err != nil {
			return nil, err
		}

		switch reference := (parsedArgType.Canonical()).(type) {
		case *symbol.TypeReference:
			if _, isGenericArgType := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
				return reference.Value == generic.Name
			}); !isGenericArgType {
				err := tr.resolveTypeReferences(parsedArgType)
				if err != nil {
					return nil, err
				}
			}
		}

		function.Arguments[argIndex].Type = parsedArgType
		function.Arguments[argIndex].Optional = paramAnnotation.Optional()
	}

	for _, returnAnnotation := range annotations.AtReturns {
		functionReturn := symbol.NewFunctionReturn()
		functionReturn.Name = returnAnnotation.Name()

		parsedReturnType, err := tr.transformType(returnAnnotation.Type())

		if err != nil {
			return nil, err
		}

		switch reference := (parsedReturnType.Canonical()).(type) {
		case *symbol.TypeReference:
			if _, isGenericArgType := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
				return reference.Value == generic.Name
			}); !isGenericArgType {
				err := tr.resolveTypeReferences(reference)
				if err != nil {
					return nil, err
				}
			}
		}

		functionReturn.Type = parsedReturnType
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
