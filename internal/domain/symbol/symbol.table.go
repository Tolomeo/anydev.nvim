package symbol

const TableKind string = "table"

type Table struct {
	Fields  []TableField `json:"fields" yaml:"fields" mapstructure:"fields"`
	Indexes []TableIndex `json:"indexes" yaml:"indexes" mapstructure:"indexes"`
	Kind    string       `json:"kind" yaml:"kind" mapstructure:"kind"`
	Name    string       `json:"name" yaml:"name" mapstructure:"name"`
}

func (t *Table) Canonical() Type {
	return t
}

func (t *Table) GetKind() string {
	return t.Kind
}

var _ Type = (*Table)(nil)

type TableField struct {
	Symbol
}

type TableIndex struct {
	Key   Type `json:"key" yaml:"key" mapstructure:"key"`
	Value Type `json:"value" yaml:"value" mapstructure:"value"`
}

type TableIndexValue Symbol

func NewTable() *Table {
	return &Table{
		Kind:    TableKind,
		Fields:  []TableField{},
		Indexes: []TableIndex{},
	}
}

func NewTableField() *TableField {
	return &TableField{}
}

func NewTableIndex() *TableIndex {
	return &TableIndex{}
}
