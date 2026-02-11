package crawl

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/definition"
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) getAnnotationOriginMap() nvim.TsNodeQueryMap {
	return nvim.TsNodeQueryMap{
		origin.Comment: []treesitter.Query{
			{
				Language: origin.AliasAnnotationQuery.Language,
				Query: fmt.Sprintf(`(
					(%s)
					(#eq? @alias.name "%s")
				) @origin.alias`, origin.AliasAnnotationQuery.Query, c.context.Target().Identifier()),
			},
			{
				Language: origin.AliasEnumeratorAnnotationQuery.Language,
				Query: fmt.Sprintf(`(
					(%s)
					(#eq? @alias.name "%s")
				) @origin.alias.enumerator`, origin.AliasEnumeratorAnnotationQuery.Query, c.context.Target().Identifier()),
			},
			{
				Language: origin.ClassAnnotationQuery.Language,
				Query: fmt.Sprintf(`(
					(%s)
					(#eq? @class.name "%s")
				) @origin.class`, origin.ClassAnnotationQuery.Query, c.context.Target().Identifier()),
			},
			{
				Language: origin.FieldAnnotationQuery.Language,
				Query: fmt.Sprintf(`(
					(%s)
					(#eq? @field.name "%s")
				) @origin.fieldannotation`, origin.FieldAnnotationQuery.Query, c.context.Target().Name()),
			},
			origin.EnumAnnotationQuery.MapQuery(func(query string) string {
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
		return origin.NewClassOrigin(location, originMatch), nil
	}

	if _, isAlias := originMatch.Find("origin.alias"); isAlias {
		return origin.NewAliasOrigin(location, originMatch), nil
	}

	if _, isAliasEnumerator := originMatch.Find("origin.alias.enumerator"); isAliasEnumerator {
		nextLines, err := buffer.NextLineIterator(uint(originMatch.LineRange().Start + 1))

		if err != nil {
			return nil, err
		}

		memberMatches := []nvim.TsQueryMatch{}

		for line, err := range nextLines {
			if err != nil {
				return nil, err
			}

			match, err := line.TsQueryOne(origin.AliasEnumeratorMemberAnnotationQuery)

			if err != nil {
				return nil, err
			}

			if match == nil {
				break
			}

			memberMatches = append(memberMatches, *match)
		}

		return origin.NewAliasEnumeratorOrigin(location, originMatch, memberMatches...), nil
	}

	if _, isFieldAnnotation := originMatch.Find("origin.fieldannotation"); isFieldAnnotation {
		return origin.NewFieldAnnotationOrigin(location, originMatch), nil
	}

	if _, isEnumAnnotation := originMatch.Find("origin.enumannotation"); isEnumAnnotation {
		enumMemberMatches := []nvim.TsQueryMatch{}

		position := originMatch.Range().Start
		dockblock, err := buffer.GetTsCommentBlockAt(uint(position.Line), uint(position.Character))

		if err != nil {
			return nil, err
		}

		membersPosition := treesitter.LineRange{
			Start: dockblock.Range.End.Line + 1,
			End:   dockblock.Range.End.Line + 1,
		}
		membersMatch, err := buffer.TsQueryOne(origin.EnumAnnotationMembersQuery.Ranged(membersPosition))

		if err != nil {
			return nil, err
		}

		if membersMatch == nil {
			return nil, fmt.Errorf("Could not find members definition for enumerator annotation origin found at <%+v>", location)
		}

		membersRange := membersMatch.LineRange()
		memberMatches, err := buffer.TsQueryAll(origin.EnumAnnotationMemberQuery.Ranged(*membersRange))

		if err != nil {
			return nil, err
		}

		if memberMatches == nil {
			return nil, fmt.Errorf("Could not find any member definition for enumerator annotation origin <%v> found at <%+v>", memberMatches, location)
		}

		enumMemberMatches = append(enumMemberMatches, *memberMatches...)
		return origin.NewEnumAnnotationOrigin(location, originMatch, enumMemberMatches...), nil
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
			origin.FunctionCallAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf("(%s) @origin.function_call", query)
			}),
			{
				Language: origin.RequireFunctionCallAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.module", origin.RequireFunctionCallAssignmentQuery.Query),
			},
			{
				Language: origin.VariableDotFieldAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.variable", origin.VariableDotFieldAssignmentQuery.Query),
			},
			{
				Language: origin.VariableAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.variable", origin.VariableAssignmentQuery.Query),
			},
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
			origin.ValueFieldAssignmentQuery.MapQuery(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @field.name "%s")
				) @origin.value`, query, c.context.Target().Name())
			}),
			origin.ValueFieldIndexAssignmentQuery.MapQuery(func(query string) string {
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

	originMatch := queryResult.Match
	originAnnotations, err := c.getDefinitionAnnotations(location)

	if err != nil {
		return nil, err
	}

	if _, isFunction := originMatch.Find("origin.function"); isFunction {
		node := definition.NewFunction(originMatch)
		return origin.NewFunctionOrigin(location, *node, originAnnotations), nil
	}

	if _, isTable := originMatch.Find("origin.table"); isTable {
		node := definition.NewTable(originMatch)
		return origin.NewTableOrigin(location, *node, originAnnotations), nil
	}

	if _, isValue := originMatch.Find("origin.value"); isValue {
		return origin.NewValueOrigin(location, originMatch, originAnnotations), nil
	}

	if _, isVariable := originMatch.Find("origin.variable"); isVariable {
		return origin.NewVariableOrigin(location, originMatch, originAnnotations), nil
	}

	if _, isFunctionCall := originMatch.Find("origin.function_call"); isFunctionCall {
		return origin.NewFunctionCallOrigin(location, originMatch, originAnnotations), nil
	}

	if _, isModuleRequire := originMatch.Find("origin.module"); isModuleRequire {
		return origin.NewRequireFunctionCallOrigin(location, originMatch, originAnnotations), nil
	}

	if _, isVirtual := originMatch.Find("origin.meta"); isVirtual {
		node := definition.NewVirtual(originMatch)
		annotations, err := buffer.GetTsNodeAnnotations(node.Root())

		if err != nil {
			return nil, err
		}

		return origin.NewVirtualOrigin(location, *node, annotations), nil
	}

	return nil, nil
}

func (c *Crawler) getDefinitionAnnotations(location nvim.Location) ([]string, error) {
	buffer, err := c.context.Nvim().OpenBuffer(location.Url)

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	annotations, err := buffer.GetTsCommentBlockAt(location.StartLine()-1, location.StartCharacter())

	if err != nil {
		return nil, err
	}

	if annotations == nil {
		c.context.Logger().Warnf("No documentation found for location <%v>", location)
		return []string{}, nil
	}

	return strings.Split(annotations.Text, "\n"), nil
}
