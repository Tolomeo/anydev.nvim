package extract

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
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
	c.Logger().Debugf("Extract kind: %s, name: %s", kind, name)
	return c.target.extractor.Extract(kind, name)
}

func (c *extractionContext) ExtractChild(parent *symbol.Table, name string) error {
	c.target.logger.Infof("Beginning the extraction of '%s' %s child target", name, c.target.target.Kind())

	childExtraction := c.target.extractor.newChildExtraction(c.target.target, name)
	err := c.target.extractor.extract(childExtraction)

	if err != nil {
		return err
	}

	namedField := symbol.NewTableField()
	namedField.Name = name
	namedField.Metadata = childExtraction.target.Meta()
	namedField.Documentation = symbol.Documentation{}
	namedField.Type = symbol.NewUnknown()

	if childExtraction.target.Type() == nil {
		parent.Fields = append(parent.Fields, *namedField)
		return nil
	}

	switch childOriginType := childExtraction.target.Origin().(type) {
	case *origin.FieldAnnotationOrigin:
		if childOriginType.Index() != nil {
			key, err := c.target.transformer.GetType(*childOriginType.Index())
			if err != nil {
				return err
			}
			indexedField := symbol.NewTableIndex()
			indexedField.Key = key
			indexedField.Value = childExtraction.target.Type()
			parent.Indexes = append(parent.Indexes, *indexedField)
			return nil
		}
	}

	namedField.Documentation = childExtraction.target.Documentation()
	namedField.Type = childExtraction.target.Type()
	parent.Fields = append(parent.Fields, *namedField)

	return nil
}

func (c *extractionContext) Follow() (symbol.Type, error) {
	c.target.logger.Info("Following")

	origin, err := c.target.crawler.FollowOriginChain()

	if err != nil {
		return nil, err
	}

	c.target.logger.Info("Follow complete")

	c.target.target.SetOriginChain(origin)

	return c.target.transformer.GetOriginType()
}

// TODO: remove
func (c *extractionContext) AddChild(parent *symbol.Table, name string, metadata symbol.Metadata, documentation []string, typ symbol.Type) {
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

	typ, err := t.transformer.GetOriginType()

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
