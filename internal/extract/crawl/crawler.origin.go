package crawl

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) getTypeOriginMap() nvim.TsNodeQueryMap {
	return nvim.TsNodeQueryMap{
		origin.AliasAnnotation: []treesitter.Query{
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
		},
		origin.ClassAnnotation: []treesitter.Query{
			{
				Language: origin.ClassAnnotationQuery.Language,
				Query: fmt.Sprintf(`(
					(%s)
					(#eq? @class.name "%s")
				) @origin.class`, origin.ClassAnnotationQuery.Query, c.context.Target().Identifier()),
			},
		},
		origin.FieldAnnotation: []treesitter.Query{
			{
				Language: origin.FieldAnnotationQuery.Language,
				Query: fmt.Sprintf(`(
					(%s)
					(#eq? @field.name "%s")
				) @origin.fieldannotation`, origin.FieldAnnotationQuery.Query, c.context.Target().Name()),
			},
		},
		origin.EnumAnnotation: []treesitter.Query{
			origin.EnumAnnotationQuery.Extend(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @enum.name "%s")
				) @origin.enumannotation`, query, c.context.Target().Identifier())
			}),
		},
	}
}

func (c *Crawler) getValueOriginQueryMap() nvim.TsNodeQueryMap {
	return nvim.TsNodeQueryMap{
		origin.AssignmentStatement: []treesitter.Query{
			{
				Language: origin.TableFieldAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.table", origin.TableFieldAssignmentQuery.Query),
			},
			{
				Language: origin.TableFieldIndexAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.table", origin.TableFieldIndexAssignmentQuery.Query),
			},
			{
				Language: origin.VirtualVariableAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.meta", origin.VirtualVariableAssignmentQuery.Query),
			},
			{
				Language: origin.FunctionFieldDotAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", origin.FunctionFieldDotAssignmentQuery.Query),
			},
			{
				Language: origin.FunctionFieldIndexAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", origin.FunctionFieldIndexAssignmentQuery.Query),
			},
			{
				Language: origin.ModuleRequireAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.module", origin.ModuleRequireAssignmentQuery.Query),
			},
			{
				Language: origin.VariableDotFieldAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.variable", origin.VariableDotFieldAssignmentQuery.Query),
			},
			{
				Language: origin.VariableAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.variable", origin.VariableAssignmentQuery.Query),
			},
		},
		origin.VariableDeclaration: []treesitter.Query{
			{
				Language: origin.FunctionVariableDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", origin.FunctionVariableDeclarationQuery.Query),
			},
			{
				Language: origin.TableVariableDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.table", origin.TableVariableDeclarationQuery.Query),
			},
		},
		origin.FunctionDeclaration: []treesitter.Query{
			{
				Language: origin.FunctionDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", origin.FunctionDeclarationQuery.Query),
			},
			{
				Language: origin.FunctionFieldMethodDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", origin.FunctionFieldMethodDeclarationQuery.Query),
			},
			{
				Language: origin.FunctionFieldDotDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", origin.FunctionFieldDotDeclarationQuery.Query),
			},
		},
		origin.Field: []treesitter.Query{
			origin.TableConstructorFieldAssignmentQuery.Extend(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @table.name "%s")
				) @origin.table`, query, c.context.Target().Name())
			}),
			origin.TableConstructorFieldIndexAssignmentQuery.Extend(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @table.name "%s")
				) @origin.table`, query, c.context.Target().Name())
			}),
			origin.ValueFieldAssignmentQuery.Extend(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @field.name "%s")
				) @origin.value`, query, c.context.Target().Name())
			}),
			origin.ValueFieldIndexAssignmentQuery.Extend(func(query string) string {
				return fmt.Sprintf(`(
					(%s)
					(#eq? @field.name "%s")
				) @origin.value`, query, c.context.Target().Name())
			}),
		},
	}
}

