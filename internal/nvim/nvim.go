package nvim

type nvim struct {
	path string
}

func (n *nvim) Path() string {
	return n.path
}

func NewNvim(path string) *nvim {
	return &nvim{
		path: path,
	}
}
