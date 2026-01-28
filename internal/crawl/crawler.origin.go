package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

func (c *Crawler) getOriginQueryMap() nvim.TsNodeQueryMap {
	return nvim.TsNodeQueryMap{
		treesitter.ASSIGNMENT_STATEMENT: []treesitter.Query{
			{
				Language: origin.TableFieldAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.table", origin.TableFieldAssignmentQuery.Query),
			},
			{
				Language: origin.TableFieldIndexAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.table", origin.TableFieldIndexAssignmentQuery.Query),
			},
			{
				Language: origin.MetaVariableAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.meta", origin.MetaVariableAssignmentQuery.Query),
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
		treesitter.VARIABLE_DECLARATION: []treesitter.Query{
			{
				Language: origin.FunctionVariableDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", origin.FunctionVariableDeclarationQuery.Query),
			},
			{
				Language: origin.TableVariableDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.table", origin.TableVariableDeclarationQuery.Query),
			},
		},
		treesitter.FUNCTION_DECLARATION: []treesitter.Query{
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
		treesitter.ALIAS_ANNOTATION: []treesitter.Query{
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
		treesitter.CLASS_ANNOTATION: []treesitter.Query{
			{
				Language: origin.ClassAnnotationQuery.Language,
				Query: fmt.Sprintf(`(
					(%s)
					(#eq? @class.name "%s")
				) @origin.class`, origin.ClassAnnotationQuery.Query, c.context.Target().Identifier()),
			},
		},
		treesitter.FIELD_ANNOTATION: []treesitter.Query{
			{
				Language: origin.FieldAnnotationQuery.Language,
				Query: fmt.Sprintf(`(
					(%s)
					(#eq? @field.name "%s")
				) @origin.field`, origin.FieldAnnotationQuery.Query, c.context.Target().Name()),
			},
		},
	}
}

func (c *Crawler) getOriginChain(locations []nvim.Location) (*origin.OriginChain, error) {
	var o origin.Origin
	var err error
	originQueryMap := c.getOriginQueryMap()

	for _, location := range locations {
		o, err = c.getOrigin(location, originQueryMap)

		if err != nil {
			return nil, err
		}

		if o != nil {
			break
		}
	}

	if o == nil {
		return nil, nil
	}

	switch ot := o.(type) {
	case *origin.ModuleOrigin:
		moduleOrigins, err := c.getModuleOriginChain(ot)
		if err != nil {
			return nil, err
		}
		return origin.NewOriginChain(o).Append(moduleOrigins), nil
	case *origin.VariableOrigin:
		variableOrigins, err := c.getVariableOriginChain(ot)
		if err != nil {
			return nil, err
		}
		return origin.NewOriginChain(o).Append(variableOrigins), nil
	}

	return origin.NewOriginChain(o), nil
}

func (c *Crawler) getOrigin(location nvim.Location, originQueryMap nvim.TsNodeQueryMap) (origin.Origin, error) {
	buffer, err := c.context.Nvim().OpenBuffer(location.Url)

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

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

	documentation, err := c.getCommentBlock(queryResult.Node, location)

	if err != nil {
		return nil, err
	}

	if _, isFunction := queryResult.Match.Find("origin.function"); isFunction {
		return origin.NewFunctionOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isTable := queryResult.Match.Find("origin.table"); isTable {
		return origin.NewTableOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isVariable := queryResult.Match.Find("origin.variable"); isVariable {
		return origin.NewVariableOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isModule := queryResult.Match.Find("origin.module"); isModule {
		return origin.NewModuleOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isClass := queryResult.Match.Find("origin.class"); isClass {
		return origin.NewClassOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isAlias := queryResult.Match.Find("origin.alias"); isAlias {
		return origin.NewAliasOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isAliasEnumerator := queryResult.Match.Find("origin.alias.enumerator"); isAliasEnumerator {
		nextLines, err := buffer.NextLineIterator(uint(queryResult.Match.LineRange().Start + 1))

		if err != nil {
			return nil, err
		}

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

			queryResult.Match = queryResult.Match.Append(*match...)
		}

		return origin.NewAliasEnumeratorOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isField := queryResult.Match.Find("origin.field"); isField {
		return origin.NewFieldOrigin(location, queryResult.Match, documentation), nil
	}

	if _, isMeta := queryResult.Match.Find("origin.meta"); isMeta {
		return origin.NewMetaOrigin(location, queryResult.Match, documentation), nil
	}

	return nil, fmt.Errorf("Unknown origin match received: location <%+v>, definition <%+v>, documentation <%+v>", location, queryResult, documentation)
}

func (c *Crawler) getModuleOriginChain(origin *origin.ModuleOrigin) (*origin.OriginChain, error) {
	moduleName := origin.Name()
	moduleLocations, err := c.findModuleDefinitionLocations(moduleName)

	if err != nil {
		return nil, err
	}

	return c.getOriginChain(*moduleLocations)
}

func (c *Crawler) getVariableOriginChain(origin *origin.VariableOrigin) (*origin.OriginChain, error) {
	assignedName := origin.Name()
	rightValueLocations, err := c.findDefinitionLocations(assignedName)

	if err != nil {
		return nil, err
	}

	return c.getOriginChain(*rightValueLocations)
}
