package origin

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

const (
	VariableDeclaration string = "variable_declaration"
	AssignmentStatement string = "assignment_statement"
	FunctionDeclaration string = "function_declaration"
	Field               string = "field"
	Comment             string = "comment"
	Documentation       string = "documentation"
	ClassAnnotation     string = "class_annotation"
	FieldAnnotation     string = "field_annotation"
	AliasAnnotation     string = "alias_annotation"
	EnumAnnotation      string = "enum_annotation"
)

type Origin interface {
	Url() string
	Line() uint
	Character() uint
	Location() string
	Definition() []string
	Annotations() []string
}

type origin struct {
	location    nvim.Location
	captures    nvim.TsQueryMatch
	annotations []string
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

func (l *origin) Location() string {
	return fmt.Sprintf("%s:%d:%d", l.Url(), l.Line(), l.Character())
}

func (l *origin) Annotations() []string {
	return l.annotations
}
