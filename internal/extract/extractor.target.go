package extract

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex"
	"github.com/Tolomeo/anydev.nvim/internal/lex/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type target struct {
	e      *extractor
	c      *crawl.Crawler
	l      *lex.Lexer
	kind   symbol.TargetKind
	parent symbol.Target
	name   string
	source symbol.Source
	symbol symbol.Symbol
}

func (t *target) Kind() symbol.TargetKind {
	return t.kind
}

func (t *target) ParentName() string {
	if t.parent == nil {
		return ""
	}

	return t.parent.Name()
}

func (t *target) Name() string {
	return t.name
}

func (t *target) Identifier() string {
	if t.parent == nil {
		return t.name
	}

	return fmt.Sprintf("%s.%s", t.parent.Identifier(), t.name)
}

func (t *target) Origin() symbol.Origin {
	return t.source.GetOrigin()
}

func (t *target) Logger() *log.Logger {
	return t.e.logger
}

func (t *target) Nvim() *nvim.Nvim {
	return t.e.Nvim()
}

func (t *target) Extract(kind symbol.TargetKind, name string) error {
	return t.e.Extract(kind, name)
}

func (t *target) ExtractChild(name string) (symbol.Symbol, error) {
	return t.e.extractChild(t, name)
}

type step func() (step, error)

func (t *target) crawl() (step, error) {
	t.Logger().Infof("Crawling '%s'", t.Identifier())

	var err error
	var source symbol.Source

	switch t.Kind() {
	case symbol.TargetKindValue:
		source, err = t.c.SourceValue()
	case symbol.TargetKindType:
		source, err = t.c.SourceType()
	}

	if err != nil {
		return nil, err
	}

	if source == nil {
		return nil, nil
	}

	t.source = source

	return t.lex, nil
}

func (t *target) lex() (step, error) {
	t.Logger().Infof("Lexing '%s'", t.Identifier())

	sym, err := t.l.Lex()

	if err != nil {
		return nil, err
	}

	t.symbol = sym

	return nil, nil
}

func (e *extractor) newChildTarget(parent symbol.Target, name string) *target {
	t := &target{
		e:      e,
		parent: parent,
		kind:   parent.Kind(),
		name:   name,
	}

	t.c = crawl.NewCrawler(t)
	t.l = lex.NewLexer(t)

	return t
}

func (e *extractor) newTarget(kind symbol.TargetKind, name string) *target {
	t := &target{
		e:    e,
		kind: kind,
		name: name,
	}

	t.c = crawl.NewCrawler(t)
	t.l = lex.NewLexer(t)

	return t
}
