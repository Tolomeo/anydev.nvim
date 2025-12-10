package lex

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/lex/lexed"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type result struct {
	Runtime map[string]lexed.Symbol `json:"runtime" yaml:"runtime"`
	Types   map[string]lexed.Symbol `json:"types" yaml:"types"`
}

func newResult() *result {
	return &result{
		Runtime: map[string]lexed.Symbol{},
		Types:   map[string]lexed.Symbol{},
	}
}

type lexingContext struct {
	path   []string
	result *result
	logger *log.Logger
}

func (l *lexingContext) Provide(path string, procedure func(path string) error) error {
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

func (l *lexingContext) Info(message string) {
	l.logger.Log(log.NewLog(log.Info, message))
}

func (l *lexingContext) Warn(message string) {
	l.logger.Log(log.NewLog(log.Warn, message))
}

func (l *lexingContext) Error(message string) {
	l.logger.Log(log.NewLog(log.Error, message))
}

func (l *lexingContext) Logs() *log.Logs {
	return l.logger.Logs()
}

func (l *lexingContext) Result() *result {
	return l.result
}

func newLexingContext() *lexingContext {
	return &lexingContext{
		result: newResult(),
		logger: log.NewLogger("_"),
	}
}
