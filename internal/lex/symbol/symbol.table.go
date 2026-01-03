package symbol

func NewTable() *Table {
	return &Table{
		Kind: TableKindTable,
	}
}

func NewTableField() *TableField {
	return &TableField {}
}
