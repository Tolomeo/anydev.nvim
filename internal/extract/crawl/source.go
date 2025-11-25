package crawl

import (
	"fmt"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type origin struct {
	url           string
	line          uint
	character     uint
	definition    []string
	documentation []string
}

type Source interface {
	Path() string
	Origin() origin
}

type NamespaceSource struct {
	path   string
	origin origin
	fields []*Source
}

func (n *NamespaceSource) Path() string {
	return n.path
}

func (n *NamespaceSource) Origin() origin {
	return n.origin
}

func (n *NamespaceSource) Fields() []*Source {
	return n.fields
}

func NewNamespaceSource(c *Crawler, path string) (*NamespaceSource, error) {
	assignmentOrigin, err := c.getOrigin(path, nvim.TS_ASSIGNMENT_STATEMENT)

	if err != nil {
		return nil, fmt.Errorf("Error crawling %s: %w", path, err)
	}

	namespace := NamespaceSource{
		path:   path,
		origin: assignmentOrigin,
	}

	children, err := c.GetFields(path)

	if err != nil {
		return nil, fmt.Errorf("Error crawling namespace %s: %w", path, err)
	}

	for _, field := range children {
		child, err := c.CrawlRuntime(path + "." + field)

		if err != nil {
			return nil, fmt.Errorf("Error crawling %s.%s: %w", path, field, err)
		}

		namespace.fields = append(namespace.fields, &child)
	}

	return &namespace, nil
}

type FunctionSource struct {
	path   string
	origin origin
}

func (f *FunctionSource) Path() string {
	return f.path
}

func (f *FunctionSource) Origin() origin {
	return f.origin
}

func NewFunctionSource(c *Crawler, path string) (*FunctionSource, error) {
	functionOrigin, err := c.getOrigin(path, nvim.TS_FUNCTION_DECLARATION, nvim.TS_ASSIGNMENT_STATEMENT)

	if err != nil {
		return nil, fmt.Errorf("Error crawling function %s: %w", path, err)
	}

	function := FunctionSource{
		path:   path,
		origin: functionOrigin,
	}

	return &function, nil
}

type VariableSource struct {
	path   string
	origin origin
}

func (v *VariableSource) Path() string {
	return v.path
}

func (v *VariableSource) Origin() origin {
	return v.origin
}

func NewVariableSource(c *Crawler, path string) (*VariableSource, error) {
	variableOrigin, err := c.getOrigin(path, nvim.TS_ASSIGNMENT_STATEMENT)

	if err != nil {
		return nil, fmt.Errorf("Error crawling %s: %w", path, err)
	}

	variable := VariableSource{
		path:   path,
		origin: variableOrigin,
	}

	return &variable, nil
}
