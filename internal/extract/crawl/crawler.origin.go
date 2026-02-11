package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) getAnnotationOriginMap() nvim.TsNodeQueryMap {
	return nvim.TsNodeQueryMap{
		origin.Comment: []treesitter.Query{
			annotation.AtAliasQuery.MapQuery(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @alias.name "%s")
				) @origin.alias`, query, c.context.Target().Identifier())
			}),
			annotation.AtAliasEnumeratorQuery.MapQuery(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @alias.name "%s")
				) @origin.alias.enumerator`, query, c.context.Target().Identifier())
			}),
			annotation.AtClassQuery.MapQuery(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @class.name "%s")
				) @origin.class`, query, c.context.Target().Identifier())
			}),
			{
				Language: origin.FieldAnnotationQuery.Language,
				Query: fmt.Sprintf(`(
					(%s)
					(#eq? @field.name "%s")
				) @origin.fieldannotation`, origin.FieldAnnotationQuery.Query, c.context.Target().Name()),
			},
			annotation.AtEnumQuery.MapQuery(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @enum.name "%s")
				) @origin.enumannotation`, query, c.context.Target().Identifier())
			}),
		},
	}
}

func (c *Crawler) getOriginChain(locations []nvim.Location) (origin.OriginChain, error) {
	var locationOrigin origin.Origin
	var err error

	for _, location := range locations {
		c.context.Logger().Verbosef("Searching for annotation origin at: %+v:%d:%d", location.Url, location.StartLine(), location.StartCharacter())

		locationOrigin, err = c.getAnnotationOrigin(location)

		if err != nil {
			return nil, err
		}

		if locationOrigin != nil {
			break
		}

		c.context.Logger().Verbosef("Searching for definition origin at: %+v:%d:%d", location.Url, location.StartLine(), location.StartCharacter())

		locationOrigin, err = c.getDefinitionOrigin(location)

		if err != nil {
			return nil, err
		}

		if locationOrigin != nil {
			break
		}
	}

	if locationOrigin == nil {
		return origin.NewOriginChain(origin.NewUnkownOrigin()), nil
	}

	c.context.Logger().Verbosef("Origin found <%T>", locationOrigin)
	return origin.NewOriginChain(locationOrigin), nil
}

func (c *Crawler) getAnnotationOrigin(location nvim.Location) (origin.Origin, error) {
	buffer, err := c.context.Nvim().OpenBuffer(location.Url)

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	originQueryMap := c.getAnnotationOriginMap()
	line, character :=
		uint(location.TargetRange.Start.Line),
		uint(location.TargetRange.Start.Character)
	queryResult, err := buffer.QueryTsNodeAt(originQueryMap, line, character)

	if err != nil {
		return nil, err
	}

	if queryResult == nil {
		return nil, nil
	}

	originMatch := queryResult.Match

	if _, isClass := originMatch.Find("origin.class"); isClass {
		node := annotation.NewAtClass(originMatch)
		return origin.NewClassOrigin(location, *node), nil
	}

	if _, isAlias := originMatch.Find("origin.alias"); isAlias {
		node := annotation.NewAtAlias(originMatch)
		return origin.NewAliasOrigin(location, *node), nil
	}

	if _, isAliasEnumerator := originMatch.Find("origin.alias.enumerator"); isAliasEnumerator {
		node := annotation.NewAtAliasEnumerator(originMatch)
		nextLines, err := buffer.NextLineIterator(uint(originMatch.LineRange().Start + 1))

		if err != nil {
			return nil, err
		}

		memberNodes := []annotation.AtAliasEnumeratorMember{}

		for line, err := range nextLines {
			if err != nil {
				return nil, err
			}

			match, err := line.TsQueryOne(annotation.AtAliasEnumeratorMemberQuery)

			if err != nil {
				return nil, err
			}

			if match == nil {
				break
			}

			memberNodes = append(memberNodes, *annotation.NewAtAliasEnumeratorMember(*match))
		}

		return origin.NewAliasEnumeratorOrigin(location, *node, memberNodes...), nil
	}

	if _, isFieldAnnotation := originMatch.Find("origin.fieldannotation"); isFieldAnnotation {
		return origin.NewFieldAnnotationOrigin(location, originMatch), nil
	}

	if _, isEnumAnnotation := originMatch.Find("origin.enumannotation"); isEnumAnnotation {
		node := annotation.NewAtEnum(originMatch)
		enumMemberNodes := []annotation.AtEnumMember{}

		position := originMatch.Range().Start
		dockblock, err := buffer.GetTsCommentBlockAt(uint(position.Line), uint(position.Character))

		if err != nil {
			return nil, err
		}

		membersPosition := treesitter.LineRange{
			Start: dockblock.Range.End.Line + 1,
			End:   dockblock.Range.End.Line + 1,
		}
		membersMatch, err := buffer.TsQueryOne(annotation.AtEnumMembersQuery.Ranged(membersPosition))

		if err != nil {
			return nil, err
		}

		if membersMatch == nil {
			return nil, fmt.Errorf("Could not find members definition for enumerator annotation origin found at <%+v>", location)
		}

		membersRange := membersMatch.LineRange()
		memberMatches, err := buffer.TsQueryAll(annotation.AtEnumMemberQuery.Ranged(*membersRange))

		if err != nil {
			return nil, err
		}

		if memberMatches == nil {
			return nil, fmt.Errorf("Could not find any member definition for enumerator annotation origin <%v> found at <%+v>", memberMatches, location)
		}

		for _, memberMatch := range *memberMatches {
			enumMemberNodes = append(enumMemberNodes, *annotation.NewAtEnumMember(memberMatch))
		}
		return origin.NewEnumAnnotationOrigin(location, *node, enumMemberNodes...), nil
	}

	return nil, nil
}

