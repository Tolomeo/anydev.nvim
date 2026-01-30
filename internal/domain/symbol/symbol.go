package symbol

import "github.com/Tolomeo/anydev.nvim/internal/domain/annotation"

type Symbol struct {
	Name          string          `json:"name" yaml:"name" mapstructure:"name"`
	Documentation Documentation   `json:"documentation" yaml:"documentation" mapstructure:"documentation"`
	Metadata      Metadata        `json:"meta" yaml:"meta" mapstructure:"meta"`
	Type          annotation.Type `json:"type" yaml:"type" mapstructure:"type"`
}

func NewSymbol(name string, meta Metadata, documentation Documentation, typ annotation.Type) Symbol {
	return Symbol{
		Name:          name,
		Metadata:      meta,
		Documentation: documentation,
		Type:          typ,
	}
}
