package extract

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/project"
)

type Options struct {
	Debug bool
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
	}
}

func (e *extractor) extract(item *extraction) error {
	e.logger.Infof("Extracting '%s'", item.target.Identifier())

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

	e.logger.Infof("The extraction of '%s' %s target yielded \n<%v>", name, kind, extraction.target.Type())

	if extraction.target.Type() == nil {
		return nil
	}

	extractionResult := symbol.NewSymbol(name, extraction.target.Meta(), extraction.target.Documentation(), extraction.target.Type())

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
