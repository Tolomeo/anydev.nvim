package crawl

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
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
				Language: symbol.TableFieldAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.table", symbol.TableFieldAssignmentQuery.Query),
			},
			{
				Language: symbol.TableFieldIndexAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.table", symbol.TableFieldIndexAssignmentQuery.Query),
			},
			{
				Language: symbol.MetaVariableAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.meta", symbol.MetaVariableAssignmentQuery.Query),
			},
			{
				Language: symbol.FunctionFieldDotAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", symbol.FunctionFieldDotAssignmentQuery.Query),
			},
			{
				Language: symbol.FunctionFieldIndexAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", symbol.FunctionFieldIndexAssignmentQuery.Query),
			},
			{
				Language: symbol.ModuleRequireAssignmentQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.module", symbol.ModuleRequireAssignmentQuery.Query),
			},
			// dotindex assignment
			// F.T = X.V
			{
				Language: "lua",
				Query: `
				(assignment_statement
					(variable_list
						name: (_)
					) @assignment.left
					(expression_list
						value: (dot_index_expression) @assignment.right 
					)
				) @origin.variable`,
			},
			// variable assignment
			// F.T = X
			{
				Language: "lua",
				Query: `
				(assignment_statement
					(variable_list
						name: (_)
					) @assignment.left
					(expression_list
						value: (identifier) @assignment.right 
					)
				) @origin.variable`,
			},
		},
		treesitter.VARIABLE_DECLARATION: []treesitter.Query{
			{
				Language: symbol.FunctionVariableDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", symbol.FunctionVariableDeclarationQuery.Query),
			},
			{
				Language: symbol.TableVariableDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.table", symbol.TableVariableDeclarationQuery.Query),
			},
		},
		treesitter.FUNCTION_DECLARATION: []treesitter.Query{
			{
				Language: symbol.FunctionDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", symbol.FunctionDeclarationQuery.Query),
			},
			{
				Language: symbol.FunctionFieldMethodDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", symbol.FunctionFieldMethodDeclarationQuery.Query),
			},
			{
				Language: symbol.FunctionFieldDotDeclarationQuery.Language,
				Query:    fmt.Sprintf("(%s) @origin.function", symbol.FunctionFieldDotDeclarationQuery.Query),
			},
		},
		treesitter.ALIAS_ANNOTATION: []treesitter.Query{
			{
				Language: "luadoc",
				Query: fmt.Sprintf(`(
					(alias_annotation
						"@alias"
						.
						(identifier) @alias.name
						.
						(%s) @alias.type
						.
						(comment)? @alias.documentation
						.
					) @origin.alias
					(#eq? @alias.name "%s")
					(#not-eq? @alias.type "")
				)`, anyTypeAnnotationQuery, c.context.Target().Identifier()),
			},
			// Luadoc matches an empty type node even when the type is not present
			// That means that enum aliases have an empty type node defined
			{
				Language: "luadoc",
				Query: fmt.Sprintf(`(
					(alias_annotation
						"@alias"
						.
						(identifier) @alias.name
						.
						(%s) @alias.emptytype
						.
						(comment)? @alias.documentation
						.
					) @origin.alias.enumerator
					(#eq? @alias.name "%s")
					(#eq? @alias.emptytype "")
				)`, anyTypeAnnotationQuery, c.context.Target().Identifier()),
			},
		},
		treesitter.CLASS_ANNOTATION: []treesitter.Query{
			{
				Language: "luadoc",
				Query: fmt.Sprintf(`(
					(class_annotation
						"@class"
						.
						"(exact)"?
						.
						(identifier) @class.name
						.
						(":" 
							. (%s) @class.parent
							("," (%s) @class.parent)*
						)?
					) @origin.class
					(#eq? @class.name "%s")
				)`, anyTypeAnnotationQuery, anyTypeAnnotationQuery, c.context.Target().Identifier()),
			},
		},
		treesitter.FIELD_ANNOTATION: []treesitter.Query{
			{
				Language: "luadoc",
				Query: fmt.Sprintf(`(
					(field_annotation
						"@field"
						.
						([
							(qualifier "public")
							(qualifier "private") @field.private
							(qualifier "protected") @field.protected
							(qualifier "package") @field.package
						 ])?
						.
						(identifier) @field.name
						.
						"?"? @field.optional
						.
						(%s) @field.type
						.
						(comment)? @field.documentation
						.
					) @origin.field
					(#eq? @field.name "%s")
				)`, anyTypeAnnotationQuery, c.context.Target().Name()),
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

func (c *Crawler) getOrigins(locations []nvim.Location) (*symbol.Origins, error) {
	var origin symbol.Origin
	var err error
	originQueryMap := c.getOriginQueryMap()

	for _, location := range locations {
		origin, err = c.getOrigin(location, originQueryMap)

		if err != nil {
			return nil, err
		}

		if origin != nil {
			break
		}
	}

	if origin == nil {
		return nil, nil
	}

	switch o := origin.(type) {
	case *symbol.ModuleOrigin:
		moduleOrigins, err := c.getModuleOrigins(o)
		if err != nil {
			return nil, err
		}
		return symbol.NewOrigins(origin).Merge(moduleOrigins), nil
	case *symbol.VariableOrigin:
		variableOrigins, err := c.getVariableOrigins(o)
		if err != nil {
			return nil, err
		}
		return symbol.NewOrigins(origin).Merge(variableOrigins), nil
	}

	return symbol.NewOrigins(origin), nil
}

func (c *Crawler) getOrigin(location nvim.Location, originQueryMap nvim.TsNodeQueryMap) (symbol.Origin, error) {
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
		return symbol.NewFunctionOrigin(location, *queryResult, documentation), nil
	}

	if _, isTable := queryResult.Match.Find("origin.table"); isTable {
		return symbol.NewTableOrigin(location, *queryResult, documentation), nil
	}

	if _, isVariable := queryResult.Match.Find("origin.variable"); isVariable {
		return symbol.NewVariableOrigin(location, *queryResult, documentation), nil
	}

	if _, isModule := queryResult.Match.Find("origin.module"); isModule {
		return symbol.NewModuleOrigin(location, *queryResult, documentation), nil
	}

	if _, isClass := queryResult.Match.Find("origin.class"); isClass {
		return symbol.NewClassOrigin(location, *queryResult, documentation), nil
	}

	if _, isAlias := queryResult.Match.Find("origin.alias"); isAlias {
		return symbol.NewAliasOrigin(location, *queryResult, documentation), nil
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

		return symbol.NewAliasEnumeratorOrigin(location, *queryResult, documentation), nil
	}

	if _, isField := queryResult.Match.Find("origin.field"); isField {
		return symbol.NewFieldOrigin(location, *queryResult, documentation), nil
	}

	if _, isMeta := queryResult.Match.Find("origin.meta"); isMeta {
		return symbol.NewMetaOrigin(location, *queryResult, documentation), nil
	}

	return nil, fmt.Errorf("Unknown origin match received: location <%+v>, definition <%+v>, documentation <%+v>", location, queryResult, documentation)
}

func (c *Crawler) getModuleOrigins(origin *symbol.ModuleOrigin) (*symbol.Origins, error) {
	moduleName := origin.GetModuleName()
	moduleLocations, err := c.findModuleDefinitionLocations(moduleName)

	if err != nil {
		return nil, err
	}

	return c.getOrigins(*moduleLocations)
}

func (c *Crawler) getVariableOrigins(origin *symbol.VariableOrigin) (*symbol.Origins, error) {
	assignedName := origin.GetAssignedName()
	rightValueLocations, err := c.findDefinitionLocations(assignedName)

	if err != nil {
		return nil, err
	}

	return c.getOrigins(*rightValueLocations)
}
