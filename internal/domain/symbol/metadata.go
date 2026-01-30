package symbol

import (
	"encoding/json"
	"reflect"
	"strings"
)

type Metadata struct {
	private    bool
	protected  bool
	package_   bool
	deprecated bool
	static     bool
}

func (m Metadata) MarshalJSON() ([]byte, error) {
	val := reflect.ValueOf(m)
	typ := reflect.TypeOf(m)

	metas := make([]string, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		if field.Kind() == reflect.Bool && field.Bool() {
			metas = append(metas, strings.TrimSuffix(strings.ToLower(structField.Name), "_"))
		}
	}

	return json.Marshal(metas)
}

func (s *Metadata) GetDeprecated() bool {
	return s.deprecated
}

func (s *Metadata) GetPackage() bool {
	return s.package_
}

func (s *Metadata) GetPrivate() bool {
	return s.private
}

func (s *Metadata) GetProtected() bool {
	return s.protected
}

func (s *Metadata) SetDeprecated(v bool) {
	s.deprecated = v
}

func (s *Metadata) SetPackage(v bool) {
	s.package_ = v

	if !v {
		return
	}

	s.private = !v
	s.protected = !v
}

func (s *Metadata) SetPrivate(v bool) {
	s.private = v

	if !v {
		return
	}

	s.package_ = !v
	s.protected = !v
}

func (s *Metadata) SetProtected(v bool) {
	s.protected = v

	if !v {
		return
	}

	s.package_ = !v
	s.private = !v
}

func (m *Metadata) GetStatic() bool {
	return m.static
}

func (m *Metadata) SetStatic(v bool) {
	m.static = true
}

func NewMetadata() *Metadata {
	meta := Metadata{}

	return &meta
}
