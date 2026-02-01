package transform

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type AtGenericAnnotation struct {
	Name string
	Type *TypeAnnotation
}

func (tr *Transformer) getAtGenericAnnotations(buffer *nvim.ScratchBuffer) ([]AtGenericAnnotation, error) {
	atGenericAnnotations := []AtGenericAnnotation{}
	matches, err := buffer.TsQueryAll(annotation.AtGenericsQuery)

	if err != nil {
		return nil, err
	}

	if matches == nil {
		return atGenericAnnotations, nil
	}

	for _, matchCaptures := range *matches {
		atGenenericsAnnotation := annotation.NewGenerics(matchCaptures)

		for _, generic := range atGenenericsAnnotation.Generics() {
			lexedGeneric := AtGenericAnnotation{
				Name: generic.Name(),
			}

			if generic.Type() != nil {
				lexedGeneric.Type = &TypeAnnotation{*generic.Type()}
			}

			atGenericAnnotations = append(atGenericAnnotations, lexedGeneric)
		}

	}

	return atGenericAnnotations, nil
}

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

type AtParamAnnotation struct {
	Name          string
	Type          TypeAnnotation
	Optional      bool
	Documentation []string
}

func (tr *Transformer) getAtParamAnnotations(buffer *nvim.ScratchBuffer) (map[string]AtParamAnnotation, error) {
	params := map[string]AtParamAnnotation{}
	matches, err := buffer.TsQueryAll(annotation.AtParamQuery)

	if err != nil {
		return nil, err
	}

	if matches == nil {
		return params, nil
	}

	for _, matchCaptures := range *matches {
		atParamAnnotation := annotation.NewAtParam(matchCaptures)
		name := atParamAnnotation.Name()
		optional := atParamAnnotation.Optional()
		type_ := TypeAnnotation{atParamAnnotation.Type()}

		params[name] = AtParamAnnotation{
			Name:     name,
			Optional: optional,
			Type:     type_,
		}
	}

	return params, nil
}

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

	// TODO: remove
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
		generic := symbol.NewFunctionGeneric(genericAnnotation.Name, nil)

		if genericAnnotation.Type != nil {
			genericType, err := tr.getType(*genericAnnotation.Type)

			if err != nil {
				return nil, err
			}

			generic.Type = genericType
		}

		function.Generics = append(function.Generics, *generic)
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