func (c *Crawler) getOriginChain(locations []nvim.Location) (origin.OriginChain, error) {
	var locationOrigin origin.Origin
	var err error

	for _, location := range locations {
		locationOrigin, err = c.getAnnotationOrigin(location)

		if err != nil {
			return nil, err
		}

		if locationOrigin != nil {
			break
		}

		locationOrigin, err = c.getDefinitionOrigin(location)

		if err != nil {
			return nil, err
		}

		if locationOrigin != nil {
			break
		}
	}

	if locationOrigin == nil {
		return nil, nil
	}

	switch ot := locationOrigin.(type) {
	case *origin.ModuleOrigin:
		moduleOrigins, err := c.getModuleOriginChain(ot)
		if err != nil {
			return nil, err
		}
		return origin.NewOriginChain(locationOrigin).Append(moduleOrigins), nil
	case *origin.VariableOrigin:
		variableOrigins, err := c.getVariableOriginChain(ot)
		if err != nil {
			return nil, err
		}
		return origin.NewOriginChain(locationOrigin).Append(variableOrigins), nil
	}

	return origin.NewOriginChain(locationOrigin), nil
}

func (c *Crawler) getAnnotationOrigin(location nvim.Location) (origin.Origin, error) {
	buffer, err := c.context.Nvim().OpenBuffer(location.Url)

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	originQueryMap := c.getTypeOriginMap()
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

	// match := queryResult.Match

	if _, isClass := queryResult.Match.Find("origin.class"); isClass {
		return origin.NewClassOrigin(location, queryResult.Match), nil
	}

	if _, isAlias := queryResult.Match.Find("origin.alias"); isAlias {
		return origin.NewAliasOrigin(location, queryResult.Match), nil
	}

	if _, isAliasEnumerator := queryResult.Match.Find("origin.alias.enumerator"); isAliasEnumerator {
		nextLines, err := buffer.NextLineIterator(uint(queryResult.Match.LineRange().Start + 1))

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

		return origin.NewAliasEnumeratorOrigin(location, queryResult.Match, memberMatches...), nil
	}

	if _, isFieldAnnotation := queryResult.Match.Find("origin.fieldannotation"); isFieldAnnotation {
		return origin.NewFieldAnnotationOrigin(location, queryResult.Match), nil
	}

	if _, isEnumAnnotation := queryResult.Match.Find("origin.enumannotation"); isEnumAnnotation {
		position := queryResult.Match.Range().Start
		dockblock, err := buffer.GetTsCommentBlockAt(uint(position.Line), uint(position.Character))

		if err != nil {
			return nil, err
		}

		enumMemberMatches := []nvim.TsQueryMatch{}

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

		enumMemberMatches = append(enumMemberMatches, *membersMatch)
		enumMemberMatches = append(enumMemberMatches, *memberMatches...)
		return origin.NewEnumAnnotationOrigin(location, queryResult.Match, enumMemberMatches...), nil
	}

	return nil, nil
}

func (c *Crawler) getDefinitionOrigin(location nvim.Location) (origin.Origin, error) {
	buffer, err := c.context.Nvim().OpenBuffer(location.Url)

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	originQueryMap := c.getValueOriginQueryMap()
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

	documentation, err := c.getAnnotations(location)

	if err != nil {
		return nil, err
	}

	if _, isFunction := queryResult.Match.Find("origin.function"); isFunction {
		return origin.NewFunctionOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isTable := queryResult.Match.Find("origin.table"); isTable {
		return origin.NewTableOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isValue := queryResult.Match.Find("origin.value"); isValue {
		return origin.NewValueOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isVariable := queryResult.Match.Find("origin.variable"); isVariable {
		return origin.NewVariableOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isModule := queryResult.Match.Find("origin.module"); isModule {
		return origin.NewModuleOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isVirtual := queryResult.Match.Find("origin.meta"); isVirtual {
		return origin.NewVirtualOrigin(location, queryResult.Match, documentation), nil
	}

	return nil, nil
}

func (c *Crawler) getModuleOriginChain(origin *origin.ModuleOrigin) (origin.OriginChain, error) {
	moduleName := origin.Name()
	moduleLocations, err := c.findModuleDefinitionLocations(moduleName)

	if err != nil {
		return nil, err
	}

	return c.getOriginChain(*moduleLocations)
}

func (c *Crawler) getVariableOriginChain(origin *origin.VariableOrigin) (origin.OriginChain, error) {
	assignedName := origin.Name()
	rightValueLocations, err := c.findDefinitionLocations(assignedName)

	if err != nil {
		return nil, err
	}

	return c.getOriginChain(*rightValueLocations)
}

func (c *Crawler) getAnnotations(location nvim.Location) ([]string, error) {
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
