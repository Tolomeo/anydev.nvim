package symbol

type WithMeta interface {
	GetPrivate() bool
	SetPrivate(bool)
	GetProtected() bool
	SetProtected(bool)
	GetDeprecated() bool
	SetDeprecated(bool)
	GetPackage() bool
	SetPackage(bool)
}

type Meta struct {
	Private    bool `json:"private" yaml:"private"`
	Protected  bool `json:"protected" yaml:"protected"`
	Package    bool `json:"package" yaml:"package"`
	Deprecated bool `json:"deprecated" yaml:"deprecated"`
}

func (s Meta) GetDeprecated() bool {
	return s.Deprecated
}

func (s Meta) GetPackage() bool {
	return s.Package
}

func (s Meta) GetPrivate() bool {
	return s.Private
}

func (s Meta) GetProtected() bool {
	return s.Protected
}

func (s Meta) SetDeprecated(v bool) {
	s.Deprecated = v
}

func (s Meta) SetPackage(v bool) {
	s.Package = v
	s.SetPrivate(!v)
	s.SetProtected(!v)
}

func (s Meta) SetPrivate(v bool) {
	s.Private = v
	s.SetPackage(!v)
	s.SetProtected(!v)
}

func (s Meta) SetProtected(v bool) {
	s.Protected = v
	s.SetPackage(!v)
	s.SetPrivate(!v)
}

var _ WithMeta = (*Meta)(nil)
