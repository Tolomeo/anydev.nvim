package symbol

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
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

type Origin interface {
	Url() string
	Line() uint
	Character() uint
	Definition() []string
	Documentation() []string
	Captures() nvim.TsQueryMatch
}

type origin struct {
	location   nvim.Location
	definition nvim.TsNodeQueryMatch
	docBlock   []string
}

func (l *origin) Url() string {
	return l.location.Url
}

func (l *origin) Line() uint {
	return l.location.StartLine()
}

func (l *origin) Character() uint {
	return l.location.StartCharacter()
}

func (l *origin) Type() string {
	return l.definition.Node.Type
}

func (l *origin) Definition() []string {
	return strings.Split(l.definition.Node.Text, "\n")
}

func (l *origin) Documentation() []string {
	return l.docBlock
}

func (l *origin) Captures() nvim.TsQueryMatch {
	return l.definition.Match
}

/*
function fn() end
function fn(arg1) end
function fn(arg1, arg2) end
function fn(arg1, arg2, ...) end
function fn(...) end
*/
var FunctionDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(function_declaration
		name: (identifier) @name
		parameters: (parameters
			(identifier)? @arg
			("," (identifier) @arg)*
			("," (vararg_expression) @vararg)?
			(vararg_expression)? @vararg
		)
	) @function`,
}

/*
local T = function() end
local M = function(arg) end
local D = function(arg, ...) end
local E = function(...) end
*/
var FunctionVariableDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(variable_declaration
		(assignment_statement
			(variable_list
				name: (identifier) @name
			)
			(expression_list
				value: (function_definition
					parameters: (parameters
						(identifier)? @arg
						("," (identifier) @arg)*
						("," (vararg_expression) @vararg)?
						(vararg_expression)? @vararg
					)
				)
			)
		)
	) @function`,
}

/*
api['fn'] = function() end
api['fn'] = function(name) end
api['fn'] = function(name, value) end
api['fn'] = function(name, value, ...) end
api['fn'] = function(...) end
*/
var FunctionFieldIndexAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (bracket_index_expression
				table: (_)
				field: (string
					content: (string_content) @name
				)
			) @access.class
		)
		(expression_list
			value: (function_definition
				parameters: (parameters
					(identifier)? @arg
					("," (identifier) @arg)*
					("," (vararg_expression) @vararg)?
					(vararg_expression)? @vararg
				)
			)
		)
	) @function`,
}

/*
api.fn = function() end
api.fn = function(name) end
api.fn = function(name, value) end
api.fn = function(name, value, ...) end
api.fn = function(...) end
*/
var FunctionFieldDotAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (dot_index_expression
				field: (identifier) @name
			) @access.class
		)
		(expression_list
			value: (function_definition
				parameters: (parameters
					(identifier)? @arg
					("," (identifier) @arg)*
					("," (vararg_expression) @vararg)?
					(vararg_expression)? @vararg
				)
			)
		)
	) @function`,
}

/*
function api.fn() end
function api.fn(name) end
function api.fn(name, value) end
function api.fn(name, value, ...) end
function api.fn(...) end
*/
var FunctionFieldDotDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(function_declaration
		name: (dot_index_expression
			field: (identifier) @name
		) @access.class
		parameters: (parameters
			(identifier)? @arg
			("," (identifier) @arg)*
			("," (vararg_expression) @vararg)?
			(vararg_expression)? @vararg
		)
	) @function`,
}

/*
function api:fn() end
function api:fn(name) end
function api:fn(name, value) end
function api:fn(name, value, ...) end
function api:fn(...) end
*/
var FunctionFieldMethodDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(function_declaration
		name: (method_index_expression
			method: (identifier) @name
		) @access.instance
		parameters: (parameters
			(identifier)? @arg
			("," (identifier) @arg)*
			("," (vararg_expression) @vararg)?
			(vararg_expression)? @vararg
		)
	) @function`,
}

type FunctionOrigin struct {
	origin
}

func (fo *FunctionOrigin) Definition() []string {
	root, _ := fo.definition.Match.Find("function")
	return strings.Split(root.Node.Text, "\n")
}

func NewFunctionOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *FunctionOrigin {
	return &FunctionOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

// local T = {}
var TableVariableDeclarationQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(variable_declaration
		(assignment_statement
			(variable_list
				name: (identifier)
			) @table.name
			(expression_list
				value: (table_constructor)
			) @table.value
		) 
	) @table`,
}

// F.T = {}
var TableFieldAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (dot_index_expression
				table: (_)
				field: (identifier) @table.name
			)
		)
		(expression_list
			value: (table_constructor) @table.value
		)
	) @table`,
}

