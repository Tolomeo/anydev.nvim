package symbol

import type_ "github.com/Tolomeo/anydev.nvim/internal/domain/type"

type Symbol struct {
	Meta
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	Documentation Documentation `json:"documentation" yaml:"documentation" mapstructure:"documentation"`
	Type type_.Type `json:"type" yaml:"type" mapstructure:"type"`
}

func NewSymbol(name string, meta Meta, documentation Documentation, typ type_.Type) Symbol {
	return Symbol{
		Name:          name,
		Meta:          meta,
		Documentation: documentation,
		Type:          typ,
	}
}
