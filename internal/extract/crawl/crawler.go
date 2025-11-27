package crawl

import (
	"errors"
	"fmt"
	"path"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/ts"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type Crawler struct {
	nvim  *nvim.Nvim
	stats *statistics
}

func (c *Crawler) Statistics() statistics {
	return *c.stats
}

func (c *Crawler) scratch(lines []string) error {
	buffer := path.Join(c.nvim.Options().Config().Dir(), "anydev.crawler.lua")

	_, err := c.nvim.Open(buffer)

	if err != nil {
		return err
	}

	err = c.nvim.SetBufferLines(lines)

	if err != nil {
		return err
	}

	return nil
}

func (c *Crawler) CrawlRuntime(path string) (Source, error) {
	source, err := c.getRuntimeSource(path)

	if err != nil {
		return nil, err
	}

	return source, nil
}

func (c *Crawler) getRuntimeSource(path string) (Source, error) {
	luaType, err := c.nvim.GetLuaTypeName(path)

	if err != nil {
		return nil, err
	}

	switch luaType {
	case "table":
		tableOrigin, err := c.getOrigin(path, ts.ASSIGNMENT_STATEMENT)

		if err != nil {
			return nil, err
		}

		tableSource := TableSource{
			path:   path,
			origin: tableOrigin,
		}

		fields, err := c.nvim.GetCompletion(path)

		if err != nil {
			return nil, err
		}

		for _, field := range fields {
			child, err := c.CrawlRuntime(path + "." + field)

			if err != nil {
				return nil, err
			}

			tableSource.fields = append(tableSource.fields, &child)
		}

		return &tableSource, nil

	case "function":
		functionOrigin, err := c.getOrigin(path, ts.FUNCTION_DECLARATION, ts.ASSIGNMENT_STATEMENT)

		if err != nil {
			return nil, err
		}

		return &FunctionSource{
			path:   path,
			origin: functionOrigin,
		}, nil

	case "boolean", "number", "string", "userdata", "thread", "nil":
		variableOrigin, err := c.getOrigin(path, ts.ASSIGNMENT_STATEMENT)

		if err != nil {
			return nil, err
		}

		return &VariableSource{
			path:   path,
			origin: variableOrigin,
		}, nil
	}

	return nil, fmt.Errorf("Unrecognized type '%s' received for path '%s'", luaType, path)
}

func (c *Crawler) getOrigin(path string, fieldType string, fieldTypes ...string) (*origin, error) {
	var orig = origin{}

	locations, err := c.getLocation(path)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		c.stats.Report(path, c.stats.MissingLocation(), c.stats.MissingOrigin(), c.stats.MissingDocumentation())
		return nil, nil
	case err != nil:
		return nil, err
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
		definitionLines, err := c.nvim.ReadTSNodeAt(position, fieldType, fieldTypes...)

		switch {
		case errors.Is(err, nvim.ErrNotFound):
			return false, nil
		case err != nil:
			return false, err
		}

		orig.Url = location.Url
		orig.Line = uint(location.TargetRange.Start.Line)
		orig.Character = uint(location.TargetRange.Start.Character)
		orig.Definition = definitionLines

		return true, nil
	})

	switch {
	case err != nil:
		return nil, err
	case locationIndex == -1:
		c.stats.Report(path, c.stats.MissingOrigin(), c.stats.MissingDocumentation())
		return nil, nil
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
		return &orig, nil
	case err != nil:
		return &orig, err
	default:
		orig.Documentation = documentationBufferLines
	}

	return &orig, nil
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

func NewCrawler(nvim *nvim.Nvim) *Crawler {
	instance := Crawler{
		nvim:  nvim,
		stats: NewStatistics(),
	}

	return &instance
}
