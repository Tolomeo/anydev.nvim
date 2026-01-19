package extract

import (
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/lex"
	"github.com/Tolomeo/anydev.nvim/internal/lex/crawl"
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
	targets []Target
	nvim    *nvim.Nvim
	crawler *crawl.Crawler
	lexer   *lex.Lexer
	logger  *log.Logger
	result  *extraction
}

func (e *extractor) Target() symbol.Target {
	return e.targets[len(e.targets)-1]
}

func (e *extractor) CurrentName() string {
	return e.Target().Name()
}

func (e *extractor) Nvim() *nvim.Nvim {
	return e.nvim
}

func (e *extractor) Result() *extraction {
	return e.result
}

func (e *extractor) Flush() {
	e.targets = []Target{}
	e.result = &extraction{
		Runtime: map[string]symbol.Symbol{},
		Types:   map[string]symbol.Symbol{},
	}
}

func (e *extractor) Logger() *log.Logger {
	return e.logger
}

func (e *extractor) initLogger(_ Options) {
	e.logger = log.NewLogger()
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

func NewExtractor(options Options) (*extractor, error) {
	options.validate()

	xtractor := &extractor{
		targets: []Target{},
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
		return nil, err
	}

	return xtractor, nil
}
