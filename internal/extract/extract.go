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
	Target string
	Debug  bool
	logger string
}

func (o Options) validate() error {
	return nil
}

type extractionItem struct {
	Name   string
	source symbol.Source
	symbol symbol.Symbol
}

type state func(e *extractor, item extractionItem) (extractionItem, state, error)

func initItem(e *extractor, item extractionItem) (extractionItem, state, error) {
	_, hasSymbol := e.result.Runtime[item.Name]

	if hasSymbol {
		return item, nil, nil
	}

	e.result.Runtime[item.Name] = symbol.NewUnknown()

	return item, crawlItem, nil
}

func crawlItem(e *extractor, item extractionItem) (extractionItem, state, error) {
	source, err := e.crawler.SourceValue(item.Name)

	if err != nil {
		return item, nil, err
	}

	item.source = source

	return item, lexItem, nil
}

func lexItem(e *extractor, item extractionItem) (extractionItem, state, error) {
	sym, err := e.lexer.Lex(item.source)

	if err != nil {
		return item, nil, err
	}

	item.symbol = sym

	return item, storeItem, nil
}

func storeItem(e *extractor, item extractionItem) (extractionItem, state, error) {
	e.result.Runtime[item.Name] = item.symbol

	return item, nil, nil
}

type result struct {
	Runtime map[string]symbol.Symbol `json:"runtime" yaml:"runtime"`
	Types   map[string]symbol.Symbol `json:"types" yaml:"types"`
}

type extractor struct {
	items   []extractionItem
	nvim    *nvim.Nvim
	crawler *crawl.Crawler
	lexer   *lex.Lexer
	logger  *log.Logger
	result  *result
}

func (e *extractor) current() extractionItem {
	return e.items[len(e.items)-1]
}

func (e *extractor) CurrentName() string {
	return e.current().Name
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

func (e *extractor) run(item extractionItem) (extractionItem, error) {
	var err error
	current := initItem
	for {
		item, current, err = current(e, item)
		if err != nil {
			return item, err
		}
		if current == nil {
			return item, nil
		}
	}
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

	e.nvim = client

	return nil
}

func (e *extractor) initCrawler(_ Options) {
	e.crawler = crawl.NewCrawler(e)
}

func (e *extractor) initLexer(_ Options) {
	e.lexer = lex.NewLexer(e)
}

func Extract(options Options) (result, error) {
	options.validate()

	xtractor := &extractor{
		items: []extractionItem{},
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

	_, err = xtractor.run(extractionItem{Name: options.Target})

	if err != nil {
		return *xtractor.result, err
	}

	return *xtractor.result, nil
}
