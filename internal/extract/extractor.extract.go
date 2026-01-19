package extract

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/lex/symbol"
)

type Target struct {
	kind   symbol.TargetKind
	parent symbol.Target
	name   string
	source symbol.Source
	symbol symbol.Symbol
}

func (e Target) Kind() symbol.TargetKind {
	return e.kind
}

func (e Target) ParentName() string {
	if e.parent == nil {
		return ""
	}

	return e.parent.Name()
}

func (e Target) Name() string {
	return e.name
}

func (e Target) Identifier() string {
	if e.parent == nil {
		return e.name
	}

	return fmt.Sprintf("%s.%s", e.parent.Identifier(), e.name)
}

func (e Target) Origin() symbol.Origin {
	return e.source.GetOrigin()
}

type state func(e *extractor, item Target) (Target, state, error)

func extractTarget(e *extractor, item Target) (Target, error) {
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

func crawlTarget(e *extractor, item Target) (Target, state, error) {
	fmt.Printf("\nCrawling '%s'\n", item.Identifier())

	var err error
	var source symbol.Source

	switch item.Kind() {
	case symbol.TargetKindValue:
		source, err = e.crawler.SourceValue()
	case symbol.TargetKindType:
		source, err = e.crawler.SourceType()
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

func lexTarget(e *extractor, item Target) (Target, state, error) {
	fmt.Printf("\nLexing '%s'\n", item.Identifier())

	sym, err := e.lexer.Lex()

	if err != nil {
		return item, nil, err
	}

	item.symbol = sym

	return item, nil, nil
}

func (e *extractor) ExtractChild(name string) (symbol.Symbol, error) {
	parent := e.Target()

	fmt.Printf("\nBeginning the extraction of '%s' . '%s' %s child target\n", parent.Name(), name, parent.Kind())

	childTarget := Target{parent: parent, kind: parent.Kind(), name: name}

	childTarget, err := extractTarget(e, childTarget)

	fmt.Printf("\nThe extraction of '%s' . '%s' %s child target yielded <%v>\n", parent.Name(), name, parent.Kind(), childTarget.symbol)

	if err != nil {
		return nil, err
	}

	if childTarget.symbol == nil {
		return symbol.NewUnknown(), nil
	}

	return childTarget.symbol, nil
}

func (e *extractor) Extract(kind symbol.TargetKind, name string) error {
	fmt.Printf("\nBeginning the extraction of '%s' %s target\n", name, kind)


	switch kind {
	case symbol.TargetKindValue:
		_, hasSymbol := e.result.Runtime[name]
		if hasSymbol {
			return nil
		}
		e.result.Runtime[name] = symbol.NewUnknown()
	case symbol.TargetKindType:
		_, hasSymbol := e.result.Types[name]
		if hasSymbol {
			return nil
		}
		e.result.Types[name] = symbol.NewUnknown()
	}

	target, err := extractTarget(e, Target{kind: kind, name: name})

	if err != nil {
		return err
	}

	fmt.Printf("\nThe extraction of '%s' %s target yielded <%v>\n", name, kind, target.symbol)

	if target.symbol == nil {
		return nil
	}

	switch target.kind {
	case symbol.TargetKindValue:
		e.result.Runtime[target.name] = target.symbol
	case symbol.TargetKindType:
		e.result.Types[target.name] = target.symbol
	}

	return nil
}
