package lex

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type functionAtAnnotations struct {
	AtGenerics  []AtGenericAnnotation
	AtParams    map[string]AtParamAnnotation
	AtReturns   []AtReturnAnnotation
	AtOverloads []AtOverloadAnnotation
}

func (l *Lexer) lexFunctionAtAnnotations(docblock []string) (functionAtAnnotations, error) {
	annotations := functionAtAnnotations{
		AtGenerics:  []AtGenericAnnotation{},
		AtParams:    map[string]AtParamAnnotation{},
		AtReturns:   []AtReturnAnnotation{},
		AtOverloads: []AtOverloadAnnotation{},
	}

	buffer, err := l.context.Nvim().NewBuffer()

	if err != nil {
		return annotations, err
	}

	defer buffer.Close()

	err = buffer.SetLines(docblock)

	if err != nil {
		return annotations, err
	}

	atGenerics, err := l.lexAtGenericAnnotations(buffer)

	if err != nil {
		return annotations, err
	}

	atParams, err := l.lexAtParamAnnotations(buffer)

	if err != nil {
		return annotations, err
	}

	atReturns, err := l.lexAtReturnAnnotations(buffer)

	if err != nil {
		return annotations, err
	}

	atOverloads, err := l.lexAtOverloadAnnotations(buffer)

	if err != nil {
		return annotations, err
	}

	annotations.AtGenerics = atGenerics
	annotations.AtParams = atParams
	annotations.AtReturns = atReturns
	annotations.AtOverloads = atOverloads
	return annotations, nil
}

func (l *Lexer) lexFunctionValue(origin *symbol.FunctionOrigin) (*symbol.Function, error) {
	function := symbol.NewFunction()

	for _, capture := range origin.Captures() {
		switch capture.Id {
		case "name":
			function.Name = &capture.Node.Text
		case "access.class":
			function.Access = &symbol.FunctionClassAccess
		case "access.instance":
			function.Access = &symbol.FunctionIstanceAccess
		case "arg":
			function.Arguments = append(function.Arguments, *symbol.NewFunctionArgument(capture.Node.Text))
		case "vararg":
			function.Arguments = append(function.Arguments, *symbol.NewFunctionArgument(capture.Node.Text))
		}
	}

	annotations, err := l.lexFunctionAtAnnotations(origin.Documentation())

	if err != nil {
		return nil, fmt.Errorf("Error lexing function %s: %w", *function.Name, err)
	}

	for _, genericAnnotation := range annotations.AtGenerics {
		genericName := genericAnnotation.Name
		genericTypes, err := slicesx.MapFunc(genericAnnotation.Types, func(genericType TypeAnnotation) (symbol.Type, error) {
			return l.lexTypeAnnotation(genericType)
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
			l.context.Logger().Info(fmt.Sprintf("Using '%s' for undocumented argument type '%s'", function.Arguments[argIndex].Type, function.Arguments[argIndex].Name))
			continue
		}

		if functionGeneric, isGeneric := slicesx.FindFunc(function.Generics, func(generic symbol.FunctionGeneric) bool {
			return generic.Name == paramAnnotation.Type.Name
		}); isGeneric {
			function.Arguments[argIndex].Type = symbol.NewReference(functionGeneric.Name)
		} else {
			argumentType, err := l.lexTypeAnnotation(paramAnnotation.Type)

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
			typ, err := l.lexTypeAnnotation(returnAnnotation.Type)

			if err != nil {
				return nil, err
			}

			functionReturn.Type = typ
		}

		function.Returns = append(function.Returns, *functionReturn)
	}

	for _, overloadAnnotation := range annotations.AtOverloads {
		overloadType, err := l.lexTypeAnnotation(overloadAnnotation.Type)

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
