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

type withMeta struct {
	Private    bool `json:"private" yaml:"private"`
	Protected  bool `json:"protected" yaml:"protected"`
	Package    bool `json:"package" yaml:"package"`
	Deprecated bool `json:"deprecated" yaml:"deprecated"`
}

func (s withMeta) GetDeprecated() bool {
	return s.Deprecated
}

func (s withMeta) GetPackage() bool {
	return s.Package
}

func (s withMeta) GetPrivate() bool {
	return s.Private
}

func (s withMeta) GetProtected() bool {
	return s.Protected
}

func (s withMeta) SetDeprecated(v bool) {
	s.Deprecated = v
}

func (s withMeta) SetPackage(v bool) {
	s.Package = v
	s.SetPrivate(!v)
	s.SetProtected(!v)
}

func (s withMeta) SetPrivate(v bool) {
	s.Private = v
	s.SetPackage(!v)
	s.SetProtected(!v)
}

func (s withMeta) SetProtected(v bool) {
	s.Protected = v
	s.SetPackage(!v)
	s.SetPrivate(!v)
}

var _ WithMeta = (*withMeta)(nil)
