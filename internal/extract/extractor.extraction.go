package extract

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/annotation"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/extract/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/extract/transform"
	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type extractionContext struct {
	target *extraction
}

func (c *extractionContext) Target() *target.Target {
	return c.target.target
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

func (c *extractionContext) ExtractChild(parent *symbol.Table, name string) (*symbol.TableField, error) {
	c.target.logger.Infof("Beginning the extraction of '%s' %s child target", name, c.target.target.Kind())

	childExtraction := c.target.extractor.newChildExtraction(c.target.target, name)
	err := c.target.extractor.extract(childExtraction)

	if err != nil {
		return nil, err
	}

	field := symbol.NewTableField()
	field.Name = name
	field.Metadata = childExtraction.target.Meta()
	field.Documentation = symbol.Documentation{}
	field.Type = symbol.NewUnknown()

	if childExtraction.target.Type() == nil {
		return field, nil
	}

	field.Documentation = childExtraction.target.Documentation()
	field.Type = childExtraction.target.Type()

	return field, nil
}

func (c *extractionContext) AddChild(parent *symbol.Table, name string, metadata symbol.Metadata, documentation []string, typ annotation.Type) {
	field := symbol.NewTableField()
	field.Name = name
	field.Metadata = metadata
	field.Documentation = documentation
	field.Type = typ
	parent.Fields = append(parent.Fields, *field)
}

type extraction struct {
	extractor   *extractor
	crawler     *crawl.Crawler
	transformer *transform.Transformer
	logger      *log.Logger
	parent      *extraction
	target      *target.Target
}

/* func (t *extraction) Target() *target.Target {
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
} */

type extractionStep func() (extractionStep, error)

func (t *extraction) getOriginChain() (extractionStep, error) {
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

	t.target.SetOriginChain(origin)

	return t.getDocumentation, nil
}

func (t *extraction) getDocumentation() (extractionStep, error) {
	t.logger.Infof("Getting origins chain")

	documentation, err := t.crawler.GetDocumentation()

	if err != nil {
		return nil, err
	}

	t.target.SetDocumentation(documentation)

	return t.getMetadata, nil
}

func (t *extraction) getMetadata() (extractionStep, error) {
	t.logger.Infof("Getting metadata")

	meta, err := t.transformer.GetMetadata()

	if err != nil {
		return nil, err
	}

	if meta == nil {
		return t.getType, nil
	}

	t.target.SetMeta(*meta)

	return t.getType, nil
}

func (t *extraction) getType() (extractionStep, error) {
	t.logger.Info("Getting symbol type")

	typ, err := t.transformer.GetType()

	if err != nil {
		return nil, err
	}

	t.target.SetType(typ)

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
	targetExtraction.transformer = transform.NewTransformer(&targetExtractionContext)
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
	targetExtraction.transformer = transform.NewTransformer(&targetExtractionContext)
	targetExtraction.logger = log.NewLogger(target.Identifier())

	return targetExtraction
}
