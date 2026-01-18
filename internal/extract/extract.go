package extract

import (
	"fmt"
	"path"
	"slices"

	"github.com/Tolomeo/anydev.nvim/internal/lex"
	"github.com/Tolomeo/anydev.nvim/internal/lex/crawl"
	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/project"
	"github.com/Tolomeo/anydev.nvim/internal/utils/log"
)

type Options struct {
	Kind  string
	Name  string
	Debug bool
}

func (o Options) validate() error {
	if o.Name == "" {
		return fmt.Errorf("Name option is required")
	}

	allowedKinds := []string{targetKindValue, targetKindType}

	if ok := slices.Contains(allowedKinds, o.Kind); !ok {
		return fmt.Errorf("Type options invalid: allowedTypes are <%s>", allowedKinds)
	}

	return nil
}

type extraction struct {
	Runtime map[string]symbol.Symbol `json:"runtime" yaml:"runtime"`
	Types   map[string]symbol.Symbol `json:"types" yaml:"types"`
}

type extractor struct {
	targets []extractionTarget
	nvim    *nvim.Nvim
	crawler *crawl.Crawler
	lexer   *lex.Lexer
	logger  *log.Logger
	result  *extraction
}

func (e *extractor) current() extractionTarget {
	return e.targets[len(e.targets)-1]
}

func (e *extractor) CurrentName() string {
	return e.current().Name()
}

func (e *extractor) CurrentSource() symbol.Source {
	return e.current().source
}

func (e *extractor) Nvim() *nvim.Nvim {
	return e.nvim
}

func (e *extractor) Result() *nvim.Nvim {
	return e.nvim
}

func (e *extractor) Logger() *log.Logger {
	return e.logger
}

func (e *extractor) Add() {

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

func (e *extractor) initCrawler(_ Options) {
	e.crawler = crawl.NewCrawler(e)
}

func (e *extractor) initLexer(_ Options) {
	e.lexer = lex.NewLexer(e)
}

func Extract(options Options) (extraction, error) {
	options.validate()

	xtractor := &extractor{
		targets: []extractionTarget{},
		result: &extraction{
			Runtime: map[string]symbol.Symbol{},
			Types:   map[string]symbol.Symbol{},
		},
	}

	xtractor.initLogger(options)

	xtractor.initCrawler(options)

	xtractor.initLexer(options)

	err := xtractor.initNvim(options)

	if err != nil {
		return *xtractor.result, err
	}

	err = xtractor.Extract(options.Kind, options.Name)

	if err != nil {
		return *xtractor.result, err
	}

	err = xtractor.nvim.Quit()

	if err != nil {
		return *xtractor.result, err
	}

	return *xtractor.result, nil
}
