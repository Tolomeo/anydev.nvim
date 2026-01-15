package symbol

func NewTable() *Table {
	return &Table{
		Kind: TableKindTable,
		Fields: []TableField{},
		Indexes: []TableIndex{},
	}
}

func NewTableField() *TableField {
	return &TableField {}
}

func NewTableIndex() *TableIndex {
	return &TableIndex{}
}
