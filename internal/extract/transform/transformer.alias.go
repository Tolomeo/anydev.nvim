package transform

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
)

func (tr *Transformer) getAliasOriginType(aliasOrigin *origin.AliasOrigin) (symbol.Type, error) {
	lexedAlias, err := tr.getType(aliasOrigin.Type())

	if err != nil {
		return nil, err
	}

	return lexedAlias, nil
}

func (tr *Transformer) getAliasEnumeratorType(aliasEnumeratorOrigin *origin.AliasEnumeratorOrigin) (symbol.Type, error) {
	types := []string{}

	for _, aliasEnumMember := range aliasEnumeratorOrigin.Members() {
		types = append(types, aliasEnumMember.Type().Text)
	}

	enumeratorType, err := tr.getType(strings.Join(types, "|"))

	if err != nil {
		return nil, err
	}

	return enumeratorType, nil
}
