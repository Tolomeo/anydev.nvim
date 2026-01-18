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
	parent *extractionItem
	name   string
	source symbol.Source
	symbol symbol.Symbol
}

func (e extractionItem) Name() string {
	return e.name
}

func (e extractionItem) Identifier() string {
	if e.parent == nil {
		return e.name
	}

	return fmt.Sprintf("%s.%s", e.parent.Identifier(), e.name)
}

type state func(e *extractor, item extractionItem) (extractionItem, state, error)

func extract(e *extractor, item extractionItem) (extractionItem, error) {
	e.items = append(e.items, item)

	var err error
	current := crawlItem
	for {
		item, current, err = current(e, item)
		if err != nil {
			return item, err
		}
		if current == nil {
			e.items = e.items[:len(e.items)-1]
			return item, nil
		}
	}
}

func crawlItem(e *extractor, item extractionItem) (extractionItem, state, error) {
	source, err := e.crawler.SourceValue(item.Identifier())

	if err != nil {
		return item, nil, err
	}

	if source == nil {
		return item, nil, nil
	}

	item.source = source

	return item, lexItem, nil
}

func lexItem(e *extractor, item extractionItem) (extractionItem, state, error) {
	sym, err := e.lexer.Lex(item.source)

	if err != nil {
		return item, nil, err
	}

	switch s := sym.(type) {
	case *symbol.Table:
		children, err := e.nvim.GetValueCompletion(item.name)

		if err != nil {
			return item, nil, err
		}

		for _, child := range children {
			childItem := extractionItem{parent: &item, name: child}
			childItem, err := extract(e, childItem)

			if err != nil {
				return item, nil, err
			}

			s.Fields = append(s.Fields, symbol.TableField{Name: child, Value: childItem.symbol})
		}

		fmt.Printf("table point %v", s)
	}

	item.symbol = sym

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
	return e.current().name
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

func (e *extractor) run(item extractionItem) error {
	_, hasSymbol := e.result.Runtime[item.name]

	if hasSymbol {
		return nil
	}

	e.result.Runtime[item.name] = symbol.NewUnknown()

	item, err := extract(e, item)

	if err != nil {
		return err
	}

	if item.symbol == nil {
		return nil
	}

	e.result.Runtime[item.name] = item.symbol

	return nil
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

	err = xtractor.run(extractionItem{name: options.Target})

	if err != nil {
		return *xtractor.result, err
	}

	err = xtractor.nvim.Quit()

	if err != nil {
		return *xtractor.result, err
	}

	return *xtractor.result, nil
}
