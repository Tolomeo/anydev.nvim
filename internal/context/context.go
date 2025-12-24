package context

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type result struct {
	Runtime map[string]symbol.Symbol `json:"runtime" yaml:"runtime"`
	Types   map[string]symbol.Symbol `json:"types" yaml:"types"`
}

type loggerProvider interface {
	SetKey(key string)
	DefaultKey()
	Info(message string)
	Warn(message string)
	Error(message string)
}

type Context struct {
	path    []string
	Nvim    *nvim.Nvim
	result  *result
	Logger  loggerProvider
}

func (l *Context) Current() string {
	return strings.Join(l.path, ".")
}

func (l *Context) Provide(path string, procedure func(path string) (error)) (error) {
	l.push(path)
	defer l.pop()

	return procedure(path)
}

func (l *Context) push(prefix string) {
	l.path = append(l.path, prefix)
	l.Logger.SetKey(l.Current())
}

func (l *Context) pop() {
	if len(l.path) > 0 {
		l.path = l.path[:len(l.path)-1]
	}
	if len(l.path) > 0 {
		l.Logger.SetKey(l.Current())
	} else {
		l.Logger.DefaultKey()
	}
}

func (l *Context) Result() *result {
	return l.result
}

func New(logger loggerProvider, client *nvim.Nvim) *Context {
	return &Context{
		Nvim:    client,
		Logger:  logger,
		result: &result{
			Runtime: map[string]symbol.Symbol{},
			Types:   map[string]symbol.Symbol{},
		},
	}
}
