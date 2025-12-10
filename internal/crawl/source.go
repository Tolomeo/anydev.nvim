package crawl

import (
	"regexp"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type origin struct {
	Url           string
	Line          uint
	Character     uint
	Definition    []string
	Documentation []string
}

func (o *origin) SetLocation(location nvim.Location) {
	o.Url = location.Url
	o.Line = uint(location.TargetRange.Start.Line)
	o.Character = uint(location.TargetRange.Start.Character)
}

func (o *origin) SetDefinition(definitionLines []string) {
	o.Definition = definitionLines
}

func (o *origin) SetDocumentation(sourceDocumentationLines []string) {
	eCommentContent := regexp.MustCompile(`^[ \t]*-{2,3}(.*)$`)

	documentation := []string{}

	for _, sourceLine := range sourceDocumentationLines {
		matches := eCommentContent.FindStringSubmatch(sourceLine)

		if len(matches) < 2 {
			documentation = append(documentation, "")
			continue
		}

		documentation = append(documentation, matches[1])
	}

	o.Documentation = documentation
}

type Source interface {
	Path() string
	Origin() *origin
}

type TableSource struct {
	path   string
	origin *origin
	fields []*Source
}

func (n *TableSource) Path() string {
	return n.path
}

func (n *TableSource) Origin() *origin {
	return n.origin
}

func (n *TableSource) Fields() []*Source {
	return n.fields
}

type FunctionSource struct {
	path   string
	origin *origin
}

func (f *FunctionSource) Path() string {
	return f.path
}

func (f *FunctionSource) Origin() *origin {
	return f.origin
}

type VariableSource struct {
	path   string
	origin *origin
}

func (v *VariableSource) Path() string {
	return v.path
}

func (v *VariableSource) Origin() *origin {
	return v.origin
}
