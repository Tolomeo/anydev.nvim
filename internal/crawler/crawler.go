package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/lsp"
	"github.com/Tolomeo/anydev.nvim/internal/slicesx"
)

var ErrNotFound = errors.New("Not found")

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
		fmt.Println("function", path)
		return struct{}{}, nil
	}

	return struct{}{}, nil
}

type definition struct {
	definition    []string
	documentation []string
}

func (c *crawler) GetDefinition(path string) (definition, error) {
	var def = definition{}

	location, err := c.GetDefinitionLocation(path)

	switch {
	case errors.Is(err, ErrNotFound):
		fmt.Printf("Documentation not found for %s\n", path)
		return def, nil
	case err != nil:
		return def, err
	}

	url, err := url.Parse(string(location.TargetUri))

	if err != nil {
		return def, err
	}

	_, err = c.nvim.Open(url.Path)

	if err != nil {
		return def, err
	}

	definitionPosition := nvim.CursorPosition{
		Line:      uint(location.TargetRange.Start.Line),
		Character: uint(location.TargetRange.Start.Character),
	}

	definitionBufferLines, err := c.nvim.GetTSAssignmentBufferLines(definitionPosition)

	if err != nil {
		return def, err
	}

	def.definition = definitionBufferLines

	documentationPosition := nvim.CursorPosition{
		Line:      uint(max(0, location.TargetRange.Start.Line-1)),
		Character: uint(0),
	}

	documentationBufferLines, err := c.nvim.GetTSCommentBlockBufferLines(documentationPosition)

	switch {
	case errors.Is(err, ErrNotFound):
		fmt.Printf("Documentation not found for %s\n", path)
		return def, nil
	case err != nil:
		return def, err
	}

	def.documentation = documentationBufferLines

	return def, nil
}

func (c *crawler) GetDefinitionLocation(path string) (lsp.DefinitionLocation, error) {
	lines := []string{"local ref = " + path}

	err := c.scratch(lines)

	if err != nil {
		return lsp.DefinitionLocation{}, err
	}

	line, character := uint(0), uint(len(lines[0]))

	lspDefinitionResponse, err := c.nvim.GetLSPDefinition(line, character)

	if err != nil {
		return lsp.DefinitionLocation{}, err
	}

	if len(lspDefinitionResponse.Result) == 0 {
		return lsp.DefinitionLocation{}, ErrNotFound
	}

	return lspDefinitionResponse.Result[0], nil
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

func (c *crawler) GetChildren(path string) ([]string, error) {
	err := c.nvim.StartLSP()

	if err != nil {
		return []string{}, fmt.Errorf("Error getting path %s children: %w", path, err)
	}

	completionPath := "lua " + path + "."
	getcompletionResult, err := c.nvim.CallFunction("getcompletion", []any{completionPath, "cmdline"})

	if err != nil {
		return []string{}, fmt.Errorf("Error getting path %s children: %w", path, err)
	}

	result, ok := slicesx.AnyToString(getcompletionResult.([]any))

	if !ok {
		return []string{}, fmt.Errorf("Error getting path %s children", path)
	}

	return result, nil
}

func (c *crawler) Debug(path string, subpath string) error {
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
}

func New(nvim *nvim.Nvim) *crawler {
	scratchBuffer := nvim.Options().Config().Dir() + "anydev.lua"

	instance := crawler{
		nvim:          nvim,
		scratchBuffer: scratchBuffer,
	}

	return &instance
}
