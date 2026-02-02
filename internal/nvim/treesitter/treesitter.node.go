package treesitter


func (n *TsNode) Contains(node TsNode) bool {
	return n.Range.Contains(node.Range)
}
