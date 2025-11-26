package crawl

type origin struct {
	url           string
	line          uint
	character     uint
	definition    []string
	documentation []string
}

type Source interface {
	Path() string
	Origin() *origin
}

type TableSource struct {
	path   string
	origin *origin
	fields []*Source
}

func (n *TableSource) Path() string {
	return n.path
}

func (n *TableSource) Origin() *origin {
	return n.origin
}

func (n *TableSource) Fields() []*Source {
	return n.fields
}

type FunctionSource struct {
	path   string
	origin *origin
}

func (f *FunctionSource) Path() string {
	return f.path
}

func (f *FunctionSource) Origin() *origin {
	return f.origin
}

type VariableSource struct {
	path   string
	origin *origin
}

func (v *VariableSource) Path() string {
	return v.path
}

func (v *VariableSource) Origin() *origin {
	return v.origin
}