// F['T'] = {}
var TableFieldIndexAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (bracket_index_expression
				table: (_)
				field: (string
					content: (string_content) @table.name
				)
			)
		)
		(expression_list
			value: (table_constructor) @table.value
		)
	) @table`,
}

type TableOrigin struct {
	origin
}

func (to *TableOrigin) Definition() []string {
	root, _ := to.definition.Match.Find("table")
	return strings.Split(root.Node.Text, "\n")
}

func NewTableOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *TableOrigin {
	return &TableOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

// F.T = X
var VariableAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: (identifier) @assignment.right 
		)
	) @variable`,
}

// F.T = X.V
var VariableDotFieldAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: (dot_index_expression) @assignment.right 
		)
	) @variable`,
}

type VariableOrigin struct {
	origin
}

func (vo *VariableOrigin) Definition() []string {
	root, _ := vo.definition.Match.Find("variable")
	return strings.Split(root.Node.Text, "\n")
}

func (vo *VariableOrigin) GetAssignedName() string {
	assignmentRightCapture, _ := slicesx.FindFunc(vo.definition.Match, func(capture treesitter.Capture) bool {
		return capture.Id == "assignment.right"
	})

	return assignmentRightCapture.Node.Text
}

func NewVariableOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *VariableOrigin {
	return &VariableOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

// F = require("T")
var ModuleRequireAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list)
		(expression_list
			value: (function_call
				name: (identifier) @require.call
				arguments: (arguments
					(string
						content: (string_content) @require.module
					)
				)
			)
		) @require
		(#eq? @require.call "require")
	) @module`,
}

type ModuleOrigin struct {
	origin
}

func (ao *ModuleOrigin) GetModuleName() string {
	moduleNameCapture, _ := slicesx.FindFunc(ao.definition.Match, func(capture treesitter.Capture) bool {
		return capture.Id == "require.module"
	})

	return moduleNameCapture.Node.Text
}

func (ao *ModuleOrigin) Definition() []string {
	root, _ := ao.definition.Match.Find("module")
	return strings.Split(root.Node.Text, "\n")
}

func NewModuleOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *ModuleOrigin {
	return &ModuleOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type ClassOrigin struct {
	origin
}

func NewClassOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *ClassOrigin {
	return &ClassOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

var AliasAnnotationQuery = treesitter.Query{
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
		)
		(#not-eq? @alias.type "")
	) @alias`, anyTypeAnnotationQuery),
}

type AliasOrigin struct {
	origin
}

func (ao *AliasOrigin) Definition() []string {
	root, _ := ao.definition.Match.Find("alias")
	return strings.Split(root.Node.Text, "\n")
}

func NewAliasOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *AliasOrigin {
	return &AliasOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

// Luadoc matches an empty type node even when the type is not present
// That means that enum aliases have an empty type node defined
var AliasEnumeratorAnnotationQuery = treesitter.Query{
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
		)
		(#eq? @alias.emptytype "")
	) @alias.enumerator`, anyTypeAnnotationQuery),
}

type AliasEnumeratorOrigin struct {
	origin
}

func (aeo *AliasEnumeratorOrigin) Definition() []string {
	root, _ := aeo.definition.Match.Find("alias")
	return strings.Split(root.Node.Text, "\n")
}

func NewAliasEnumeratorOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *AliasEnumeratorOrigin {
	return &AliasEnumeratorOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type FieldOrigin struct {
	origin
}

func (fo *FieldOrigin) GetName() string {
	fieldName, _ := fo.definition.Match.Find("field.name")
	return fieldName.Node.Text
}

func (fo *FieldOrigin) GetType() string {
	fieldType, _ := fo.definition.Match.Find("field.type")
	return fieldType.Node.Text
}

func NewFieldOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *FieldOrigin {
	return &FieldOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

// local F = ...
var MetaVariableAssignmentQuery = treesitter.Query{
	Language: "lua",
	Query: `
	(assignment_statement
		(variable_list
			name: (_)
		) @assignment.left
		(expression_list
			value: [
				(vararg_expression) @assignment.right
			] 
		)
	) @meta`,
}

type MetaOrigin struct {
	origin
}

func (mo *MetaOrigin) Definition() []string {
	root, _ := mo.definition.Match.Find("meta")
	return strings.Split(root.Node.Text, "\n")
}

func NewMetaOrigin(location nvim.Location, definition nvim.TsNodeQueryMatch, documentation []string) *MetaOrigin {
	return &MetaOrigin{
		origin: origin{
			location:   location,
			definition: definition,
			docBlock:   documentation,
		},
	}
}

type Origins []Origin

func (o *Origins) Last() Origin {
	return (*o)[len(*o)-1]
}

func (o *Origins) First() Origin {
	return (*o)[0]
}

func (o *Origins) Merge(o2 *Origins) *Origins {
	*o = append(*o, *o2...)
	return o
}

func NewOrigins(origins ...Origin) *Origins {
	t := Origins{}

	for _, l := range origins {
		t = append(t, l)
	}

	return &t
}
