package symbol

func NewTable() *Table {
	return &Table{
		Kind: TableKindTable,
		Fields: []TableField{},
	}
}

func NewTableField() *TableField {
	return &TableField {}
}
