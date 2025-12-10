package lex

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type result struct {
	Runtime map[string]lexed.Symbol `json:"runtime" yaml:"runtime"`
	Types   map[string]lexed.Symbol `json:"types" yaml:"types"`
}

type lexingContext struct {
	path    []string
	nvim    *nvim.Nvim
	crawler *crawl.Crawler
	result  *result
	logger  *log.Logger
}

func (l *lexingContext) provide(path string, procedure func(path string) error) error {
	l.push(path)
	defer l.pop()

	return procedure(path)
}

func (l *lexingContext) push(prefix string) {
	l.path = append(l.path, prefix)
	l.logger.SetKey(strings.Join(l.path, ""))
}

func (l *lexingContext) pop() {
	if len(l.path) > 0 {
		l.path = l.path[:len(l.path)-1]
		l.logger.SetKey(strings.Join(l.path, ""))
	}
}

func (l *lexingContext) Result() *result {
	return l.result
}

func NewLexingContext(logger *log.Logger, client *nvim.Nvim, crawler *crawl.Crawler) *lexingContext {
	return &lexingContext{
		nvim:    client,
		crawler: crawler,
		logger:  logger,
		result: &result{
			Runtime: map[string]lexed.Symbol{},
			Types:   map[string]lexed.Symbol{},
		},
	}
}
