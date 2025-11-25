package crawler

import (
	"fmt"

	// "github.com/Tolomeo/anydev.nvim/internal/crawler/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type origin struct {
	url string
	line uint
	character uint
	definition    []string
	documentation []string
}

type Source interface {
	Path() string
	Origin() origin
}

type Namespace struct {
	path   string
	origin origin
	fields []*Source
}

func (n *Namespace) Path() string {
	return n.path
}

func (n *Namespace) Origin() origin {
	return n.origin
}

func (n *Namespace) Fields() []*Source {
	return n.fields
}

func NewNamespace(c *Crawler, path string) (*Namespace, error) {
	assignmentOrigin, err := c.getOrigin(path, nvim.TS_ASSIGNMENT_STATEMENT)

	if err != nil {
		return nil, fmt.Errorf("Error crawling %s: %w", path, err)
	}

	namespace := Namespace{
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

type Function struct {
	path   string
	origin origin
}

func (f *Function) Path() string {
	return f.path
}

func (f *Function) Origin() origin {
	return f.origin
}

func NewFunction(c *Crawler, path string) (*Function, error) {
	functionOrigin, err := c.getOrigin(path, nvim.TS_FUNCTION_DECLARATION, nvim.TS_ASSIGNMENT_STATEMENT)

	if err != nil {
		return nil, fmt.Errorf("Error crawling function %s: %w", path, err)
	}

	function := Function{
		path:   path,
		origin: functionOrigin,
	}

	return &function, nil
}

type Variable struct {
	path   string
	origin origin
}

func (v *Variable) Path() string {
	return v.path
}

func (v *Variable) Origin() origin {
	return v.origin
}

func NewVariable(c *Crawler, path string) (*Variable, error) {
	variableOrigin, err := c.getOrigin(path, nvim.TS_ASSIGNMENT_STATEMENT)

	if err != nil {
		return nil, fmt.Errorf("Error crawling %s: %w", path, err)
	}

	variable := Variable{
		path:   path,
		origin: variableOrigin,
	}

	return &variable, nil
}
