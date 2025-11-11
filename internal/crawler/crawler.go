package crawler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type crawler struct {
	nvim          *nvim.Nvim
	scratchBuffer string
}

func (c *crawler) scratch(lines []string) error {
	_, err := c.nvim.Open(c.scratchBuffer)

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
	runtimeType, err := c.getRuntimeTypeName(path)

	if err != nil {
		return nil, err
	}

	fmt.Println("findRuntimeSymbol", path, runtimeType)

	switch runtimeType {
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
	}

	return struct{}{}, nil
}

func (c *crawler) GetAssignmentOrigin(path string) (origin, error) {
	var def = origin{}

	pathLocation, err := c.getLocation(path)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		fmt.Printf("Location not found for %s\n", path)
		return def, nil
	case err != nil:
		return def, err
	}

	_, err = c.nvim.Open(pathLocation.Url)

	if err != nil {
		return def, err
	}

	assignmentPosition := nvim.CursorPosition{
		Line:      uint(pathLocation.TargetRange.Start.Line),
		Character: uint(pathLocation.TargetRange.Start.Character),
	}

	definitionBufferLines, err := c.nvim.GetAssignmentStatementAt(assignmentPosition)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		fmt.Printf("Definition not found for %s\n", path)
		return def, nil
	case err != nil:
		return def, err
	default:
		def.definition = definitionBufferLines
	}

	documentationPosition := nvim.CursorPosition{
		Line:      max(0, assignmentPosition.Line-1),
		Character: assignmentPosition.Character,
	}
	documentationBufferLines, err := c.nvim.GetCommentBlockAt(documentationPosition)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		fmt.Printf("Documentation not found for %s\n", path)
		return def, nil
	case err != nil:
		return def, err
	default:
		def.documentation = documentationBufferLines
	}

	return def, nil

}

func (c *crawler) GetDeclarationOrigin(path string) (origin, error) {
	var def = origin{}

	pathLocation, err := c.getLocation(path)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		fmt.Printf("Location not found for %s\n", path)
		return def, nil
	case err != nil:
		return def, err
	}

	_, err = c.nvim.Open(pathLocation.Url)

	if err != nil {
		return def, err
	}

	declarationPosition := nvim.CursorPosition{
		Line:      uint(pathLocation.TargetRange.Start.Line),
		Character: uint(pathLocation.TargetRange.Start.Character),
	}
	declarationBufferLines, err := c.nvim.GetDeclarationStatementAt(declarationPosition)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		fmt.Printf("Declaration statement not found for %s\n", path)
		return def, nil
	case err != nil:
		return def, err
	default:
		def.definition = declarationBufferLines
	}

	documentationPosition := nvim.CursorPosition{
		Line:      max(0, declarationPosition.Line-1),
		Character: declarationPosition.Character,
	}
	documentationBufferLines, err := c.nvim.GetCommentBlockAt(documentationPosition)

	switch {
	case errors.Is(err, nvim.ErrNotFound):
		fmt.Printf("Documentation not found for %s\n", path)
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

func (c *crawler) getRuntimeTypeName(path string) (string, error) {
	runtimePath := path
	parts := strings.Split(runtimePath, ".")

	switch len(parts) {
	case 1:
	default:
		tail := parts[len(parts)-1]
		// https://www.lua.org/manual/5.1/manual.html#2.1
		switch tail {
		case "and", "break", "do", "else", "elseif", "end", "false", "for", "function", "if", "in", "local", "nil", "not", "or", "repeat", "return", "then", "true", "until", "while":
			head := parts[:len(parts)-1]
			runtimePath = strings.Join(head, ".") + "['" + tail + "']"
		}
	}

	luaCode := fmt.Sprintf("return type(%s)", runtimePath)

	result, err := c.nvim.ExecLua(luaCode, []any{})

	if err != nil {
		return "", fmt.Errorf("Error getting the type of %s: %w", path, err)
	}

	typeName, ok := result.(string)

	if !ok {
		return "", fmt.Errorf("Error getting the type of %s: Error converting the result to a string", runtimePath)
	}

	return typeName, nil
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
		nvim:          nvim,
		scratchBuffer: scratchBuffer,
	}

	return &instance
}
