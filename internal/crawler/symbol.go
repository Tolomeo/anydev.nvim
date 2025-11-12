package crawler

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/crawler/symbol"
)

type origin struct {
	location string
	definition    []string
	documentation []string
}

func (o *origin) SetLocation(url string, line uint, character uint) {
	o.location = fmt.Sprintf(`%s:%v:%v`, url, line, character)
}

type source struct {
	path   string
	origin origin
}

type Symbol any

type Namespace struct {
	source source
	fields []*Symbol
	symbol symbol.NamespaceSymbol
}

func (n *Namespace) Path() string {
	return n.source.path
}

func (n *Namespace) Fields() []*Symbol {
	return n.fields
}

func (n *Namespace) AddField(field *Symbol) {
	n.fields = append(n.fields, field)
}

func NewNamespace(c *crawler, path string) (*Namespace, error) {
	assignmentOrigin, err := c.GetVariableOrigin(path)

	if err != nil {
		return nil, fmt.Errorf("Error crawling %s: %w", path, err)
	}

	namespace := Namespace{
		source: source{
			path:   path,
			origin: assignmentOrigin,
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
	declarationOrigin, err := c.GetFunctionOrigin(path)

	if err != nil {
		return nil, fmt.Errorf("Error crawling function %s: %w", path, err)
	}

	function := Function{
		source: source{
			path:   path,
			origin: declarationOrigin,
		},
	}

	return &function, nil
}

type Variable struct {
	source source
}

func NewVariable(c *crawler, path string) (*Variable, error) {
	variableOrigin, err := c.GetVariableOrigin(path)

	if err != nil {
		return nil, fmt.Errorf("Error crawling %s: %w", path, err)
	}

	variable := Variable{
		source: source{
			path:   path,
			origin: variableOrigin,
		},
	}

	return &variable, nil
}
