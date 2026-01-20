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
	extractor *extractor
	crawler   *crawl.Crawler
	lexer     *lex.Lexer
	logger    *log.Logger
	kind      symbol.TargetKind
	parent    symbol.Target
	name      string
	source    symbol.Source
	symbol    symbol.Symbol
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
	return t.logger
}

func (t *target) Nvim() *nvim.Nvim {
	return t.extractor.Nvim()
}

func (t *target) Extract(kind symbol.TargetKind, name string) error {
	return t.extractor.Extract(kind, name)
}

func (t *target) ExtractChild(name string) (symbol.Symbol, error) {
	return t.extractor.extractChild(t, name)
}

type step func() (step, error)

func (t *target) crawl() (step, error) {
	t.logger.Infof("Crawling '%s'", t.Identifier())

	var err error
	var source symbol.Source

	switch t.Kind() {
	case symbol.TargetKindValue:
		source, err = t.crawler.SourceValue()
	case symbol.TargetKindType:
		source, err = t.crawler.SourceType()
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
	t.logger.Infof("Lexing '%s'", t.Identifier())

	sym, err := t.lexer.Lex()

	if err != nil {
		return nil, err
	}

	t.symbol = sym

	return nil, nil
}

func (e *extractor) newChildTarget(parent symbol.Target, name string) *target {
	t := &target{
		extractor: e,
		parent:    parent,
		kind:      parent.Kind(),
		name:      name,
	}

	t.crawler = crawl.NewCrawler(t)
	t.lexer = lex.NewLexer(t)
	t.logger = log.NewLogger(t.Identifier())

	return t
}

func (e *extractor) newTarget(kind symbol.TargetKind, name string) *target {
	t := &target{
		extractor: e,
		kind:      kind,
		name:      name,
	}

	t.crawler = crawl.NewCrawler(t)
	t.lexer = lex.NewLexer(t)
	t.logger = log.NewLogger(t.Identifier())

	return t
}
