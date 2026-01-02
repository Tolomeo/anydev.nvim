package lex

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/cache"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

var lexedAnnotationsCache = cache.NewCache[*lexedAnnotations]()

type lexedAnnotations struct {
	tipe      symbol.Symbol
	private   bool
	protected bool
	params    map[string]symbol.FunctionArg
	overloads []symbol.FunctionOverload
	generics  []symbol.FunctionGeneric
	returns   []symbol.FunctionReturn
}

var typeAnnotationQuery string = fmt.Sprintf(`
	(documentation
		(type_annotation
			%s @type
		)
	) @typeannotation
`, anyTypeQuery)

func (l *Lexer) lexTypeAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	captures, err := buffer.TsQueryOne(treesitter.Query{Language: "luadoc", Query: typeAnnotationQuery})

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

	lexedType, err := l.lexType(types[0])

	if err != nil {
		return false, err
	}

	annotations.tipe = lexedType
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

func (l *Lexer) lexOverloadAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, err := buffer.SafeTsQueryAll(treesitter.Query{Language: "luadoc", Query: overloadAnnotationQuery})

	// fmt.Printf("\n Overload matches: %+v\n", matches)

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, match := range *matches {
		if match.HasError {
			l.context.Logger.Warn(fmt.Sprintf("Skipping overload annotation in '%s' because it contains syntax errors", l.context.Current()))
			continue
		}

		overload := symbol.FunctionOverload{}

		for _, capture := range match.Captures {
			switch capture.Id {
			case "documentation":
				overload.Documentation = []string{capture.Node.Text}
			case "type":
				overloadType, err := l.lexType(capture.Node.Text)

				if err != nil {
					return false, fmt.Errorf("Error overload annotation type: %w", err)
				}

				overloadFunctionType, isFunctionType := overloadType.(*symbol.Function)

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

func (l *Lexer) lexGenericAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(treesitter.Query{Language: "luadoc", Query: genericAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		lexedGeneric := symbol.FunctionGeneric{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "generic.name":
				lexedGeneric.Name = capture.Node.Text
			case "generic.type":
				genericType, err := l.lexType(capture.Node.Text)

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

func (l *Lexer) lexParamAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, hasMatches, err := slicesx.MapFindFunc(
		mapx.Values(paramAnnotationQueries),
		func(paramAnnotationQuery string) (*[]nvim.TsQueryMatch, bool, error) {
			paramMatches, err := buffer.TsQueryAll(treesitter.Query{Language: "luadoc", Query: paramAnnotationQuery})

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
		lexedParam := symbol.FunctionArg{}

		for _, matchCapture := range matchCaptures {
			switch matchCapture.Id {
			case "name":
				lexedParam.Name = matchCapture.Node.Text
			case "optional":
				lexedParam.Optional = true
			case "documentation":
				lexedParam.Documentation = []string{matchCapture.Node.Text}
			case "type":
				if generic, isGeneric := slicesx.FindFunc(annotations.generics, func(generic symbol.FunctionGeneric) bool {
					return generic.Name == matchCapture.Node.Text
				}); isGeneric {
					lexedParam.Type = symbol.NewReference(generic.Name)
					continue
				}

				lexedParamType, err := l.lexType(matchCapture.Node.Text)

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

func (l *Lexer) lexReturnAnnotations(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	matches, err := buffer.TsQueryAll(treesitter.Query{Language: "luadoc", Query: returnAnnotationQuery})

	switch {
	case err != nil:
		return false, err
	case matches == nil:
		return false, nil
	}

	for _, matchCaptures := range *matches {
		functionReturn := symbol.FunctionReturn{}

		for _, capture := range matchCaptures {
			switch capture.Id {
			case "return.name":
				functionReturn.Name = &capture.Node.Text
			case "return.documentation":
				functionReturn.Documentation = []string{capture.Node.Text}
			case "return.type":
				if generic, isGeneric := slicesx.FindFunc(annotations.generics, func(generic symbol.FunctionGeneric) bool {
					return generic.Name == capture.Node.Text
				}); isGeneric {
					functionReturn.Type = symbol.NewReference(generic.Name)
					continue
				}

				lexedReturnType, err := l.lexType(capture.Node.Text)

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

func (l *Lexer) lexPrivateAnnotation(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	captures, err := buffer.TsQueryOne(treesitter.Query{Language: "luadoc", Query: privateAnnotationQuery})

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

func (l *Lexer) lexProtectedAnnotation(buffer *nvim.Buffer, annotations *lexedAnnotations) (bool, error) {
	captures, err := buffer.TsQueryOne(treesitter.Query{Language: "luadoc", Query: protectedAnnotationQuery})

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

	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	annotations := lexedAnnotations{
		generics:  []symbol.FunctionGeneric{},
		params:    make(map[string]symbol.FunctionArg),
		overloads: []symbol.FunctionOverload{},
		returns:   []symbol.FunctionReturn{},
	}

	err = buffer.SetLines(dockblock)

	if err != nil {
		return nil, err
	}

	// Generics are lexed ahead of other annotations, which could read them
	_, err = l.lexGenericAnnotations(buffer, &annotations)

	if err != nil {
		return nil, fmt.Errorf("Error lexing generic annotations: %w", err)
	}

	_, err = l.lexPrivateAnnotation(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing private annotation: %w", err)
	}

	_, err = l.lexProtectedAnnotation(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing private annotation: %w", err)
	}

	_, err = l.lexParamAnnotations(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing param annotation: %w", err)
	}

	_, err = l.lexOverloadAnnotations(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing overload annotation: %w", err)
	}

	_, err = l.lexReturnAnnotations(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing return annotation: %w", err)
	}

	_, err = l.lexTypeAnnotations(buffer, &annotations)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error lexing type annotation: %w", err)
	}

	lexedAnnotationsCache.Set(&annotations, dockblock...)
	return &annotations, nil
}

// Luadoc matches an empty type node even when the type is not present
// So those false positives are excluded with the not-eq predicate
var simpleAliasQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(alias_annotation
		(identifier) @alias.name
		%s @alias.type
		(comment)? @alias.documentation
		(#not-eq? @alias.type "")
	)`, anyTypeQuery),
}

func (l *Lexer) lexSimpleAlias(buffer *nvim.Buffer, source *symbol.TypeSource) (symbol.Symbol, error) {
	match, err := buffer.TsQueryOne(simpleAliasQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	typeCapture, typeCaptureFound := slicesx.FindFunc(*match, func(capture treesitter.Capture) bool {
		return capture.Id == "alias.type"
	})

	if !typeCaptureFound {
		return nil, fmt.Errorf("Error retrieving type value from alias annotation '%s'", source.Origin.DocumentationText())
	}

	lexedAliasType, err := l.lexType(typeCapture.Node.Text)

	if err != nil {
		return nil, err
	}

	return lexedAliasType, nil
}

var enumAliasQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(continuation
		%s @enumAlias.type
	)
`, anyTypeQuery)}

func (l *Lexer) lexEnumAlias(buffer *nvim.Buffer, source *symbol.TypeSource) (*symbol.Union, error) {
	matches, err := buffer.TsQueryAll(enumAliasQuery)

	switch {
	case err != nil:
		return nil, err
	case matches == nil:
		return nil, nil
	}

	enumMembers := []symbol.Symbol{}

	for _, matchCaptures := range *matches {
		for _, capture := range matchCaptures {
			switch capture.Id {
			case "enumAlias.type":
				lexedType, err := l.lexType(capture.Node.Text)

				if err != nil {
					return nil, err
				}

				enumMembers = append(enumMembers, lexedType)
			}
		}
	}

	// fmt.Println(source.Path)
	/* for _, member := range enumMembers {
		fmt.Println(member)

	} */

	if len(enumMembers) < 1 {
		return nil, fmt.Errorf("Could not retrieve enum members from enum alias '%s'", source.Path)
	}

	return newUnionType(enumMembers), nil
}

func (l *Lexer) lexAlias(source *symbol.TypeSource) (symbol.Symbol, error) {
	// fmt.Println(source.Path)

	if source.Origin.Type() != "alias_annotation" {
		return nil, nil
	}

	buffer, err := l.context.Nvim.NewBuffer()

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	// Replacing all dots in the alias name with underscores
	// because apparently luadoc would not permit to use dots in identifiers
	origin := source.Origin
	name := source.Path
	definitionText := origin.DefinitionText()
	patchedDefinitionText := strings.Replace(definitionText, name, strings.ReplaceAll(name, ".", "_"), 1)
	documentationText := origin.DocumentationText()
	patchedDocumentationLines := strings.Split(
		strings.Replace(documentationText, definitionText, patchedDefinitionText, 1),
		"\n",
	)

	buffer.SetLines(patchedDocumentationLines)

	simpleAliasType, err := l.lexSimpleAlias(buffer, source)

	switch {
	case err != nil:
		return nil, err
	case simpleAliasType != nil:
		return simpleAliasType, nil
	}

	enumAliasType, err := l.lexEnumAlias(buffer, source)

	switch {
	case err != nil:
		return nil, err
	case enumAliasType != nil:
		return enumAliasType, nil
	}

	return nil, fmt.Errorf("Could not lex alias type '%s'", source.Path)
}
