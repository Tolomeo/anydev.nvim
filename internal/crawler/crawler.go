package crawler

import (
	"errors"
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type Crawler struct {
	nvim   *nvim.Nvim
	buffer string
	stats  *statistics
}

func (c *Crawler) Statistics() statistics {
	return *c.stats
}

func (c *Crawler) scratch(lines []string) error {
	_, err := c.nvim.Open(c.buffer)

	if err != nil {
		return err
	}

	err = c.nvim.SetBufferLines(lines)

	if err != nil {
		return err
	}

	return nil
}

func (c *Crawler) Crawl(path string) (Source, error) {
	foundSymbol, err := c.get(path)

	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", foundSymbol)

	return foundSymbol, nil
}

func (c *Crawler) get(path string) (Source, error) {
	luaType, err := c.nvim.GetLuaTypeName(path)

	if err != nil {
		return nil, err
	}

	switch luaType {
	case "table":
		namespace, err := NewNamespace(c, path)

		if err != nil {
			return nil, err
		}

		return namespace, nil

	case "function":
		function, err := NewFunction(c, path)

		if err != nil {
			return nil, err
		}

		return function, nil
	case "boolean", "number", "string", "userdata", "thread", "nil":
		variable, err := NewVariable(c, path)

		if err != nil {
			return nil, err
		}

		return variable, nil
	}

	return nil, fmt.Errorf("Unrecognized type '%s' received for path '%s'", luaType, path)
}

func (c *Crawler) GetOrigin(path string, fieldType string, fieldTypes ...string) (origin, error) {
	var orig = origin{}

	locations, err := c.getLocation(path)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		c.stats.Report(path, c.stats.MissingLocation(), c.stats.MissingOrigin(), c.stats.MissingDocumentation())
		return orig, nil
	case err != nil:
		return orig, err
	}

	locationIndex, err := slicesx.IndexFunc(locations, func(location nvim.Location) (bool, error) {
		_, err = c.nvim.Open(location.Url)

		if err != nil {
			return false, err
		}

		position := nvim.CursorPosition{
			Line:      uint(location.TargetRange.Start.Line),
			Character: uint(location.TargetRange.Start.Character),
		}
		bufferLines, err := c.nvim.ReadTSNodeAt(position, fieldType, fieldTypes...)

		switch {
		case errors.Is(err, nvim.ErrNotFound):
			return false, nil
		case err != nil:
			return false, err
		}

		orig.definition = bufferLines
		orig.SetLocation(
			location.Url,
			uint(location.TargetRange.Start.Line),
			uint(location.TargetRange.Start.Character),
		)

		return true, nil
	})

	switch {
	case err != nil:
		return orig, err
	case locationIndex == -1:
		c.stats.Report(path, c.stats.MissingOrigin(), c.stats.MissingDocumentation())
		return orig, nil
	}

	functionLocation := locations[locationIndex]

	documentationPosition := nvim.CursorPosition{
		Line:      uint(max(0, functionLocation.TargetRange.Start.Line-1)),
		Character: uint(functionLocation.TargetRange.Start.Character),
	}
	documentationBufferLines, err := c.nvim.ReadCommentBlockAt(documentationPosition)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		c.stats.Report(path, c.stats.MissingDocumentation())
		return orig, nil
	case err != nil:
		return orig, err
	default:
		orig.documentation = documentationBufferLines
	}

	return orig, nil
}

func (c *Crawler) getLocation(path string) ([]nvim.Location, error) {
	lines := []string{"local ref = " + path}

	err := c.scratch(lines)

	if err != nil {
		return []nvim.Location{}, err
	}

	line, character := uint(0), uint(len(lines[0]))

	locations, err := c.nvim.GetDefinitionLocation(line, character)

	if err != nil {
		return []nvim.Location{}, err
	}

	return locations, nil
}

func (c *Crawler) GetFields(path string) ([]string, error) {
	return c.nvim.GetCompletion(path)
}

/* func (c *crawler) Debug(path string, subpath string) error {
	statement := "local ref = " + path + "." + subpath

	err := c.nvim.SetBufferLines([]string{
		statement,
	})

	if err != nil {
		return err
	}

	lines, err := c.nvim.GetBufferLines()

	if err != nil {
		return err
	}

	fmt.Println(lines)

	documentSymbols, err := c.nvim.GetLSPDocumentSymbols()

	if err != nil {
		return err
	}

	documentSymbol, ok := slicesx.FindFunc(documentSymbols.Result, func(s lsp.DocumentSymbol) bool {
		return s.Name == "ref"
	})

	if !ok {
		return fmt.Errorf("Error retrieving ref from document symbols: %+v", documentSymbols)
	}

	fmt.Println("Hover")
	c.nvim.GetLSPHover(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("Declaration")
	c.nvim.GetLSPDeclaration(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("Definition")
	c.nvim.GetLSPDefinition(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("TypeDefinition")
	c.nvim.GetLSPTypeDefinition(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))
	fmt.Println("Implementation")
	c.nvim.GetLSPImplementation(uint(documentSymbol.Range.End.Line), uint(documentSymbol.Range.End.Character))

	return nil
} */

func New(nvim *nvim.Nvim) *Crawler {
	buffer := nvim.Options().Config().Dir() + "anydev.lua"

	instance := Crawler{
		nvim:   nvim,
		buffer: buffer,
		stats:  NewStatistics(),
	}

	return &instance
}
