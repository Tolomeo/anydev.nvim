package crawl

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
)

var typeAnnotationQueries = map[string]string{
	"builtin_type":         "(builtin_type)",
	"identifier":           "(identifier)",
	"array_type":           "(array_type)",
	"table_type":           "(table_type)",
	"table_literal_type":   "(table_literal_type)",
	"union_type":           "(union_type)",
	"parenthesized_type":   "(parenthesized_type)",
	"tuple_type":           "(tuple_type)",
	"function_type":        "(function_type)",
	"member_type":          "(member_type)",
	"optional_type":        "(optional_type)",
	"literal_type":         "(literal_type)",
	"numeric_literal_type": "(numeric_literal_type)",
	"custom_type":          "(custom_type)",
}

var anyTypeAnnotationQuery = fmt.Sprintf(`[%s]`, strings.Join(mapx.Values(typeAnnotationQueries), " "))

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

var enumAliasEnumeratorMemberQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(continuation
		(%s) @alias.type
	)
`, anyTypeAnnotationQuery)}

func (c *Crawler) getOrigins(locations []nvim.Location) (*origin.Origins, error) {
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
		moduleOrigins, err := c.getModuleOrigins(ot)
		if err != nil {
			return nil, err
		}
		return origin.NewOrigins(o).Merge(moduleOrigins), nil
	case *origin.VariableOrigin:
		variableOrigins, err := c.getVariableOrigins(ot)
		if err != nil {
			return nil, err
		}
		return origin.NewOrigins(o).Merge(variableOrigins), nil
	}

	return origin.NewOrigins(o), nil
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
		return origin.NewFunctionOrigin(location, *queryResult, documentation), nil
	}

	if _, isTable := queryResult.Match.Find("origin.table"); isTable {
		return origin.NewTableOrigin(location, *queryResult, documentation), nil
	}

	if _, isVariable := queryResult.Match.Find("origin.variable"); isVariable {
		return origin.NewVariableOrigin(location, *queryResult, documentation), nil
	}

	if _, isModule := queryResult.Match.Find("origin.module"); isModule {
		return origin.NewModuleOrigin(location, *queryResult, documentation), nil
	}

	if _, isClass := queryResult.Match.Find("origin.class"); isClass {
		return origin.NewClassOrigin(location, *queryResult, documentation), nil
	}

	if _, isAlias := queryResult.Match.Find("origin.alias"); isAlias {
		return origin.NewAliasOrigin(location, *queryResult, documentation), nil
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

			match, err := line.TsQueryOne(enumAliasEnumeratorMemberQuery)

			if err != nil {
				return nil, err
			}

			if match == nil {
				break
			}

			queryResult.Match = queryResult.Match.Append(*match...)
		}

		return origin.NewAliasEnumeratorOrigin(location, *queryResult, documentation), nil
	}

	if _, isField := queryResult.Match.Find("origin.field"); isField {
		return origin.NewFieldOrigin(location, *queryResult, documentation), nil
	}

	if _, isMeta := queryResult.Match.Find("origin.meta"); isMeta {
		return origin.NewMetaOrigin(location, *queryResult, documentation), nil
	}

	return nil, fmt.Errorf("Unknown origin match received: location <%+v>, definition <%+v>, documentation <%+v>", location, queryResult, documentation)
}

func (c *Crawler) getModuleOrigins(origin *origin.ModuleOrigin) (*origin.Origins, error) {
	moduleName := origin.GetModuleName()
	moduleLocations, err := c.findModuleDefinitionLocations(moduleName)

	if err != nil {
		return nil, err
	}

	return c.getOrigins(*moduleLocations)
}

func (c *Crawler) getVariableOrigins(origin *origin.VariableOrigin) (*origin.Origins, error) {
	assignedName := origin.GetAssignedName()
	rightValueLocations, err := c.findDefinitionLocations(assignedName)

	if err != nil {
		return nil, err
	}

	return c.getOrigins(*rightValueLocations)
}
