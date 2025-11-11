package crawler

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/crawler/symbol"
	// "github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type Symbol any

type source struct {
	// docString  []string
	path       string
	definition definition
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

func NewNamespace(c *crawler, path string) (*Namespace, error) {
	definition, err := c.GetDefinition(path)

	if err != nil {
		return nil, fmt.Errorf("Error crawling %s: %w", path, err)
	}

	namespace := Namespace{
		source: source{
			path:       path,
			definition: definition,
		},
		symbol: symbol.NamespaceSymbol{},
	}

	children, err := c.GetChildren(path)

	if err != nil {
		return nil, fmt.Errorf("Error crawling %s: %w", path, err)
	}

	for _, field := range children {
		child, err := c.Crawl(path + "." + field)

		if err != nil {
			return nil, fmt.Errorf("Error crawling %s.%s: %w", path, field, err)
		}

		namespace.AddField(&child)
	}

	return &namespace, nil
}
