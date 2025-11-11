package crawler

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/crawler/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
)

type definition struct {
	definition    []string
	documentation []string
}

type location struct {
	lsp.DefinitionLocation
	Url string
}

type source struct {
	// docString  []string
	path       string
	definition definition
	fields     []*Symbol
}

type Symbol any

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
	definition, err := c.GetAssignmentDescription(path)

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

	children, err := c.GetFields(path)

	if err != nil {
		return nil, fmt.Errorf("Error crawling namespace %s: %w", path, err)
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

type Function struct {
	source source
}

func NewFunction(c *crawler, path string) (*Function, error) {
	definition, err := c.GetDeclarationDefinition(path)

	if err != nil {
		return nil, fmt.Errorf("Error crawling function %s: %w", path, err)
	}

	function := Function{
		source: source{
			path:       path,
			definition: definition,
		},
	}

	return &function, nil
}
