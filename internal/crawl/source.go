package crawl

import (
	"regexp"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type origin struct {
	Url           string
	Line          uint
	Character     uint
	Definition    string
	Documentation []string
}

func (o *origin) SetLocation(location nvim.Location) {
	o.Url = location.Url
	o.Line = uint(location.TargetRange.Start.Line)
	o.Character = uint(location.TargetRange.Start.Character)
}

func (o *origin) SetDefinition(definitionLines string) {
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

type Source struct {
	path   string
	origin *origin
	fields []*Source
}

func (s *Source) Path() string {
	return s.path
}

func (s *Source) Origin() *origin {
	return s.origin
}

func (s *Source) Fields() []*Source {
	return s.fields
}
