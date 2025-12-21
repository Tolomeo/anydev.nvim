package lex

import (
	"fmt"

	// "github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/ts"
	"github.com/Tolomeo/anydev.nvim/internal/utils/cache"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

var lexedAnnotationsCache = cache.NewCache[*lexedAnnotations]()

type lexedAnnotations struct {
	Type      lexed.Symbol
	private   bool
	protected bool
	params    map[string]lexed.FunctionArg
	overloads []lexed.FunctionOverload
	generics  []lexed.FunctionGeneric
	returns   []lexed.FunctionReturn
}

var typeAnnotationQuery string = fmt.Sprintf(`
	(documentation
		(type_annotation
			%s @type
		)
	) @typeannotation
`, anyTypeQuery)

func (l *Lexer) lexTypeAnnotations(annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQueryOne(nvim.TsQueryConfig{Language: "luadoc", Query: typeAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	types := []string{}

	for _, capture := range *captures {
		switch capture.Id {
		case "type":
			types = append(types, capture.Node.Text)
		}
	}

	switch {
	case len(types) > 1:
		return false, fmt.Errorf("Error reading type annotation: too many annotations received (%d), expected 1", len(types))
	case len(types) < 1:
		return false, fmt.Errorf("Error reading type annotation: type value not found")
	}

	var lexedType lexed.Symbol

	err = l.temp(func() error {
		lexed, err := l.lexType(types[0])

		if err != nil {
			return err
		}

		lexedType = lexed
		return nil
	})

	if err != nil {
		return false, err
	}

	annotations.Type = lexedType
	return true, nil
}

var overloadAnnotationQuery string = fmt.Sprintf(`
	(documentation 
		(overload_annotation 
			%s @type 
			(comment)? @documentation
		)
	)
`, typeQueries["function_type"])

func (l *Lexer) lexOverloadAnnotations(annotations *lexedAnnotations) (bool, error) {
	matches, err := l.context.nvim.TsQueryAll(nvim.TsQueryConfig{Language: "luadoc", Query: overloadAnnotationQuery})

	// fmt.Printf("\n Overload matches: %+v\n", matches)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		overload := lexed.FunctionOverload{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "documentation":
				overload.Documentation = []string{capture.Node.Text}
			case "type":
				var overloadType lexed.Symbol

				err = l.temp(func() error {
					lexedType, err := l.lexType(capture.Node.Text)

					if err != nil {
						return err
					}

					overloadType = lexedType
					return nil
				})

				if err != nil {
					return false, fmt.Errorf("Error overload annotation type: %w", err)
				}

				overloadFunctionType, isFunctionType := overloadType.(*lexed.Function)

				if !isFunctionType {
					return false, fmt.Errorf("Error lexing overload annotation type: lexed type '%+v' is not a function", overloadFunctionType)
				}

				overload.Generics = overloadFunctionType.Generics
				overload.Args = overloadFunctionType.Args
				overload.Documentation = overloadFunctionType.Documentation
				overload.Return = overloadFunctionType.Return
			}
		}

		annotations.overloads = append(annotations.overloads, overload)
	}

	return true, nil
}

var genericAnnotationQuery string = fmt.Sprintf(`
	(documentation 
		(generic_annotation
			(identifier) @generic.name
			parent_type: 
				%s? @generic.type
		)
	) @generic
`, anyTypeQuery)

func (l *Lexer) lexGenericAnnotations(annotations *lexedAnnotations) (bool, error) {
	matches, err := l.context.nvim.TsQueryAll(nvim.TsQueryConfig{Language: "luadoc", Query: genericAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		lexedGeneric := lexed.FunctionGeneric{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "generic.name":
				lexedGeneric.Name = capture.Node.Text
			case "generic.type":
				var genericType lexed.Symbol

				err = l.temp(func() error {
					lexedType, err := l.lexType(capture.Node.Text)

					if err != nil {
						return err
					}

					genericType = lexedType
					return nil
				})

				if err != nil {
					return false, err
				}

				lexedGeneric.Types = append(lexedGeneric.Types, genericType)
			}
		}

		if lexedGeneric.Name == "" {
			return false, fmt.Errorf("Could not retrieve generic name for generic annotation '%v'", lexedGeneric)
		}

		annotations.generics = append(annotations.generics, lexedGeneric)
	}

	return true, nil
}

var paramAnnotationQueries = map[string]string{
	"arg": fmt.Sprintf(`
		(documentation
			(param_annotation
				(identifier) @name
				"?"? @optional
				%s @type
				(comment)? @documentation
			)
		) @param
	`, anyTypeQuery),
	"vararg": fmt.Sprintf(`
		(documentation
			(param_annotation 
				"..." @name
				%s @type
				(comment)? @documentation
			)
		) @param
	`, anyTypeQuery),
}

