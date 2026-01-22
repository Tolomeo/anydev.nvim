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
	origin    symbol.Origins
	meta      symbol.Meta
	type_     symbol.Type
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

func (t *target) Origin() symbol.Origins {
	return t.origin
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

func (t *target) ExtractChild(parent *symbol.Table, name string) error {
	t.logger.Infof("Beginning the extraction of '%s' %s child target", name, t.Kind())

	childTarget := t.extractor.newChildTarget(t, name)
	err := t.extractor.extract(childTarget)

	if err != nil {
		return err
	}

	if childTarget.type_ == nil {
		parent.Fields = append(parent.Fields, symbol.NewSymbol(name, symbol.Meta{}, symbol.Documentation{}, symbol.NewUnknown()))
		return nil
	}

	parent.Fields = append(parent.Fields, symbol.NewSymbol(name, symbol.Meta{}, childTarget.Origin().Documentation(), childTarget.type_))
	return nil
}

type step func() (step, error)

func (t *target) getOrigins() (step, error) {
	t.logger.Info("Crawling")

	origin, err := t.crawler.GetOrigins()

	if err != nil {
		return nil, err
	}

	if origin == nil {
		t.logger.Warn("Crawling complete, but no origin found")
		return nil, nil
	}

	t.logger.Info("Crawling complete")

	t.origin = origin

	return t.getMeta, nil
}

func (t *target) getMeta() (step, error) {
	t.logger.Infof("Crawling origin meta information")

	meta, err := t.crawler.GetMeta()

	if err != nil {
		return nil, err
	}

	t.meta = meta

	return t.getType, nil
}

func (t *target) getType() (step, error) {
	t.logger.Infof("Lexing '%s'", t.Identifier())

	sym, err := t.lexer.Lex()

	if err != nil {
		return nil, err
	}

	t.type_ = sym

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
