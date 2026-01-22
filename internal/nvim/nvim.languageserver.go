package nvim

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/internal/scripts"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/languageserver"
	"github.com/Tolomeo/anydev.nvim/internal/utils/anyx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

type Location struct {
	languageserver.Location
	Url string
}

func (l *Location) StartLine() uint {
	return uint(l.TargetRange.Start.Line)
}

func (l *Location) StartCharacter() uint {
	return uint(l.TargetRange.Start.Character)
}

func (n *Nvim) startLSP() error {
	script, err := scripts.Read("start-lsp")

	if err != nil {
		return fmt.Errorf("Error starting lua lsp: %v", err)
	}

	_, err = n.execLua(script, []any{30000})

	if err != nil {
		return fmt.Errorf("Error starting lua lsp: %v", err)
	}

	return nil
}

func (n *Nvim) getDocumentSymbols() (*languageserver.TextDocumentDocumentSymbolResponse, error) {
	documentSymbols := languageserver.TextDocumentDocumentSymbolResponse{}

	err := n.startLSP()

	if err != nil {
		return nil, err
	}

	script, err := scripts.Read("get-lsp-document-symbols")

	if err != nil {
		return nil, fmt.Errorf("Error getting document symbols: %v", err)
	}

	result, err := n.execLua(script, []any{})

	if err != nil {
		return nil, fmt.Errorf("Error getting document symbols: %v", err)
	}

	if result == nil {
		return nil, nil
	}

	stringResult, ok := result.(string)

	if !ok {
		return nil, fmt.Errorf("Error reading document symbols response: %v", result)
	}

	err = documentSymbols.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling document symbols response: %v", err)
	}

	return &documentSymbols, nil
}

func (n *Nvim) getHover(line uint, character uint) (*languageserver.TextDocumentHoverResponse, error) {
	hover := languageserver.TextDocumentHoverResponse{}

	err := n.startLSP()

	if err != nil {
		return nil, err
	}

	script, err := scripts.Read("get-lsp-hover")

	if err != nil {
		return nil, err
	}

	result, err := n.execLua(script, []any{line, character, 15000})

	if err != nil {
		return nil, fmt.Errorf("Error getting lsp hover response: %v", err)
	}

	if result == nil {
		return nil, nil
	}

	stringResult, ok := result.(string)

	if !ok {
		return nil, fmt.Errorf("Error reading lsp hover response: %v", result)
	}

	// fmt.Println(stringResult)

	err = hover.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling lsp hover response: %v", err)
	}

	return &hover, nil
}

func (n *Nvim) getLSPDefinitions(line uint, character uint) (*[]languageserver.Location, error) {
	err := n.startLSP()

	if err != nil {
		return nil, err
	}

	script, err := scripts.Read("get-lsp-definition-locations")

	if err != nil {
		return nil, err
	}

	result, err := n.execLua(script, []any{line, character, 15000})

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error getting lsp definition: %v", err)
	case result == nil:
		return nil, nil
	}

	stringResult, ok := result.(string)

	if !ok {
		return nil, fmt.Errorf("Error reading lsp definition response: %v", result)
	}

	response := languageserver.TextDocumentDefinitionResponse{}
	err = response.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling lsp definition response: %w", err)
	}

	return &response.Result, nil
}

func (n *Nvim) getLspTypeDefinitions(line uint, character uint) (*[]languageserver.Location, error) {
	err := n.startLSP()

	if err != nil {
		return nil, err
	}

	script, err := scripts.Read("get-lsp-type-definition-locations")

	if err != nil {
		return nil, err
	}

	result, err := n.execLua(script, []any{line, character, 15000})

	switch {
	case err != nil:
		return nil, fmt.Errorf("Error getting lsp type definition: %v", err)
	case result == nil:
		return nil, nil
	}

	stringResult, ok := result.(string)

	if !ok {
		return nil, fmt.Errorf("Error reading lsp definition response: %v", result)
	}

	response := languageserver.TextDocumentTypeDefinitionResponse{}
	err = response.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling lsp type definition response: %w", err)
	}

	return &response.Result, nil
}

