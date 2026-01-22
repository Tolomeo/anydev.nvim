package extract

import (
	// "fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex"
	"github.com/Tolomeo/anydev.nvim/internal/lex/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type context struct {
	target *extractorTarget
}

func (c *context) Target() *symbol.Target {
	return c.target.Target()
}

func (c *context) Logger() *log.Logger {
	return c.target.logger
}

func (c *context) Nvim() *nvim.Nvim {
	return c.target.extractor.Nvim()
}

func (c *context) Extract(kind symbol.TargetKind, name string) error {
	return c.target.extractor.Extract(kind, name)
}

func (c *context) ExtractChild(parent *symbol.Table, name string) error {
	c.target.logger.Infof("Beginning the extraction of '%s' %s child target", name, c.target.Target().Kind())

	childTarget := c.target.extractor.newChildTarget(c.target.Target(), name)
	err := c.target.extractor.extract(childTarget)

	if err != nil {
		return err
	}

	if childTarget.Target().Type() == nil {
		parent.Fields = append(parent.Fields, symbol.NewSymbol(name, childTarget.target.Meta(), symbol.Documentation{}, symbol.NewUnknown()))
		return nil
	}

	parent.Fields = append(parent.Fields, symbol.NewSymbol(name, symbol.Meta{}, childTarget.Target().Origin().Documentation(), childTarget.Target().Type()))
	return nil
}

type extractorTarget struct {
	extractor *extractor
	crawler   *crawl.Crawler
	lexer     *lex.Lexer
	logger    *log.Logger
	parent    *extractorTarget
	target    *symbol.Target
}

func (t *extractorTarget) Target() *symbol.Target {
	return t.target
}

func (t *extractorTarget) Logger() *log.Logger {
	return t.logger
}

func (t *extractorTarget) Nvim() *nvim.Nvim {
	return t.extractor.Nvim()
}

func (t *extractorTarget) Extract(kind symbol.TargetKind, name string) error {
	return t.extractor.Extract(kind, name)
}

func (t *extractorTarget) ExtractChild(parent *symbol.Table, name string) error {
	t.logger.Infof("Beginning the extraction of '%s' %s child target", name, t.target.Kind())

	childTarget := t.extractor.newChildTarget(t.target, name)
	err := t.extractor.extract(childTarget)

	if err != nil {
		return err
	}

	if childTarget.Target().Type() == nil {
		parent.Fields = append(parent.Fields, symbol.NewSymbol(name, childTarget.target.Meta(), symbol.Documentation{}, symbol.NewUnknown()))
		return nil
	}

	parent.Fields = append(parent.Fields, symbol.NewSymbol(name, symbol.Meta{}, childTarget.Target().Origin().Documentation(), childTarget.Target().Type()))
	return nil
}

type step func() (step, error)

func (t *extractorTarget) getOrigins() (step, error) {
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

	t.Target().SetOrigins(origin)

	return t.getMeta, nil
}

func (t *extractorTarget) getMeta() (step, error) {
	t.logger.Infof("Crawling origin meta information")

	meta, err := t.crawler.GetMeta()

	if err != nil {
		return nil, err
	}

	t.Target().SetMeta(meta)

	return t.getType, nil
}

func (t *extractorTarget) getType() (step, error) {
	t.logger.Infof("Lexing '%s'", t.Target().Identifier())

	typ, err := t.lexer.Lex()

	if err != nil {
		return nil, err
	}

	t.Target().SetType(typ)

	return nil, nil
}

func (e *extractor) newChildTarget(parent *symbol.Target, name string) *extractorTarget {
	target := parent.NewChild(name)
	targetExtraction := &extractorTarget{
		extractor: e,
		target:    target,
	}
	targetExtractionContext := context{
		target: targetExtraction,
	}

	targetExtraction.crawler = crawl.NewCrawler(&targetExtractionContext)
	targetExtraction.lexer = lex.NewLexer(&targetExtractionContext)
	targetExtraction.logger = log.NewLogger(target.Identifier())

	return targetExtraction
}

func (e *extractor) newTarget(kind symbol.TargetKind, name string) *extractorTarget {
	target := symbol.NewTarget(kind, name)
	targetExtraction := &extractorTarget{
		extractor: e,
		target:    target,
	}
	targetExtractionContext := context{
		target: targetExtraction,
	}

	targetExtraction.crawler = crawl.NewCrawler(&targetExtractionContext)
	targetExtraction.lexer = lex.NewLexer(&targetExtractionContext)
	targetExtraction.logger = log.NewLogger(target.Identifier())

	return targetExtraction
}
