package transform

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
)

func (tr *Transformer) getAliasOriginType(aliasOrigin *origin.AliasOrigin) (annotation.Type, error) {
	lexedAlias, err := tr.getType(aliasOrigin.Type())

	if err != nil {
		return nil, err
	}

	return lexedAlias, nil
}

func (tr *Transformer) getAliasEnumeratorType(aliasEnumeratorOrigin *origin.AliasEnumeratorOrigin) (annotation.Type, error) {
	enumeratorType, err := tr.getType(strings.Join(aliasEnumeratorOrigin.Types(), "|"))

	if err != nil {
		return nil, err
	}

	return enumeratorType, nil
}
