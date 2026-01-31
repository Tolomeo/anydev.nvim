package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type AtReturnAnnotation struct {
	Name          string
	Type          TypeAnnotation
	Documentation []string
}

func (tr *Transformer) getAtReturnAnnotations(buffer *nvim.ScratchBuffer) ([]AtReturnAnnotation, error) {
	atReturns := []AtReturnAnnotation{}
	matches, err := buffer.TsQueryAll(annotation.AtReturnQuery)

	if err != nil {
		return atReturns, err
	}

	if matches == nil {
		return atReturns, nil
	}

	for _, matchCaptures := range *matches {
		atReturnAnnotation := annotation.NewAtReturn(matchCaptures)

		returnAnnotation := AtReturnAnnotation{
			Name: atReturnAnnotation.Name(),
			Type: TypeAnnotation{atReturnAnnotation.Type()},
		}

		atReturns = append(atReturns, returnAnnotation)
	}

	return atReturns, nil
}

type functionAtAnnotations struct {
	AtGenerics  []AtGenericAnnotation
	AtParams    map[string]AtParamAnnotation
	AtReturns   []AtReturnAnnotation
	AtOverloads []AtOverloadAnnotation
}

func (tr *Transformer) getFunctionAtAnnotations(docblock []string) (functionAtAnnotations, error) {
	annotations := functionAtAnnotations{
		AtGenerics:  []AtGenericAnnotation{},
		AtParams:    map[string]AtParamAnnotation{},
		AtReturns:   []AtReturnAnnotation{},
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

	atGenerics, err := tr.getAtGenericAnnotations(buffer)

	if err != nil {
		return annotations, err
	}

	atParams, err := tr.getAtParamAnnotations(buffer)

	if err != nil {
		return annotations, err
	}

	atReturns, err := tr.getAtReturnAnnotations(buffer)

	if err != nil {
		return annotations, err
	}

	atOverloads, err := tr.getAtOverloadAnnotations(buffer)

	if err != nil {
		return annotations, err
	}

	annotations.AtGenerics = atGenerics
	annotations.AtParams = atParams
	annotations.AtReturns = atReturns
	annotations.AtOverloads = atOverloads
	return annotations, nil
}

func (tr *Transformer) getFunctionOriginType(functionOrigin *origin.FunctionOrigin) (*symbol.Function, error) {
	function := symbol.NewFunction()

	function.Name = functionOrigin.Name()

	if functionOrigin.Static() {
		function.Access = &symbol.FunctionClassAccess
	} else {
		function.Access = &symbol.FunctionIstanceAccess
	}

	for _, arg := range functionOrigin.Args() {
		function.Arguments = append(function.Arguments, *symbol.NewFunctionArgument(arg))
	}

	annotations, err := tr.getFunctionAtAnnotations(functionOrigin.Annotations())

	if err != nil {
		return nil, fmt.Errorf("Error lexing function %s: %w", function.Name, err)
	}

	for _, genericAnnotation := range annotations.AtGenerics {
		genericName := genericAnnotation.Name
		genericTypes, err := slicesx.MapFunc(genericAnnotation.Types, func(genericType TypeAnnotation) (annotation.Type, error) {
			return tr.getType(genericType)
		})

		if err != nil {
			return nil, err
		}

		function.Generics = append(function.Generics, *symbol.NewFunctionGeneric(genericName, genericTypes...))
	}

	for argIndex := range function.Arguments {
		name := function.Arguments[argIndex].Name

		paramAnnotation, hasParamAnnotation := annotations.AtParams[name]

		if !hasParamAnnotation {
			tr.context.Logger().Info(fmt.Sprintf("Using '%s' for undocumented argument type '%s'", function.Arguments[argIndex].Type, function.Arguments[argIndex].Name))
			continue
		}

		if functionGeneric, isGeneric := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
			return generic.Name == paramAnnotation.Type.Name
		}); isGeneric {
			function.Arguments[argIndex].Type = symbol.NewReference(functionGeneric.Name)
		} else {
			argumentType, err := tr.getType(paramAnnotation.Type)

			if err != nil {
				return nil, err
			}

			function.Arguments[argIndex].Type = argumentType
		}

		function.Arguments[argIndex].Optional = paramAnnotation.Optional
		function.Arguments[argIndex].Documentation = paramAnnotation.Documentation
	}

	for _, returnAnnotation := range annotations.AtReturns {
		functionReturn := symbol.NewFunctionReturn()
		functionReturn.Name = returnAnnotation.Name
		functionReturn.Documentation = returnAnnotation.Documentation

		if functionGeneric, isGeneric := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
			return generic.Name == returnAnnotation.Type.Name
		}); isGeneric {
			functionReturn.Type = symbol.NewReference(functionGeneric.Name)
		} else {
			typ, err := tr.getType(returnAnnotation.Type)

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
