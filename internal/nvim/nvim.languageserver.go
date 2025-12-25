package nvim

import (
	"fmt"
	"net/url"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/internal/scripts"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/languageserver"
	"github.com/Tolomeo/anydev.nvim/internal/utils/anyx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (n *Nvim) startLSP() error {
	script, err := scripts.Read("start-lsp")

	if err != nil {
		return fmt.Errorf("Error starting lua lsp: %v", err)
	}

	_, err = n.ExecLua(script, []any{30000})

	if err != nil {
		return fmt.Errorf("Error starting lua lsp: %v", err)
	}

	return nil
}

func (n *Nvim) GetDocumentSymbols() (*languageserver.TextDocumentDocumentSymbolResponse, error) {
	documentSymbols := languageserver.TextDocumentDocumentSymbolResponse{}

	err := n.startLSP()

	if err != nil {
		return nil, err
	}

	script, err := scripts.Read("get-lsp-document-symbols")

	if err != nil {
		return nil, fmt.Errorf("Error getting document symbols: %v", err)
	}

	result, err := n.ExecLua(script, []any{})

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

func (n *Nvim) GetHover(line uint, character uint) (*languageserver.TextDocumentHoverResponse, error) {
	hover := languageserver.TextDocumentHoverResponse{}

	err := n.startLSP()

	if err != nil {
		return nil, err
	}

	script, err := scripts.Read("get-lsp-hover")

	if err != nil {
		return nil, err
	}

	result, err := n.ExecLua(script, []any{line, character, 15000})

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

func (n *Nvim) GetLSPDefinitions(line uint, character uint) (*[]languageserver.DefinitionLocation, error) {
	err := n.startLSP()

	if err != nil {
		return nil, err
	}

	script, err := scripts.Read("get-lsp-definition-locations")

	if err != nil {
		return nil, err
	}

	result, err := n.ExecLua(script, []any{line, character, 15000})

	if err != nil {
		return nil, fmt.Errorf("Error getting lsp definition: %v", err)
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

func (n *Nvim) GetDefinitionLocations(line uint, character uint) (*[]Location, error) {
	lspDefinitions, err := n.GetLSPDefinitions(line, character)

	switch {
	case err != nil:
		return nil, err
	case lspDefinitions == nil:
		return nil, nil
	case len(*lspDefinitions) == 0:
		return nil, nil
	}

	locations, err := slicesx.MapFunc(*lspDefinitions, func(lspLocation languageserver.DefinitionLocation) (Location, error) {
		location := Location{
			DefinitionLocation: lspLocation,
		}

		url, err := url.Parse(string(location.DefinitionLocation.TargetUri))

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

func (n *Nvim) GetCompletion(head string) ([]string, error) {
	err := n.startLSP()

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", head, err)
	}

	cmd := "lua " + head + "."
	getcompletionResult, err := n.CallFunction("getcompletion", []any{cmd, "cmdline"})

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", head, err)
	}

	result, err := anyx.ToSliceOf[string](getcompletionResult)

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", head, err)
	}

	return result, nil
}
