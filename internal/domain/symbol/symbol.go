package symbol

type Symbol struct {
	Meta
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	Documentation Documentation `json:"documentation" yaml:"documentation" mapstructure:"documentation"`
	Type Type `json:"type" yaml:"type" mapstructure:"type"`
}

func NewSymbol(name string, meta Meta, documentation Documentation, type_ Type) Symbol {
	return Symbol{
		Name:          name,
		Meta:          meta,
		Documentation: documentation,
		Type:          type_,
	}
}
