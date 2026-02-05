package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

var EnumAnnotationQuery = annotation.AtEnumQuery

var EnumAnnotationMembersQuery = annotation.AtEnumMembersQuery

var EnumAnnotationMemberQuery = annotation.AtEnumMemberQuery

type EnumeratorAnnotationOrigin struct {
	origin
	enumAnnotation        annotation.AtEnum
	enumAnnotationMembers []annotation.AtEnumMember
}

func (eo *EnumeratorAnnotationOrigin) Definition() []string {
	return strings.Split(eo.enumAnnotation.Text(), "\n")
}

func (eo *EnumeratorAnnotationOrigin) Name() string {
	return eo.enumAnnotation.Name()
}

func (eo *EnumeratorAnnotationOrigin) Members() []annotation.AtEnumMember {
	return eo.enumAnnotationMembers
}

func NewEnumAnnotationOrigin(location nvim.Location, enumMatch nvim.TsQueryMatch, memberMatches ...nvim.TsQueryMatch) *EnumeratorAnnotationOrigin {
	enumOrigin := EnumeratorAnnotationOrigin{
		origin: origin{
			location:    location,
			captures:    enumMatch,
			annotations: []string{},
		},
		enumAnnotation:        *annotation.NewAtEnum(enumMatch),
		enumAnnotationMembers: []annotation.AtEnumMember{},
	}

	for _, memberMatches := range memberMatches {
		enumOrigin.enumAnnotationMembers = append(enumOrigin.enumAnnotationMembers, *annotation.NewAtEnumMember(memberMatches))
	}

	return &enumOrigin
}
