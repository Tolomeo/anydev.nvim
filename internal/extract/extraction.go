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
	return c.extraction.currentTarget()
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
	c.extraction.logger.Infof("Beginning the extraction of '%s' %s child target", name, c.extraction.currentTarget().Kind())

	childExtraction := c.extraction.extractor.newChildExtraction(c.extraction.currentTarget(), name)
	err := c.extraction.extractor.extract(childExtraction)

	if err != nil {
		return err
	}

	metadata := childExtraction.currentTarget().Meta()
	documentation := childExtraction.currentTarget().Documentation()
	var type_ symbol.Type
	if childExtraction.currentTarget().Type() != nil {
		type_ = childExtraction.currentTarget().Type()
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

	switch childOriginType := childExtraction.currentTarget().Origin().(type) {
	case *origin.FieldAnnotationOrigin:
		if childOriginType.Index() != nil {
			key, err := c.extraction.transformer.GetType(*childOriginType.Index())
			if err != nil {
				return err
			}
			indexedField := symbol.NewTableIndex()
			indexedField.Key = key
			indexedField.Value = childExtraction.currentTarget().Type()
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

	c.extraction.target = append(c.extraction.target, c.extraction.currentTarget().Follow(name))

	origin, err := c.extraction.crawler.FollowOriginChain()

	if err != nil {
		return nil, err
	}

	c.extraction.logger.Info("Follow complete")

	c.extraction.target = c.extraction.target[:len(c.extraction.target)-1]
	c.extraction.currentTarget().SetOriginChain(origin)

	return c.extraction.transformer.GetOriginType()
}

type extraction struct {
	extractor   *extractor
	crawler     *crawl.Crawler
	transformer *transform.Transformer
	logger      *log.Logger
	parent      *extraction
	target      []*target.Target
}

type extractionStep func() (extractionStep, error)

func (t *extraction) currentTarget() *target.Target {
	return t.target[len(t.target)-1]
}

func (t *extraction) getOriginChain() (extractionStep, error) {
	t.logger.Info("Crawling symbol origin")

	symbolOrigin, err := t.crawler.GetOriginChain()

	if err != nil {
		return nil, err
	}

	t.logger.Infof("Crawling symbol origin yielded <%T>", symbolOrigin.Last())

	t.currentTarget().SetOriginChain(symbolOrigin)

	return t.getDocumentation, nil
}

func (t *extraction) getDocumentation() (extractionStep, error) {
	t.logger.Infof("Getting origins chain")

	documentation, err := t.crawler.GetDocumentation()

	if err != nil {
		return nil, err
	}

	t.currentTarget().SetDocumentation(documentation)

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

	t.currentTarget().SetMeta(*meta)

	return t.getType, nil
}

func (t *extraction) getType() (extractionStep, error) {
	t.logger.Info("Getting symbol type")

	typ, err := t.transformer.GetOriginType()

	if err != nil {
		return nil, err
	}

	t.currentTarget().SetType(typ)

	return nil, nil
}
