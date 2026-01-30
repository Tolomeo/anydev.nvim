package extract

import (
	"github.com/Tolomeo/anydev.nvim/internal/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/lex"
	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
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

	field := symbol.NewTableField()
	field.Name = name
	field.Metadata = childExtraction.target.Meta()
	field.Documentation = symbol.Documentation{}
	field.Type = symbol.NewUnknown()

	if childExtraction.Target().Type() == nil {
		parent.Fields = append(parent.Fields, *field)
		return nil
	}

	field.Documentation = childExtraction.Target().Documentation()
	field.Type = childExtraction.Target().Type()

	parent.Fields = append(parent.Fields, *field)
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

	childExtraction := t.extractor.newChildExtraction(t.target, name)
	err := t.extractor.extract(childExtraction)

	if err != nil {
		return err
	}

	field := symbol.NewTableField()
	field.Name = name
	field.Metadata = symbol.Metadata{}
	field.Documentation = symbol.Documentation{}
	field.Type = symbol.NewUnknown()

	if childExtraction.Target().Type() == nil {
		parent.Fields = append(parent.Fields, *field)
		return nil
	}

	field.Documentation = childExtraction.Target().Documentation()
	field.Type = childExtraction.target.Type()

	parent.Fields = append(parent.Fields, *field)
	return nil
}

type extractionStep func() (extractionStep, error)

func (t *extraction) getOrigins() (extractionStep, error) {
	t.logger.Info("Crawling")

	origin, err := t.crawler.GetOriginChain()

	if err != nil {
		return nil, err
	}

	if origin == nil {
		t.logger.Warn("Crawling complete, but no origin found")
		return nil, nil
	}

	t.logger.Info("Crawling complete")

	t.Target().SetOriginChain(origin)

	return t.getDocumentation, nil
}

func (t *extraction) getDocumentation() (extractionStep, error) {
	t.logger.Infof("Crawling origin documentation")

	documentation, err := t.crawler.GetDocumentation()

	if err != nil {
		return nil, err
	}

	t.Target().SetDocumentation(documentation)

	return t.getMeta, nil
}

func (t *extraction) getMeta() (extractionStep, error) {
	t.logger.Infof("Crawling origin meta information")

	meta, err := t.lexer.GetMetadata()

	if err != nil {
		return nil, err
	}

	if meta == nil {
		return t.getType, nil
	}

	t.Target().SetMeta(*meta)

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
