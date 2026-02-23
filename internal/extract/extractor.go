package extract

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/extract/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/extract/transform"
	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/project"
)

type Override struct {
	Definition map[string][]string
}

type Options struct {
	Debug    bool
	Override Override
}

func (o Options) validate() error {
	return nil
}

type result struct {
	Runtime map[string]symbol.Symbol `json:"runtime" yaml:"runtime"`
	Types   map[string]symbol.Symbol `json:"types" yaml:"types"`
	Modules map[string]symbol.Symbol `json:"modules" yaml:"modules"`
}

type extractor struct {
	options     Options
	extractions []*extraction
	nvim        *nvim.Nvim
	logger      *log.Logger
	result      *result
}

func (e *extractor) Nvim() *nvim.Nvim {
	return e.nvim
}

func (e *extractor) Result() *result {
	return e.result
}

func (e *extractor) Flush() {
	e.extractions = []*extraction{}
	e.result = &result{
		Runtime: map[string]symbol.Symbol{},
		Types:   map[string]symbol.Symbol{},
		Modules: map[string]symbol.Symbol{},
	}
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

	e.logger.Debugf("%+v", e.result)

	placeholder := symbol.NewSymbol(name, symbol.Metadata{}, symbol.Documentation{}, symbol.NewUnknown())

	switch kind {
	case target.TargetKindValue:
		_, hasSymbol := e.result.Runtime[name]
		if hasSymbol {
			e.logger.Infof("Skipping extraction of '%s' %s target: already processed", name, kind)
			return nil
		}
		e.result.Runtime[name] = placeholder
	case target.TargetKindType:
		_, hasSymbol := e.result.Types[name]
		if hasSymbol {
			e.logger.Infof("Skipping extraction of '%s' %s target: already processed", name, kind)
			return nil
		}
		e.result.Types[name] = placeholder
	case target.TargetKindModule:
		_, hasSymbol := e.result.Modules[name]
		if hasSymbol {
			e.logger.Infof("Skipping extraction of '%s' %s target: already processed", name, kind)
			return nil
		}
		e.result.Modules[name] = placeholder
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
		e.result.Runtime[name] = extractionResult
	case target.TargetKindType:
		e.result.Types[name] = extractionResult
	case target.TargetKindModule:
		e.result.Modules[name] = extractionResult
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
		targets:    []*target.Target{extractionTarget},
	}
	targetExtractionContext := extractionContext{
		extraction: targetExtraction,
	}

	targetExtraction.crawler = crawl.NewCrawler(&targetExtractionContext)
	targetExtraction.transformer = transform.NewTransformer(&targetExtractionContext)
	targetExtraction.logger = log.NewLogger(extractionTarget.Identifier())

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
		targets:    []*target.Target{childExtractionTarget},
	}
	targetExtractionContext := extractionContext{
		extraction: childExtraction,
	}

	childExtraction.crawler = crawl.NewCrawler(&targetExtractionContext)
	childExtraction.transformer = transform.NewTransformer(&targetExtractionContext)
	childExtraction.logger = log.NewLogger(childExtractionTarget.Identifier())

	return childExtraction
}

func (e *extractor) initLogger(_ Options) {
	e.logger = log.NewLogger("")
}

func (e *extractor) initNvim(options Options) error {
	configDir, err := project.GetConfigDir()
	tmpDir, err := project.GetTmpDir()

	if err != nil {
		return fmt.Errorf("Error reading project directories: %w", err)
	}

	nvimConfig := nvim.NewConfig(configDir)

	var client *nvim.Nvim

	if !options.Debug {
		client, err = nvim.New(nvimConfig)
	} else {
		client, err = nvim.New(
			nvimConfig,
			nvim.WithArguments(
				fmt.Sprintf("-V%d%s", 1, path.Join(tmpDir, "nvim.verbosefile")),
				"--listen", path.Join(tmpDir, "nvim.server.pipe"),
			),
		)
	}

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

func NewExtractor(options Options) (*extractor, error) {
	options.validate()

	xtractor := &extractor{
		options:     options,
		extractions: []*extraction{},
		result: &result{
			Runtime: map[string]symbol.Symbol{},
			Types:   map[string]symbol.Symbol{},
			Modules: map[string]symbol.Symbol{},
		},
	}

	xtractor.initLogger(options)

	err := xtractor.initNvim(options)

	if err != nil {
		return nil, err
	}

	return xtractor, nil
}
