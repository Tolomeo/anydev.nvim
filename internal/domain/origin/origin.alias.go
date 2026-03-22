package origin

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type AliasOrigin struct {
	origin
	node annotation.AtAlias
}

func (ao *AliasOrigin) Definition() []string {
	return strings.Split(ao.node.Root().Text, "\n")
}

func (ao *AliasOrigin) Name() string {
	return ao.node.Name().Text
}

func (ao *AliasOrigin) Type() string {
	return ao.node.Type().Text
}

func NewAliasOrigin(location nvim.Location, node annotation.AtAlias) *AliasOrigin {
	aliasOrigin := AliasOrigin{
		origin: origin{
			location: location,
		},
		node: node,
	}

	return &aliasOrigin
}

type AliasEnumeratorOrigin struct {
	origin
	node        annotation.AtAliasEnumerator
	memberNodes []annotation.AtAliasEnumeratorMember
}

func (aeo *AliasEnumeratorOrigin) Definition() []string {
	return strings.Split(aeo.node.Root().Text, "\n")
}

func (aeo *AliasEnumeratorOrigin) Name() string {
	return aeo.node.Name().Text
}

func (aeo *AliasEnumeratorOrigin) Members() []annotation.AtAliasEnumeratorMember {
	return aeo.memberNodes
}

func NewAliasEnumeratorOrigin(location nvim.Location, node annotation.AtAliasEnumerator, memberNodes ...annotation.AtAliasEnumeratorMember) *AliasEnumeratorOrigin {
	aliasEnumeratorOrigin := AliasEnumeratorOrigin{
		origin: origin{
			location: location,
		},
		node:        node,
		memberNodes: memberNodes,
	}

	return &aliasEnumeratorOrigin
}
