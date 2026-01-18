package extract

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

type targetKind string

const (
	TargetKindValue targetKind = "value"
	TargetKindType  targetKind = "type"
)

type extractionTarget struct {
	kind   targetKind
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

func extractTarget(e *extractor, item extractionTarget) (extractionTarget, error) {
	fmt.Printf("\nExtracting '%s'\n", item.Identifier())

	e.targets = append(e.targets, item)

	var err error
	current := crawlTarget

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

func crawlTarget(e *extractor, item extractionTarget) (extractionTarget, state, error) {
	fmt.Printf("\nCrawling '%s'\n", item.Identifier())

	var err error
	var source symbol.Source

	switch item.kind {
	case TargetKindValue:
		source, err = e.crawler.SourceValue(item.Identifier())
	case TargetKindType:
		source, err = e.crawler.SourceType(item.Name(), item.ParentName())
	}

	if err != nil {
		return item, nil, err
	}

	if source == nil {
		return item, nil, nil
	}

	item.source = source

	return item, lexTarget, nil
}

func lexTarget(e *extractor, item extractionTarget) (extractionTarget, state, error) {
	fmt.Printf("\nLexing '%s'\n", item.Identifier())

	sym, err := e.lexer.Lex(item.source)

	if err != nil {
		return item, nil, err
	}

	item.symbol = sym

	return item, nil, nil
}

func (e *extractor) ExtractChild(name string) (symbol.Symbol, error) {
	parent := e.current()

	fmt.Printf("\nBeginning the extraction of '%s' . '%s' %s child target\n", parent.name, name, parent.kind)

	childTarget := extractionTarget{parent: &parent, kind: parent.kind, name: name}

	childTarget, err := extractTarget(e, childTarget)

	fmt.Printf("\nThe extraction of '%s' . '%s' %s child target yielded <%v>\n", parent.name, name, parent.kind, childTarget.symbol)

	if err != nil {
		return nil, err
	}

	if childTarget.symbol == nil {
		return symbol.NewUnknown(), nil
	}

	return childTarget.symbol, nil
}

func (e *extractor) Extract(kind string, name string) error {
	fmt.Printf("\nBeginning the extraction of '%s' %s target\n", name, kind)

	extractionTargetKind := targetKind(kind)

	switch extractionTargetKind {
	case TargetKindValue:
		_, hasSymbol := e.result.Runtime[name]
		if hasSymbol {
			return nil
		}
		e.result.Runtime[name] = symbol.NewUnknown()
	case TargetKindType:
		_, hasSymbol := e.result.Types[name]
		if hasSymbol {
			return nil
		}
		e.result.Types[name] = symbol.NewUnknown()
	}

	target, err := extractTarget(e, extractionTarget{kind: extractionTargetKind, name: name})

	if err != nil {
		return err
	}

	fmt.Printf("\nThe extraction of '%s' %s target yielded <%v>\n", name, kind, target.symbol)

	if target.symbol == nil {
		return nil
	}

	switch target.kind {
	case TargetKindValue:
		e.result.Runtime[target.name] = target.symbol
	case TargetKindType:
		e.result.Types[target.name] = target.symbol
	}

	return nil
}
