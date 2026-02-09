package symbol

type Symbol struct {
	Name          string        `json:"name" yaml:"name" mapstructure:"name"`
	Documentation Documentation `json:"documentation" yaml:"documentation" mapstructure:"documentation"`
	Metadata      Metadata      `json:"meta" yaml:"meta" mapstructure:"meta"`
	Type          Type          `json:"type" yaml:"type" mapstructure:"type"`
}

func NewSymbol(name string, meta Metadata, documentation Documentation, typ Type) Symbol {
	return Symbol{
		Name:          name,
		Metadata:      meta,
		Documentation: documentation,
		Type:          typ,
	}
}
