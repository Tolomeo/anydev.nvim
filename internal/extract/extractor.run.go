package extract

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

const (
	targetKindValue string = "value"
	targetKindType  string = "type"
)

type extractionTarget struct {
	kind   string
	parent *extractionTarget
	name   string
	source symbol.Source
	symbol symbol.Symbol
}

func (e extractionTarget) ParentName() string {
	if e.parent == nil {
		return ""
	}

	return e.parent.Name()
}

func (e extractionTarget) Name() string {
	return e.name
}

func (e extractionTarget) Identifier() string {
	if e.parent == nil {
		return e.name
	}

	return fmt.Sprintf("%s.%s", e.parent.Identifier(), e.name)
}

type state func(e *extractor, item extractionTarget) (extractionTarget, state, error)

func extract(e *extractor, item extractionTarget) (extractionTarget, error) {
	e.targets = append(e.targets, item)

	var err error
	current := crawlItem
	for {
		item, current, err = current(e, item)

		e.targets[len(e.targets)-1] = item

		if err != nil {
			return item, err
		}

		if current == nil {
			e.targets = e.targets[:len(e.targets)-1]
			return item, nil
		}
	}
}

func crawlItem(e *extractor, item extractionTarget) (extractionTarget, state, error) {
	var err error
	var source symbol.Source

	switch item.kind {
	case targetKindValue:
		source, err = e.crawler.SourceValue(item.Identifier())
	case targetKindType:
		source, err = e.crawler.SourceType(item.Name(), item.ParentName())
	}

	if err != nil {
		return item, nil, err
	}

	if source == nil {
		return item, nil, nil
	}

	item.source = source

	return item, lexItem, nil
}

func lexItem(e *extractor, item extractionTarget) (extractionTarget, state, error) {
	sym, err := e.lexer.Lex(item.source)

	if err != nil {
		return item, nil, err
	}

	item.symbol = sym

	return item, nil, nil
}

func (e *extractor) ExtractChild(name string) (symbol.Symbol, error) {
	parent := e.current()
	childTarget := extractionTarget{parent: &parent, kind: parent.kind, name: name}

	childTarget, err := extract(e, childTarget)

	if err != nil {
		return nil, err
	}

	return childTarget.symbol, nil
}

func (e *extractor) Extract(kind string, name string) error {
	switch kind {
	case targetKindValue:
		_, hasSymbol := e.result.Runtime[name]
		if hasSymbol {
			return nil
		}
		e.result.Runtime[name] = symbol.NewUnknown()
	case targetKindType:
		_, hasSymbol := e.result.Types[name]
		if hasSymbol {
			return nil
		}
		e.result.Types[name] = symbol.NewUnknown()
	}

	target, err := extract(e, extractionTarget{kind: kind, name: name})

	if err != nil {
		return err
	}

	if target.symbol == nil {
		return nil
	}

	switch target.kind {
	case targetKindValue:
		e.result.Runtime[target.name] = target.symbol
	case targetKindType:
		e.result.Types[target.name] = target.symbol
	}

	return nil
}
