package transform

import (
	"errors"
	"reflect"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

func IsEmpty(s any) (bool, error) {
	val := reflect.ValueOf(s)

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return false, errors.New("input must be a struct or a pointer to a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		switch field.Kind() {
		case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Chan, reflect.Func:
			if !field.IsNil() {
				return false, nil
			}
		}
	}

	return true, nil
}

type functionAtAnnotations struct {
	AtGenerics  []annotation.AtGenerics
	AtParams    map[string]annotation.AtParam
	AtReturns   []annotation.AtReturn
	AtOverloads []annotation.AtOverload
}

func (tr *Transformer) getFunctionAtAnnotations(docblock []string) (*functionAtAnnotations, error) {
	annotations := functionAtAnnotations{}

	buffer, err := tr.context.Nvim().OpenScratchBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	err = buffer.SetLines(docblock)

	if err != nil {
		return nil, err
	}

	atGenericMatches, err := buffer.SafeTsQueryAll(annotation.AtGenericsQuery)

	if err != nil {
		return nil, err
	}

	if atGenericMatches != nil {
		annotations.AtGenerics = []annotation.AtGenerics{}

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
		return nil, err
	}

	if atParamMatches != nil {
		annotations.AtParams = map[string]annotation.AtParam{}

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
		return nil, err
	}

	if atReturnMatches != nil {
		annotations.AtReturns = []annotation.AtReturn{}

		for _, match := range *atReturnMatches {
			annotations.AtReturns = append(annotations.AtReturns, *annotation.NewAtReturn(match))
		}
	}

	atOverloadMatches, err := buffer.SafeTsQueryAll(annotation.AtOverloadQuery)

	if err != nil {
		return nil, err
	}

	if atOverloadMatches != nil {
		annotations.AtOverloads = []annotation.AtOverload{}

		for _, match := range *atOverloadMatches {
			if match.HasError {
				tr.context.Logger().Errorf("Skipping @overload annotation <%v> because it contains syntax errors", match.Captures)
				continue
			}

			annotations.AtOverloads = append(annotations.AtOverloads, *annotation.NewOverload(match.Captures))
		}
	}

	isEmpty, err := IsEmpty(annotations)

	if err != nil {
		return nil, err
	}

	if isEmpty {
		return nil, nil
	}

	return &annotations, nil
}

type typeAtAnnotations struct {
	AtType     *annotation.AtType
	AtClass    *annotation.AtClass
	AtOverload *annotation.AtOverload
	AtModule   *annotation.AtModule
}

func (tr *Transformer) getTypeAtAnnotations(buffer nvim.Buffer) (*typeAtAnnotations, error) {
	annotations := typeAtAnnotations{}

	atTypeMatch, err := buffer.SafeTsQueryOne(annotation.AtTypeQuery)

	if err != nil {
		return nil, err
	}

	if atTypeMatch != nil {
		annotations.AtType = annotation.NewAtType(atTypeMatch.Captures)
		return &annotations, nil
	}

	atClassMatch, err := buffer.TsQueryOne(annotation.AtClassQuery)

	if err != nil {
		return nil, err
	}

	if atClassMatch != nil {
		annotations.AtClass = annotation.NewAtClass(*atClassMatch)
		return &annotations, nil
	}

	atModuleMatch, err := buffer.SafeTsQueryOne(annotation.AtModuleQuery)

	if err != nil {
		return nil, err
	}

	if atModuleMatch != nil {
		annotations.AtModule = annotation.NewAtModule(atModuleMatch.Captures)
		/* if atModuleMatch.HasError {
			tr.context.Logger().Errorf("Skipping @module annotation <%v> because it contains syntax errors", atModuleMatch.Captures)
		} else {
		} */
		return &annotations, nil
	}

	atOverloadMatch, err := buffer.TsQueryOne(annotation.AtOverloadQuery)

	if err != nil {
		return nil, err
	}

	if atOverloadMatch != nil {
		annotations.AtOverload = annotation.NewOverload(*atOverloadMatch)
		return &annotations, nil
	}

	return nil, nil
}