func (n *Nvim) getTypeDefinitionLocations(line uint, character uint) (*[]Location, error) {
	lspTypeDefinitions, err := n.getLspTypeDefinitions(line, character)

	switch {
	case err != nil:
		return nil, err
	case lspTypeDefinitions == nil:
		return nil, nil
	case len(*lspTypeDefinitions) == 0:
		return nil, nil
	}

	locations, err := slicesx.MapFunc(*lspTypeDefinitions, func(lspLocation languageserver.Location) (Location, error) {
		location := Location{
			Location: lspLocation,
		}

		url, err := url.Parse(string(location.Location.TargetUri))

		if err != nil {
			return Location{}, err
		}

		location.Url = url.Path

		return location, nil
	})

	if err != nil {
		return nil, err
	}

	return &locations, nil
}

func (n *Nvim) getDefinitionLocations(line uint, character uint) (*[]Location, error) {
	lspDefinitions, err := n.getLSPDefinitions(line, character)

	switch {
	case err != nil:
		return nil, err
	case lspDefinitions == nil:
		return nil, nil
	case len(*lspDefinitions) == 0:
		return nil, nil
	}

	locations, err := slicesx.MapFunc(*lspDefinitions, func(lspLocation languageserver.Location) (Location, error) {
		location := Location{
			Location: lspLocation,
		}

		url, err := url.Parse(string(location.Location.TargetUri))

		if err != nil {
			return Location{}, err
		}

		location.Url = url.Path

		return location, nil
	})

	if err != nil {
		return nil, err
	}

	return &locations, nil
}

func (n *Nvim) GetTypeCompletion(name string) ([]string, error) {
	buffer, err := n.NewBuffer()

	if err != nil {
		return []string{}, err
	}

	defer buffer.Close()

	err = n.startLSP()

	if err != nil {
		return []string{}, err
	}

	script, err := scripts.Read("get-lsp-completion")

	if err != nil {
		return nil, err
	}

	annotation := fmt.Sprintf("---@type %s", name)
	ref := "local ref"
	trigger := "ref."
	err = buffer.SetLines([]string{annotation, ref, trigger})

	if err != nil {
		return nil, err
	}

	line, character := uint(2), uint(len(trigger))
	result, err := n.execLua(script, []any{line, character, 15000})

	if err != nil {
		return []string{}, fmt.Errorf("Error getting lsp definition: %v", err)
	}

	stringResult, ok := result.(string)

	if !ok {
		return []string{}, fmt.Errorf("Error reading lsp definition response: %v", result)
	}

	response := languageserver.TextDocumentCompletionResponse{}
	err = response.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return []string{}, fmt.Errorf("Error unmarshalling lsp definition response: %w", err)
	}

	if response.Result.IsIncomplete {
		return []string{}, fmt.Errorf("Error reading lsp type completion: the completion response is marked as incomplete")
	}

	completion, _ := slicesx.MapFunc(response.Result.Items, func(item languageserver.CompletionItem) (string, error) {
		if item.InsertText == nil {
			return item.Label, nil
		}

		return *item.InsertText, nil
	})

	return completion, nil
}

func (n *Nvim) GetValueCompletion(value string) ([]string, error) {
	err := n.startLSP()

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", value, err)
	}

	cmd := fmt.Sprintf("lua %s.", value)
	getcompletionResult, err := n.callFunction("getcompletion", []any{cmd, "cmdline"})

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", value, err)
	}

	result, err := anyx.ToSliceOf[string](getcompletionResult)

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", value, err)
	}

	return result, nil
}

func (n *Nvim) GetValueType(value string) (string, error) {
	runtimePath := value
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

	result, err := n.execLua(luaCode, []any{})

	if err != nil {
		return "", fmt.Errorf("Error getting the type of %s: %w", value, err)
	}

	typeName, ok := result.(string)

	if !ok {
		return "", fmt.Errorf("Error getting the type of %s: Error converting the result to a string", runtimePath)
	}

	return typeName, nil
}
