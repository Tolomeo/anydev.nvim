package crawler

import (
	"github.com/Tolomeo/anydev.nvim/internal/crawler/symbol"
	// "github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type Symbol any

type source struct {
	// docString  []string
	path       string
	definition []string
	fields     []*Symbol
}

type Namespace struct {
	source source
	symbol symbol.NamespaceSymbol
}

func (n *Namespace) Path() string {
	return n.source.path
}

func (n *Namespace) Fields() []*Symbol {
	return n.source.fields
}

func (n *Namespace) AddField(field *Symbol) {
	n.source.fields = append(n.source.fields, field)
}

/* func FindNamespace(crawler *nvim.Nvim, path string) {
	runtimeType, err := crawler.getRuntimeTypeName(path)

	if err != nil {
		return nil, err
	}

} */

func NewNamespace(path string, definition []string) *Namespace {
	namespace := Namespace{
		source: source{
			path:       path,
			definition: definition,
		},
		symbol: symbol.NamespaceSymbol{},
	}

	return &namespace
}
