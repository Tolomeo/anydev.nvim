package symbol

import (
	"encoding/json"
	"reflect"
	"strings"
)

type WithMeta interface {
	GetPrivate() bool
	SetPrivate(bool)
	GetProtected() bool
	SetProtected(bool)
	GetDeprecated() bool
	SetDeprecated(bool)
	GetPackage() bool
	SetPackage(bool)
	GetStatic() bool
	SetStatic(bool)
}

type Meta struct {
	private    bool
	protected  bool
	package_   bool
	deprecated bool
	static     bool
}

func (m Meta) MarshalJSON() ([]byte, error) {
	val := reflect.ValueOf(m)
	typ := reflect.TypeOf(m)

	metas := make([]string, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		if field.Kind() == reflect.Bool && field.Bool() {
			metas = append(metas, strings.ToLower(structField.Name))
		}
	}

	return json.Marshal(metas)
}

func (s *Meta) GetDeprecated() bool {
	return s.deprecated
}

func (s *Meta) GetPackage() bool {
	return s.package_
}

func (s *Meta) GetPrivate() bool {
	return s.private
}

func (s *Meta) GetProtected() bool {
	return s.protected
}

func (s *Meta) SetDeprecated(v bool) {
	s.deprecated = v
}

func (s *Meta) SetPackage(v bool) {
	s.package_ = v

	if !v {
		return
	}

	s.private = !v
	s.protected = !v
}

func (s *Meta) SetPrivate(v bool) {
	s.private = v

	if !v {
		return
	}

	s.package_ = !v
	s.protected = !v
}

func (s *Meta) SetProtected(v bool) {
	s.protected = v

	if !v {
		return
	}

	s.package_ = !v
	s.private = !v
}

func (m *Meta) GetStatic() bool {
	return m.static
}

func (m *Meta) SetStatic(v bool) {
	m.static = true
}

var _ WithMeta = (*Meta)(nil)

func NewMeta() *Meta {
	meta := Meta{}

	return &meta

}
