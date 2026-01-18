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
	Source symbol.Source
	Debug  bool
	logger string
}

func (o Options) validate() error {
	return nil
}

type state func(e *extractor, src symbol.Source) (state, error)

func crawlSource(e *extractor, src symbol.Source) (state, error) {
	return nil, fmt.Errorf("Not implemented")
}

func lexSource(e *extractor, src symbol.Source) (state, error) {
	return nil, fmt.Errorf("Not implemented")
}

func run(e *extractor, source symbol.Source) (symbol.Source, error) {
	var err error
	current := crawlSource
	for {
		current, err = current(e, source)
		if err != nil {
			return source, err
		}
		if current == nil {
			return source, nil
		}
	}
}

type result struct {
	Runtime map[string]symbol.Symbol `json:"runtime" yaml:"runtime"`
	Types   map[string]symbol.Symbol `json:"types" yaml:"types"`
}

type extractor struct {
	source  []symbol.Source
	nvim    *nvim.Nvim
	crawler *crawl.Crawler
	lexer   *lex.Lexer
	logger  *log.Logger
	result  *result
}

func (e *extractor) Source() symbol.Source {
	return e.source[len(e.source)-1]
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

	return nil
}

func (e *extractor) initCrawler(_ Options) {
	e.crawler = crawl.NewCrawler(e)
}

func (e *extractor) initLexer(_ Options) {
	e.lexer = lex.NewLexer()
}

func Extract(options Options) (result, error) {
	options.validate()

	xtractor := &extractor{
		source: []symbol.Source{options.Source},
		result: &result{
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

	run(xtractor, xtractor.Source())

	return *xtractor.result, nil
}
