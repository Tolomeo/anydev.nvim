package crawler

import (
	"errors"
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type statistics struct {
	locationNotFound      map[string]struct{}
	originNotFound        map[string]struct{}
	documentationNotFound map[string]struct{}
}

func (s *statistics) LocationNotFound(path string) {
	s.locationNotFound[path] = struct{}{}
}

func (s *statistics) OriginNotFound(path string) {
	s.originNotFound[path] = struct{}{}
}

func (s *statistics) DocumentationNotFound(path string) {
	s.documentationNotFound[path] = struct{}{}
}

type crawler struct {
	nvim   *nvim.Nvim
	buffer string
	stats  *statistics
}

func (c *crawler) Statistics() statistics {
	return *c.stats
}

func (c *crawler) scratch(lines []string) error {
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

func (c *crawler) Crawl(path string) (Symbol, error) {
	foundSymbol, err := c.get(path)

	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", foundSymbol)

	return foundSymbol, nil
}

func (c *crawler) get(path string) (Symbol, error) {
	luaType, err := c.nvim.GetLuaTypeName(path)

	if err != nil {
		return nil, err
	}

	fmt.Println("findRuntimeSymbol", path, luaType)

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

func (c *crawler) GetVariableOrigin(path string) (origin, error) {
	var def = origin{}

	pathLocation, err := c.getLocation(path)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		c.stats.LocationNotFound(path)
		c.stats.OriginNotFound(path)
		c.stats.DocumentationNotFound(path)
		return def, nil
	case err != nil:
		return def, err
	default:
		def.SetLocation(
			pathLocation.Url,
			uint(pathLocation.TargetRange.Start.Line),
			uint(pathLocation.TargetRange.Start.Character),
		)
	}

	_, err = c.nvim.Open(pathLocation.Url)

	if err != nil {
		return def, err
	}

	variablePosition := nvim.CursorPosition{
		Line:      uint(pathLocation.TargetRange.Start.Line),
		Character: uint(pathLocation.TargetRange.Start.Character),
	}
	variableBufferLines, err := c.nvim.ReadTSNodeAt(variablePosition, nvim.TS_ASSIGNMENT_STATEMENT)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		c.stats.OriginNotFound(path)
		c.stats.DocumentationNotFound(path)
		return def, nil
	case err != nil:
		return def, err
	default:
		def.definition = variableBufferLines
	}

	documentationPosition := nvim.CursorPosition{
		Line:      max(0, variablePosition.Line-1),
		Character: uint(0),
	}
	documentationBufferLines, err := c.nvim.ReadCommentBlockAt(documentationPosition)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		c.stats.DocumentationNotFound(path)
		return def, nil
	case err != nil:
		return def, err
	default:
		def.documentation = documentationBufferLines
	}

	return def, nil

}

func (c *crawler) GetFunctionOrigin(path string) (origin, error) {
	var def = origin{}

	pathLocation, err := c.getLocation(path)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		c.stats.LocationNotFound(path)
		c.stats.OriginNotFound(path)
		c.stats.DocumentationNotFound(path)
		return def, nil
	case err != nil:
		return def, err
	default:
		def.SetLocation(
			pathLocation.Url,
			uint(pathLocation.TargetRange.Start.Line),
			uint(pathLocation.TargetRange.Start.Character),
		)
	}

	_, err = c.nvim.Open(pathLocation.Url)

	if err != nil {
		return def, err
	}

	functionPosition := nvim.CursorPosition{
		Line:      uint(pathLocation.TargetRange.Start.Line),
		Character: uint(pathLocation.TargetRange.Start.Character),
	}
	functionBufferLines, err := c.nvim.ReadTSNodeAt(functionPosition, nvim.TS_FUNCTION_DECLARATION, nvim.TS_ASSIGNMENT_STATEMENT)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		c.stats.OriginNotFound(path)
		c.stats.DocumentationNotFound(path)
		return def, nil
	case err != nil:
		return def, err
	default:
		def.definition = functionBufferLines
	}

	documentationPosition := nvim.CursorPosition{
		Line:      max(0, functionPosition.Line-1),
		Character: functionPosition.Character,
	}
	documentationBufferLines, err := c.nvim.ReadCommentBlockAt(documentationPosition)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		c.stats.DocumentationNotFound(path)
		return def, nil
	case err != nil:
		return def, err
	default:
		def.documentation = documentationBufferLines
	}

	return def, nil
}

func (c *crawler) getLocation(path string) (nvim.Location, error) {
	lines := []string{"local ref = " + path}

	err := c.scratch(lines)

	if err != nil {
		return nvim.Location{}, err
	}

	line, character := uint(0), uint(len(lines[0]))

	location, err := c.nvim.GetDefinitionLocation(line, character)

	if err != nil {
		return nvim.Location{}, err
	}

	return location, nil
}

func (c *crawler) GetFields(path string) ([]string, error) {
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

func New(nvim *nvim.Nvim) *crawler {
	scratchBuffer := nvim.Options().Config().Dir() + "anydev.lua"

	instance := crawler{
		nvim:   nvim,
		buffer: scratchBuffer,
		stats:  &statistics{
			locationNotFound: map[string]struct{}{},
			originNotFound: map[string]struct{}{},
			documentationNotFound: map[string]struct{}{},
		},
	}

	return &instance
}
