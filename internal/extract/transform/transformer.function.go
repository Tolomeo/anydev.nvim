package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type AtOverloadAnnotation struct {
	Type          TypeAnnotation
	Documentation []string
}

func (tr *Transformer) getAtOverloadAnnotations(buffer *nvim.ScratchBuffer) ([]AtOverloadAnnotation, error) {
	atOverloads := []AtOverloadAnnotation{}
	matches, err := buffer.SafeTsQueryAll(annotation.AtOverloadQuery)

	// fmt.Printf("\n Overload matches: %+v\n", matches)

	if err != nil {
		return atOverloads, err
	}

	if matches == nil {
		return atOverloads, nil
	}

	for _, match := range *matches {
		if match.HasError {
			tr.context.Logger().Warn(fmt.Sprintf("Skipping overload annotation in '%s' because it contains syntax errors", tr.context.Target().Name()))
			continue
		}

		atOverloadAnnotation := annotation.NewOverload(match.Captures)
		overload := AtOverloadAnnotation{}
		overload.Type = TypeAnnotation{atOverloadAnnotation.Type()}
		overload.Documentation = atOverloadAnnotation.Documentation()

		atOverloads = append(atOverloads, overload)
	}

	return atOverloads, nil
}

type functionAtAnnotations struct {
	AtGenerics  []annotation.AtGenerics
	AtParams    map[string]annotation.AtParam
	AtReturns   []annotation.AtReturn
	AtOverloads []AtOverloadAnnotation
}

func (tr *Transformer) getFunctionAtAnnotations(docblock []string) (functionAtAnnotations, error) {
	annotations := functionAtAnnotations{
		AtGenerics:  []annotation.AtGenerics{},
		AtParams:    map[string]annotation.AtParam{},
		AtReturns:   []annotation.AtReturn{},
		AtOverloads: []AtOverloadAnnotation{},
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

	atGenericMatches, err := buffer.TsQueryAll(annotation.AtGenericsQuery)

	if err != nil {
		return annotations, err
	}

	if atGenericMatches != nil {
		for _, match := range *atGenericMatches {
			annotations.AtGenerics = append(annotations.AtGenerics, *annotation.NewGenerics(match))
		}
	}

	atParamMatches, err := buffer.TsQueryAll(annotation.AtParamQuery)

	if err != nil {
		return annotations, err
	}

	if atParamMatches != nil {
		for _, match := range *atParamMatches {
			atParamAnnotation := annotation.NewAtParam(match)
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

	atOverloads, err := tr.getAtOverloadAnnotations(buffer)

	if err != nil {
		return annotations, err
	}

	annotations.AtOverloads = atOverloads
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
				genericType, err := tr.getType(TypeAnnotation{*genericAnnotation.Type()})

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
			tr.context.Logger().Info(fmt.Sprintf("Using '%s' for undocumented argument type '%s'", function.Arguments[argIndex].Type, function.Arguments[argIndex].Name))
			continue
		}

		if functionGeneric, isGeneric := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
			return generic.Name == paramAnnotation.Type()
		}); isGeneric {
			function.Arguments[argIndex].Type = symbol.NewReference(functionGeneric.Name)
		} else {
			argumentType, err := tr.getType(TypeAnnotation{paramAnnotation.Type()})

			if err != nil {
				return nil, err
			}

			function.Arguments[argIndex].Type = argumentType
		}

		function.Arguments[argIndex].Optional = paramAnnotation.Optional()
	}

	for _, returnAnnotation := range annotations.AtReturns {
		functionReturn := symbol.NewFunctionReturn()
		functionReturn.Name = returnAnnotation.Name()

		if functionGeneric, isGeneric := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
			return generic.Name == returnAnnotation.Type()
		}); isGeneric {
			functionReturn.Type = symbol.NewReference(functionGeneric.Name)
		} else {
			typ, err := tr.getType(TypeAnnotation{returnAnnotation.Type()})

			if err != nil {
				return nil, err
			}

			functionReturn.Type = typ
		}

		function.Returns = append(function.Returns, *functionReturn)
	}

	for _, overloadAnnotation := range annotations.AtOverloads {
		overloadType, err := tr.getType(overloadAnnotation.Type)

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
		// functionOverload.Documentation = overloadFunctionType.Documentation
		functionOverload.Returns = overloadFunctionType.Returns

		function.Overloads = append(function.Overloads, *functionOverload)
	}

	return function, nil
}