func (c *Crawler) getDefinitionOriginQueryMap() nvim.TsNodeQueryMap {
	return nvim.TsNodeQueryMap{
		origin.AssignmentStatement: []treesitter.Query{
			definition.MetatableAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.table", query)
			}),
			definition.TableFieldAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.table", query)
			}),
			definition.TableFieldIndexAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.table", query)
			}),
			definition.VirtualVariableAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.meta", query)
			}),
			definition.FunctionFieldDotAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.function", query)
			}),
			definition.FunctionFieldIndexAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.function", query)
			}),
			definition.FunctionCallAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.function_call", query)
			}),
			definition.RequireFunctionCallAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.module", query)
			}),
			definition.VariableDotFieldAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.variable", query)
			}),
			definition.VariableAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.variable", query)
			}),
			definition.TableDeclarationQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.table", query)
			}),
			definition.TableReturnAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.table", query)
			}),
		},
		origin.VariableDeclaration: []treesitter.Query{
			definition.FunctionVariableDeclarationQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.function", query)
			}),
		},
		origin.FunctionDeclaration: []treesitter.Query{
			definition.FunctionDeclarationQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.function", query)
			}),
			definition.FunctionFieldMethodDeclarationQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.function", query)
			}),
			definition.FunctionFieldDotDeclarationQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.function", query)
			}),
		},
		origin.Field: []treesitter.Query{
			definition.VirtualFieldAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @virtual.name "%s")
				) @origin.meta`, query, c.context.Target().Name())
			}),
			definition.TableConstructorFieldAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @table.name "%s")
				) @origin.table`, query, c.context.Target().Name())
			}),
			definition.TableConstructorFieldIndexAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @table.name "%s")
				) @origin.table`, query, c.context.Target().Name())
			}),
			definition.ValueFieldAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @field.name "%s")
				) @origin.value`, query, c.context.Target().Name())
			}),
			definition.ValueFieldIndexAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @field.name "%s")
				) @origin.value`, query, c.context.Target().Name())
			}),
		},
	}
}

func (c *Crawler) getDefinitionOrigin(location nvim.Location) (origin.Origin, error) {
	url, line, character :=
		location.Url,
		uint(location.TargetRange.Start.Line),
		uint(location.TargetRange.Start.Character)

	buffer, err := c.context.Nvim().OpenBuffer(url)

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	queryMap := c.getDefinitionOriginQueryMap()
	queryResult, err := buffer.QueryTsNodeAt(queryMap, line, character)

	if err != nil {
		return nil, err
	}

	// c.context.Logger().Debugf("Definition match: %+v", queryResult)

	if queryResult == nil {
		return nil, nil
	}

	match := queryResult.Match

	if _, isFunction := match.Find("origin.function"); isFunction {
		node := definition.NewFunction(match)
		annotations, err := buffer.GetTsNodeAnnotations(node.Root())

		if err != nil {
			return nil, err
		}

		return origin.NewFunctionOrigin(location, *node, annotations), nil
	}

	if _, isTable := match.Find("origin.table"); isTable {
		node := definition.NewTable(match)
		annotations, err := buffer.GetTsNodeAnnotations(node.Root())

		if err != nil {
			return nil, err
		}

		return origin.NewTableOrigin(location, *node, annotations), nil
	}

	if _, isValue := match.Find("origin.value"); isValue {
		node := definition.NewValue(match)
		annotations, err := buffer.GetTsNodeAnnotations(node.Root())

		if err != nil {
			return nil, err
		}

		return origin.NewValueOrigin(location, *node, annotations), nil
	}

	if _, isVariable := match.Find("origin.variable"); isVariable {
		node := definition.NewVariable(match)
		annotations, err := buffer.GetTsNodeAnnotations(node.Root())

		if err != nil {
			return nil, err
		}

		return origin.NewVariableOrigin(location, *node, annotations), nil
	}

	if _, isFunctionCall := match.Find("origin.function_call"); isFunctionCall {
		node := definition.NewFunctionCall(match)
		annotations, err := buffer.GetTsNodeAnnotations(node.Root())

		if err != nil {
			return nil, err
		}

		return origin.NewFunctionCallOrigin(location, *node, annotations), nil
	}

	if _, isRequireFunctionCall := match.Find("origin.module"); isRequireFunctionCall {
		node := definition.NewRequireFunctionCall(match)
		annotations, err := buffer.GetTsNodeAnnotations(node.Root())

		if err != nil {
			return nil, err
		}

		return origin.NewRequireFunctionCallOrigin(location, *node, annotations), nil
	}

	if _, isVirtual := match.Find("origin.meta"); isVirtual {
		node := definition.NewVirtual(match)
		annotations, err := buffer.GetTsNodeAnnotations(node.Root())

		if err != nil {
			return nil, err
		}

		return origin.NewVirtualOrigin(location, *node, annotations), nil
	}

	return nil, nil
}
