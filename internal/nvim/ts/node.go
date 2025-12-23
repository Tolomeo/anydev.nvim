package ts

const (
	VARIABLE_DECLARATION string = "variable_declaration"
	ASSIGNMENT_STATEMENT string = "assignment_statement"
	FUNCTION_DECLARATION string = "function_declaration"
	COMMENT              string = "comment"
	DOCUMENTATION        string = "documentation"
	CLASS_ANNOTATION     string = "class_annotation"
)

func (n *TsNode) Contains(node TsNode) bool {
	if n.Range.Start.Line == node.Range.Start.Line &&
		n.Range.Start.Character > node.Range.Start.Character {
		return false
	}

	if n.Range.End.Line == node.Range.End.Line &&
		n.Range.End.Character < node.Range.End.Character {
		return false
	}

	return n.Range.Start.Line <= node.Range.Start.Line &&
		n.Range.End.Line >= node.Range.End.Line
}
