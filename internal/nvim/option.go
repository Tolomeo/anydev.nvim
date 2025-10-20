package nvim

type nvimOptions struct {
	path string
	vimrc string
}

func (o *nvimOptions) Set(opts ...nvimOptionProvider) {
	for _, opt := range opts {
		opt(o)
	}
}

type nvimOptionProvider func(*nvimOptions)

func WithPath(path string) nvimOptionProvider {
	return func(o *nvimOptions) {
		o.path = path
	}
}

func WithVimrc(vimrc string) nvimOptionProvider {
	return func(o *nvimOptions) {
		o.vimrc = vimrc
	}
}

