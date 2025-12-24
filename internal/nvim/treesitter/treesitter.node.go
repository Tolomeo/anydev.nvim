package treesitter

const (
	VARIABLE_DECLARATION string = "variable_declaration"
	ASSIGNMENT_STATEMENT string = "assignment_statement"
	FUNCTION_DECLARATION string = "function_declaration"
	COMMENT              string = "comment"
	DOCUMENTATION        string = "documentation"
	CLASS_ANNOTATION     string = "class_annotation"
	ALIAS_ANNOTATION     string = "alias_annotation"
)

func (n *TsNode) Contains(node TsNode) bool {
	return n.Range.Contains(node.Range)
}
