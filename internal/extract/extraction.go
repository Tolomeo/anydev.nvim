package extract

import (
	"strconv"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/extract/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/extract/transform"
	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type extractionContext struct {
	extraction *extraction
}

func (c *extractionContext) Target() *target.Target {
	return c.extraction.target()
}

func (c *extractionContext) TargetDefinitionOverride() *[]string {
	id := c.extraction.target().Identifier()
	originOverrides := c.extraction.extractor.options.Override.Definition

	if originOverrides == nil {
		return nil
	}

	targetOriginOverride, hasOriginOverride := originOverrides[id]

	if !hasOriginOverride {
		return nil
	}

	return &targetOriginOverride
}

func (c *extractionContext) Logger() *log.Logger {
	return c.extraction.logger
}

func (c *extractionContext) Nvim() *nvim.Nvim {
	return c.extraction.extractor.Nvim()
}

func (c *extractionContext) Extract(kind target.TargetKind, name string) error {
	return c.extraction.extractor.Extract(kind, name)
}

func (c *extractionContext) ExtractChild(parent *symbol.Table, name string) error {
	c.extraction.logger.Infof("Beginning the extraction of '%s' %s child target", name, c.extraction.target().Kind())
	c.extraction.extractor.record(name)

	childExtraction := c.extraction.extractor.newChildExtraction(c.extraction.target(), name)
	err := c.extraction.extractor.extract(childExtraction)

	if err != nil {
		return err
	}

	metadata := childExtraction.target().Meta()
	documentation := childExtraction.target().Documentation()
	var type_ symbol.Type
	if childExtraction.target().Type() != nil {
		type_ = childExtraction.target().Type()
	} else {
		type_ = symbol.NewUnknown()
	}

	if _, err := strconv.Atoi(name); err == nil {
		key, err := c.extraction.transformer.GetType(name)
		if err != nil {
			return err
		}
		indexedField := symbol.NewTableIndex()
		indexedField.Key = key
		indexedField.Value = type_
		parent.Indexes = append(parent.Indexes, *indexedField)
		return nil
	}

	switch childOriginType := childExtraction.target().Origin().(type) {
	case *origin.FieldAnnotationOrigin:
		if childOriginType.Index() != nil {
			key, err := c.extraction.transformer.GetType(*childOriginType.Index())
			if err != nil {
				return err
			}
			indexedField := symbol.NewTableIndex()
			indexedField.Key = key
			indexedField.Value = childExtraction.target().Type()
			parent.Indexes = append(parent.Indexes, *indexedField)
			return nil
		}
	}

	namedField := symbol.NewTableField()
	namedField.Name = name
	namedField.Metadata = metadata
	namedField.Documentation = documentation
	namedField.Type = type_

	parent.Fields = append(parent.Fields, *namedField)

	return nil
}

func (c *extractionContext) Follow(name string) (symbol.Type, error) {
	c.extraction.logger.Infof("Following <%s>", name)

	followTarget := target.NewTarget(c.extraction.target().Kind(), name)
	followTarget.SetOriginChain(c.extraction.target().OriginChain())
	c.extraction.targets = append(c.extraction.targets, followTarget)

	origin, err := c.extraction.crawler.FollowOriginChain()

	if err != nil {
		return nil, err
	}

	c.extraction.logger.Info("Follow complete")

	c.extraction.targets = c.extraction.targets[:len(c.extraction.targets)-1]
	c.extraction.target().SetOriginChain(origin)

	return c.extraction.transformer.TransformOrigin()
}

type extraction struct {
	extractor   *extractor
	crawler     *crawl.Crawler
	transformer *transform.Transformer
	logger      *log.Logger
	parent      *extraction
	targets     []*target.Target
}

type extractionStep func() (extractionStep, error)

func (t *extraction) target() *target.Target {
	return t.targets[len(t.targets)-1]
}

func (t *extraction) getOriginChain() (extractionStep, error) {
	t.logger.Info("Crawling symbol origin")

	symbolOrigin, err := t.crawler.GetOriginChain()

	if err != nil {
		return nil, err
	}

	t.logger.Infof("Crawling symbol origin yielded <%T>", symbolOrigin.Last())

	t.target().SetOriginChain(symbolOrigin)

	return t.getDocumentation, nil
}

func (t *extraction) getDocumentation() (extractionStep, error) {
	t.logger.Infof("Crawling documentation")

	documentation, err := t.crawler.GetDocumentation()

	if err != nil {
		return nil, err
	}

	t.target().SetDocumentation(documentation)

	return t.getMetadata, nil
}

func (t *extraction) getMetadata() (extractionStep, error) {
	t.logger.Infof("Crawling metadata")

	meta, err := t.transformer.GetMetadata()

	if err != nil {
		return nil, err
	}

	if meta == nil {
		return t.getType, nil
	}

	t.target().SetMeta(*meta)

	return t.getType, nil
}

func (t *extraction) getType() (extractionStep, error) {
	t.logger.Info("Transforming symbol type")

	typ, err := t.transformer.TransformOrigin()

	if err != nil {
		return nil, err
	}

	t.target().SetType(typ)

	return nil, nil
}
