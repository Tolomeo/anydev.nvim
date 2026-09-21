package extract

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/extract/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/extract/export"
	"github.com/Tolomeo/anydev.nvim/internal/extract/transform"
	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type Override struct {
	Definition map[string][]string
}

type Options struct {
	Debug    bool
	LogLevel uint
	OutDir   string
	TmpDir   string
	Override Override
}

// TODO: validate options?
func (o Options) validate() error {
	return nil
}

type extractor struct {
	options     Options
	extractions []*extraction
	manifest    []string
	exporter    *export.Exporter
	nvim        *nvim.Nvim
	logger      *log.Logger
}

func (e *extractor) Nvim() *nvim.Nvim {
	return e.nvim
}

func (e *extractor) extract(item *extraction) error {
	e.logger.Infof("Extracting '%s'", item.target().Identifier())

	e.extractions = append(e.extractions, item)

	var err error
	current := item.getOriginChain

	for {
		current, err = current()

		if err != nil {
			return err
		}

		if current == nil {
			e.extractions = e.extractions[:len(e.extractions)-1]
			return nil
		}
	}
}

func (e *extractor) Extract(kind target.TargetKind, name string) error {
	e.logger.Infof("Beginning the extraction of '%s' %s target", name, kind)
	e.record(name)

	// e.logger.Debugf("%+v", e.result)

	placeholder := symbol.NewSymbol(name, symbol.Metadata{}, symbol.Documentation{}, symbol.NewUnknown())

	switch kind {
	case target.TargetKindValue:
		hasExport := e.exporter.Exported("", name)
		if hasExport {
			e.logger.Infof("Skipping extraction of '%s' %s target: already exported", name, kind)
			return nil
		}
		if err := e.exporter.ExportJson("", name, placeholder); err != nil {
			return err
		}
	default:
		hasExport := e.exporter.Exported(string(kind), name)
		if hasExport {
			e.logger.Infof("Skipping extraction of '%s' %s target: already exported", name, kind)
			return nil
		}
		if err := e.exporter.ExportJson(string(kind), name, placeholder); err != nil {
			return err
		}
	}

	extraction := e.newExtraction(kind, name)
	err := e.extract(extraction)

	if err != nil {
		return err
	}

	metadata := extraction.target().Meta()
	documentation := extraction.target().Documentation()
	type_ := extraction.target().Type()

	e.logger.Infof("The extraction of '%s' %s target yielded \n<%v>", name, kind, type_)

	if extraction.target().Type() == nil {
		return nil
	}

	extractionResult := symbol.NewSymbol(name, metadata, documentation, type_)

	switch kind {
	case target.TargetKindValue:
		if err := e.exporter.ExportJson("", name, extractionResult); err != nil {
			return err
		}
	default:
		if err := e.exporter.ExportJson(string(kind), name, extractionResult); err != nil {
			return err
		}
	}

	return nil
}

func (e *extractor) newExtractionTarget(kind target.TargetKind, name string) *target.Target {
	return target.NewTarget(kind, name)
}

func (e *extractor) newExtraction(kind target.TargetKind, name string) *extraction {
	extractionTarget := e.newExtractionTarget(kind, name)
	targetExtraction := &extraction{
		extractor: e,
		targets:   []*target.Target{extractionTarget},
	}
	targetExtractionContext := extractionContext{
		extraction: targetExtraction,
	}

	targetExtraction.crawler = crawl.NewCrawler(&targetExtractionContext)
	targetExtraction.transformer = transform.NewTransformer(&targetExtractionContext)
	targetExtraction.logger = log.NewLogger(extractionTarget.Identifier(), e.options.LogLevel)

	return targetExtraction
}

func (e *extractor) newChildExtractionTarget(parent *target.Target, name string) *target.Target {
	var childExtractionTarget *target.Target

	switch parent.Kind() {
	case target.TargetKindModule:
		childExtractionTarget = parent.ChildOfKind(target.TargetKindValue, name)
	default:
		childExtractionTarget = parent.Child(name)
	}

	return childExtractionTarget
}

func (e *extractor) newChildExtraction(parent *target.Target, name string) *extraction {
	childExtractionTarget := e.newChildExtractionTarget(parent, name)
	childExtraction := &extraction{
		extractor: e,
		targets:   []*target.Target{childExtractionTarget},
	}
	targetExtractionContext := extractionContext{
		extraction: childExtraction,
	}

	childExtraction.crawler = crawl.NewCrawler(&targetExtractionContext)
	childExtraction.transformer = transform.NewTransformer(&targetExtractionContext)
	childExtraction.logger = log.NewLogger(childExtractionTarget.Identifier(), e.options.LogLevel)

	return childExtraction
}

func (e *extractor) record(name string) {
	e.manifest = append(e.manifest, name)
}

func (e *extractor) initLogger(_ Options) error {
	e.logger = log.NewLogger("", e.options.LogLevel)
	return nil
}

func (e *extractor) initNvim(options Options) error {
	arguments := []string{"--embed", "--headless", "-i", "NONE"}

	if options.Debug {
		verboselevel, verbosefile := 1, path.Join(options.TmpDir, "nvim.verbosefile")
		pipefile := path.Join(options.TmpDir, "nvim.server.pipe")

		arguments = append(arguments, []string{
			fmt.Sprintf("-V%d%s", verboselevel, verbosefile),
			"--listen", pipefile,
		}...)
	}

	client, err := nvim.New(nvim.Options{Arguments: arguments, LogLevel: options.LogLevel})

	if err != nil {
		return fmt.Errorf("Error initialising nvim client: %v", err)
	}

	err = client.Start()

	if err != nil {
		return fmt.Errorf("Error starting nvim client: %w", err)
	}

	e.nvim = client
	return nil
}

func (e *extractor) initExporter(options Options) error {
	e.exporter = export.NewExporter(options.OutDir)
	return nil
}

func (e *extractor) Destroy() error {
	err := e.nvim.Quit()

	if err != nil {
		return err
	}

	e.exporter.ExportTxt("", "manifest", e.manifest)
	e.manifest = []string{}

	return nil
}

func NewExtractor(options Options) (*extractor, error) {
	options.validate()

	newExtractor := &extractor{
		options:     options,
		extractions: []*extraction{},
		manifest:    []string{},
	}

	err := newExtractor.initLogger(options)

	if err != nil {
		return nil, err
	}

	err = newExtractor.initNvim(options)

	if err != nil {
		return nil, err
	}

	err = newExtractor.initExporter(options)

	if err != nil {
		return nil, err
	}

	return newExtractor, nil
}