func (l *Lexer) lexParamAnnotations(annotations *lexedAnnotations) (bool, error) {
	matches, hasMatches, err := slicesx.MapFindFunc(
		mapx.Values(paramAnnotationQueries),
		func(paramAnnotationQuery string) (*[][]ts.Capture, bool, error) {
			paramMatches, err := l.context.nvim.TsQueryAll(nvim.TsQueryConfig{Language: "luadoc", Query: paramAnnotationQuery})

			switch {
			case err != nil:
				return nil, false, err
			case paramMatches == nil:
				return nil, false, nil
			}

			return paramMatches, true, nil
		},
	)

	switch {
	case err != nil:
		return false, err
	case !hasMatches:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		lexedParam := lexed.FunctionArg{}

		for _, matchCapture := range matchCaptures {
			switch matchCapture.Id {
			case "name":
				lexedParam.Name = matchCapture.Node.Text
			case "optional":
				lexedParam.Optional = true
			case "documentation":
				lexedParam.Documentation = []string{matchCapture.Node.Text}
			case "type":
				if generic, isGeneric := slicesx.FindFunc(annotations.generics, func(generic lexed.FunctionGeneric) bool {
					return generic.Name == matchCapture.Node.Text
				}); isGeneric {
					lexedParam.Type = newReferenceType(generic.Name)
					continue
				}

				var lexedParamType lexed.Symbol

				err = l.temp(func() error {
					lexedType, err := l.lexType(matchCapture.Node.Text)

					if err != nil {
						return err
					}

					lexedParamType = lexedType
					return nil
				})

				if err != nil {
					return false, fmt.Errorf("Error lexing type annotations : %w", err)
				}

				lexedParam.Type = lexedParamType
			}
		}

		if lexedParam.Name == "" {
			return false, fmt.Errorf("Could not retrieve param name for param annotation")
		}

		annotations.params[lexedParam.Name] = lexedParam
	}

	return true, nil
}

var returnAnnotationQuery string = fmt.Sprintf(`
	(documentation
		(return_annotation
			%s @return.type
			(comment)? @return.documentation
		)
	) @return
`, anyTypeQuery)

func (l *Lexer) lexReturnAnnotations(annotations *lexedAnnotations) (bool, error) {
	matches, err := l.context.nvim.TsQueryAll(nvim.TsQueryConfig{Language: "luadoc", Query: returnAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		functionReturn := lexed.FunctionReturn{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "return.name":
				functionReturn.Name = &capture.Node.Text
			case "return.documentation":
				functionReturn.Documentation = []string{capture.Node.Text}
			case "return.type":
				if generic, isGeneric := slicesx.FindFunc(annotations.generics, func(generic lexed.FunctionGeneric) bool {
					return generic.Name == capture.Node.Text
				}); isGeneric {
					functionReturn.Type = newReferenceType(generic.Name)
					continue
				}

				var lexedReturnType lexed.Symbol

				err = l.temp(func() error {
					lexedType, err := l.lexType(capture.Node.Text)

					if err != nil {
						return err
					}

					lexedReturnType = lexedType
					return nil
				})

				if err != nil {
					return false, err
				}

				functionReturn.Type = lexedReturnType
			}
		}

		annotations.returns = append(annotations.returns, functionReturn)
	}

	return true, nil
}

var privateAnnotationQuery string = `
	(documentation 
		(qualifier_annotation) @qualifier
		(#match? @qualifier "\\@private")
	)
`

func (l *Lexer) lexPrivateAnnotation(annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQueryOne(nvim.TsQueryConfig{Language: "luadoc", Query: privateAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	annotations.private = true

	return true, nil
}

var protectedAnnotationQuery string = `
	(documentation 
		(qualifier_annotation) @qualifier
		(#match? @qualifier "\\@protected")
	)
`

func (l *Lexer) lexProtectedAnnotation(annotations *lexedAnnotations) (bool, error) {
	captures, err := l.context.nvim.TsQueryOne(nvim.TsQueryConfig{Language: "luadoc", Query: protectedAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case captures == nil:
		return false, nil
	}

	annotations.protected = true

	return true, nil
}

func (l *Lexer) lexAnnotations(dockblock []string) (*lexedAnnotations, error) {
	if cachedAnnotations, cached := lexedAnnotationsCache.Get(dockblock...); cached {
		return cachedAnnotations, nil
	}

	/* annotationLines, err := slicesx.FilterFunc(dockblock, func(docLine string) (bool, error) {
		err := l.scratch([]string{docLine})

		if err != nil {
			return false, err
		}

		captures, err := l.context.nvim.SafeTsQuery(nvim.TsQueryConfig{Language: "luadoc", Query: annotationQuery})

		switch {
		case errors.Is(nvim.ErrSafeTSQueryNoMatch, err):
			l.context.logger.Error(fmt.Sprintf("Skipping annotation line '%s' containing syntax errors", docLine))
			return false, nil
		case err != nil:
			return false, err
		case captures == nil:
			return false, nil
		}

		return true, nil
	})

	if err != nil {
		return nil, err
	} */

	annotations := lexedAnnotations{
		generics:  []lexed.FunctionGeneric{},
		params:    make(map[string]lexed.FunctionArg),
		overloads: []lexed.FunctionOverload{},
		returns:   []lexed.FunctionReturn{},
	}

	err := l.scratch(dockblock)

	if err != nil {
		return nil, err
	}

	// Generics are lexed ahead of other annotations, which could read them
	_, err = l.lexGenericAnnotations(&annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing generic annotations: %w", err)
	}

	_, err = l.lexPrivateAnnotation(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing private annotation: %w", err)
	}

	_, err = l.lexProtectedAnnotation(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing private annotation: %w", err)
	}

	_, err = l.lexParamAnnotations(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing param annotation: %w", err)
	}

	_, err = l.lexOverloadAnnotations(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing overload annotation: %w", err)
	}

	_, err = l.lexReturnAnnotations(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing return annotation: %w", err)
	}

	_, err = l.lexTypeAnnotations(&annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type annotation: %w", err)
	}

	lexedAnnotationsCache.Set(&annotations, dockblock...)
	return &annotations, nil
}
