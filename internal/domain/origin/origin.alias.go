package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

var AliasAnnotationQuery = annotation.AtAliasQuery

type AliasOrigin struct {
	origin
	aliasAnnotation annotation.AtAlias
}

func (ao *AliasOrigin) Definition() []string {
	return strings.Split(ao.aliasAnnotation.Match(), "\n")
}

func (ao *AliasOrigin) Name() string {
	return ao.aliasAnnotation.Name()
}

func (ao *AliasOrigin) Type() string {
	return ao.aliasAnnotation.Type()
}

func NewAliasOrigin(location nvim.Location, captures nvim.TsQueryMatch) *AliasOrigin {
	aliasOrigin := AliasOrigin{
		origin: origin{
			location:    location,
			captures:    captures,
			annotations: []string{},
		},
		aliasAnnotation: *annotation.NewAtAlias(captures),
	}

	return &aliasOrigin
}

var AliasEnumeratorAnnotationQuery = annotation.AtAliasEnumeratorQuery

var AliasEnumeratorMemberAnnotationQuery = annotation.AtAliasEnumeratorMemberQuery

type AliasEnumeratorOrigin struct {
	origin
	enumeratorAnnotation        annotation.AtAliasEnumerator
	enumeratorMemberAnnotations []annotation.AtAliasEnumeratorMember
}

func (aeo *AliasEnumeratorOrigin) Definition() []string {
	return strings.Split(aeo.enumeratorAnnotation.Match(), "\n")
}

func (aeo *AliasEnumeratorOrigin) Name() string {
	return aeo.enumeratorAnnotation.Name()
}

func (aeo *AliasEnumeratorOrigin) Members() []annotation.AtAliasEnumeratorMember {
	return aeo.enumeratorMemberAnnotations
}

func NewAliasEnumeratorOrigin(location nvim.Location, aliasCaptures nvim.TsQueryMatch, membersCaptures ...nvim.TsQueryMatch) *AliasEnumeratorOrigin {
	aliasEnumeratorOrigin := AliasEnumeratorOrigin{
		origin: origin{
			location:    location,
			captures:    aliasCaptures,
			annotations: []string{},
		},
		enumeratorAnnotation:        *annotation.NewAtAliasEnumerator(aliasCaptures),
		enumeratorMemberAnnotations: []annotation.AtAliasEnumeratorMember{},
	}

	for _, memberCaptures := range membersCaptures {
		aliasEnumeratorOrigin.enumeratorMemberAnnotations =
			append(aliasEnumeratorOrigin.enumeratorMemberAnnotations, *annotation.NewAtAliasEnumeratorMember(memberCaptures))
	}

	return &aliasEnumeratorOrigin
}
