package context

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
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
	path   []string
	nvim   *nvim.Nvim
	result *result
	logger *log.Logger
}

func (l *Context) Current() string {
	return strings.Join(l.path, ".")
}

func (l *Context) Nvim() *nvim.Nvim {
	return l.nvim
}

func (l *Context) Logger() *log.Logger {
	return l.logger
}

func (l *Context) Push(path string, procedure func(path string) error) error {
	l.path = append(l.path, path)
	l.Logger().SetKey(l.Current())
	defer func() {
		if len(l.path) > 0 {
			l.path = l.path[:len(l.path)-1]
		}
		if len(l.path) > 0 {
			l.Logger().SetKey(l.Current())
		} else {
			l.Logger().DefaultKey()
		}
	}()

	return procedure(path)
}

func (l *Context) Fork(path string, procedure func(path string) error) error {
	currentPath := l.path
	l.path = []string{path}
	l.Logger().SetKey(l.Current())

	defer func() {
		l.path = currentPath
		l.Logger().SetKey(l.Current())
	}()

	return procedure(path)
}

func (l *Context) Result() *result {
	return l.result
}

func New(logger *log.Logger, client *nvim.Nvim) *Context {
	return &Context{
		nvim:   client,
		logger: logger,
		result: &result{
			Runtime: map[string]symbol.Symbol{},
			Types:   map[string]symbol.Symbol{},
		},
	}
}
