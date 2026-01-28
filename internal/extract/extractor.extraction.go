package extract

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/lex"
	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type extractionContext struct {
	target *extraction
}

func (c *extractionContext) Target() *target.Target {
	return c.target.Target()
}

func (c *extractionContext) Logger() *log.Logger {
	return c.target.logger
}

func (c *extractionContext) Nvim() *nvim.Nvim {
	return c.target.extractor.Nvim()
}

func (c *extractionContext) Extract(kind target.TargetKind, name string) error {
	return c.target.extractor.Extract(kind, name)
}

func (c *extractionContext) ExtractChild(parent *symbol.Table, name string) error {
	c.target.logger.Infof("Beginning the extraction of '%s' %s child target", name, c.target.Target().Kind())

	childExtraction := c.target.extractor.newChildExtraction(c.target.Target(), name)
	err := c.target.extractor.extract(childExtraction)

	if err != nil {
		return err
	}

	if childExtraction.Target().Type() == nil {
		parent.Fields = append(parent.Fields, symbol.NewSymbol(name, childExtraction.target.Meta(), symbol.Documentation{}, symbol.NewUnknown()))
		return nil
	}

	parent.Fields = append(parent.Fields, symbol.NewSymbol(name, symbol.Meta{}, childExtraction.Target().Origin().Documentation(), childExtraction.Target().Type()))
	return nil
}

type extraction struct {
	extractor *extractor
	crawler   *crawl.Crawler
	lexer     *lex.Lexer
	logger    *log.Logger
	parent    *extraction
	target    *target.Target
}

func (t *extraction) Target() *target.Target {
	return t.target
}

func (t *extraction) Logger() *log.Logger {
	return t.logger
}

func (t *extraction) Nvim() *nvim.Nvim {
	return t.extractor.Nvim()
}

func (t *extraction) Extract(kind target.TargetKind, name string) error {
	return t.extractor.Extract(kind, name)
}

func (t *extraction) ExtractChild(parent *symbol.Table, name string) error {
	t.logger.Infof("Beginning the extraction of '%s' %s child target", name, t.target.Kind())

	childTarget := t.extractor.newChildExtraction(t.target, name)
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

type extractionStep func() (extractionStep, error)

func (t *extraction) getOrigins() (extractionStep, error) {
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

func (t *extraction) getMeta() (extractionStep, error) {
	t.logger.Infof("Crawling origin meta information")

	meta, err := t.crawler.GetMeta()

	if err != nil {
		return nil, err
	}

	t.Target().SetMeta(meta)

	return t.getType, nil
}

func (t *extraction) getType() (extractionStep, error) {
	t.logger.Infof("Lexing '%s'", t.Target().Identifier())

	typ, err := t.lexer.Lex()

	if err != nil {
		return nil, err
	}

	t.Target().SetType(typ)

	return nil, nil
}

func (e *extractor) newChildExtraction(parent *target.Target, name string) *extraction {
	target := parent.NewChild(name)
	targetExtraction := &extraction{
		extractor: e,
		target:    target,
	}
	targetExtractionContext := extractionContext{
		target: targetExtraction,
	}

	targetExtraction.crawler = crawl.NewCrawler(&targetExtractionContext)
	targetExtraction.lexer = lex.NewLexer(&targetExtractionContext)
	targetExtraction.logger = log.NewLogger(target.Identifier())

	return targetExtraction
}

func (e *extractor) newExtraction(kind target.TargetKind, name string) *extraction {
	target := target.NewTarget(kind, name)
	targetExtraction := &extraction{
		extractor: e,
		target:    target,
	}
	targetExtractionContext := extractionContext{
		target: targetExtraction,
	}

	targetExtraction.crawler = crawl.NewCrawler(&targetExtractionContext)
	targetExtraction.lexer = lex.NewLexer(&targetExtractionContext)
	targetExtraction.logger = log.NewLogger(target.Identifier())

	return targetExtraction
}
