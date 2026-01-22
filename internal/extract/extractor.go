package extract

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/project"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type Options struct {
	Debug bool
}

func (o Options) validate() error {
	return nil
}

type extraction struct {
	Runtime map[string]symbol.Symbol `json:"runtime" yaml:"runtime"`
	Types   map[string]symbol.Symbol `json:"types" yaml:"types"`
}

type extractor struct {
	targets []*target
	nvim    *nvim.Nvim
	logger  *log.Logger
	result  *extraction
}

func (e *extractor) Nvim() *nvim.Nvim {
	return e.nvim
}

func (e *extractor) Result() *extraction {
	return e.result
}

func (e *extractor) Flush() {
	e.targets = []*target{}
	e.result = &extraction{
		Runtime: map[string]symbol.Symbol{},
		Types:   map[string]symbol.Symbol{},
	}
}

func (e *extractor) extract(item *target) error {
	e.logger.Infof("Extracting '%s'", item.Identifier())

	e.targets = append(e.targets, item)

	var err error
	current := item.getOrigins

	for {
		current, err = current()

		e.targets[len(e.targets)-1] = item

		if err != nil {
			return err
		}

		if current == nil {
			e.targets = e.targets[:len(e.targets)-1]
			return nil
		}
	}
}

func (e *extractor) Extract(kind symbol.TargetKind, name string) error {
	e.logger.Infof("Beginning the extraction of '%s' %s target", name, kind)

	switch kind {
	case symbol.TargetKindValue:
		_, hasSymbol := e.result.Runtime[name]
		if hasSymbol {
			e.logger.Infof("Skipping extraction of '%s' %s target: already processed", name, kind)
			return nil
		}
		e.result.Runtime[name] = symbol.NewSymbol(name, symbol.Meta{}, symbol.Documentation{}, symbol.NewUnknown())
	case symbol.TargetKindType:
		_, hasSymbol := e.result.Types[name]
		if hasSymbol {
			e.logger.Infof("Skipping extraction of '%s' %s target: already processed", name, kind)
			return nil
		}
		e.result.Types[name] = symbol.NewSymbol(name, symbol.Meta{}, symbol.Documentation{}, symbol.NewUnknown())
	}

	target := e.newTarget(kind, name)

	err := e.extract(target)

	if err != nil {
		return err
	}

	e.logger.Infof("The extraction of '%s' %s target yielded \n<%v>", name, kind, target.type_)

	if target.type_ == nil {
		return nil
	}

	switch target.kind {
	case symbol.TargetKindValue:
		e.result.Runtime[target.name] = symbol.NewSymbol(name, target.meta, target.Origin().Documentation(), target.type_)
	case symbol.TargetKindType:
		e.result.Types[target.name] = symbol.NewSymbol(name, target.meta, target.Origin().Documentation(), target.type_)
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
				fmt.Sprintf("-V%d%s", 10, path.Join(tmpDir, "nvim.verbosefile")),
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
		targets: []*target{},
		result: &extraction{
			Runtime: map[string]symbol.Symbol{},
			Types:   map[string]symbol.Symbol{},
		},
	}

	xtractor.initLogger(options)

	err := xtractor.initNvim(options)

	if err != nil {
		return nil, err
	}

	return xtractor, nil
}
