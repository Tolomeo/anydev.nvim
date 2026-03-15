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

// https://www.lua.org/manual/5.4/manual.html#2.4
var metamethods = map[string]struct{}{
	"__add":      {},
	"__sub":      {},
	"__mul":      {},
	"__div":      {},
	"__mod":      {},
	"__pow":      {},
	"__unm":      {},
	"__idiv":     {},
	"__band":     {},
	"__bor":      {},
	"__bxor":     {},
	"__bnot":     {},
	"__shl":      {},
	"__shr":      {},
	"__concat":   {},
	"__len":      {},
	"__eq":       {},
	"__lt":       {},
	"__le":       {},
	"__index":    {},
	"__newindex": {},
	"__call":     {},
}

// https://www.lua.org/manual/5.1/manual.html#2.1
var keywords = map[string]struct{}{
	"and":      {},
	"break":    {},
	"do":       {},
	"else":     {},
	"elseif":   {},
	"end":      {},
	"false":    {},
	"for":      {},
	"function": {},
	"if":       {},
	"in":       {},
	"local":    {},
	"nil":      {},
	"not":      {},
	"or":       {},
	"repeat":   {},
	"return":   {},
	"then":     {},
	"true":     {},
	"until":    {},
	"while":    {},
}

type Location struct {
	languageserver.Location
	Url string
}

func NewLocation(url string, startLine, startCharacter, endLine, endCharacter int) Location {
	return Location{
		Location: languageserver.Location{
			TargetRange: languageserver.Range{
				Start: languageserver.Position{
					Line:      float64(startLine),
					Character: float64(startCharacter),
				},
				End: languageserver.Position{
					Line:      float64(endLine),
					Character: float64(endCharacter),
				},
			},
		},
		Url: url,
	}
}

func (l *Location) StartLine() uint {
	return uint(l.TargetRange.Start.Line)
}

func (l *Location) StartCharacter() uint {
	return uint(l.TargetRange.Start.Character)
}

/* func (n *Nvim) getDocumentSymbols() (*languageserver.TextDocumentDocumentSymbolResponse, error) {
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
} */

func (n *Nvim) getLspHover(line uint, character uint) (*languageserver.MarkupContent, error) {
	script, err := scripts.Read("get-lsp-hover")

	if err != nil {
		return nil, err
	}

	hover := languageserver.TextDocumentHoverResponse{}

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

	return &hover.Result.Contents, nil
}

func (n *Nvim) getLSPDefinitions(line uint, character uint) (*[]languageserver.Location, error) {
	script, err := scripts.Read("get-lsp-definition-locations")

	if err != nil {
		return nil, err
	}

	result, err := n.execLua(script, []any{line, character, 15000})

	if err != nil {
		return nil, fmt.Errorf("Error getting lsp definition: %v", err)
	}

	if result == nil {
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
	script, err := scripts.Read("get-lsp-type-definition-locations")

	if err != nil {
		return nil, err
	}

	result, err := n.execLua(script, []any{line, character, 15000})

	if err != nil {
		return nil, fmt.Errorf("Error getting lsp type definition locations: %v", err)
	}

	if result == nil {
		return nil, nil
	}

	stringResult, ok := result.(string)

	if !ok {
		return nil, fmt.Errorf("Error reading lsp type definition response: %v", result)
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

	if err != nil {
		return nil, err
	}

	if lspTypeDefinitions == nil {
		return nil, nil
	}

	if len(*lspTypeDefinitions) == 0 {
		return nil, nil
	}

	// removing any results pointing to scratch buffers
	typeDefinitions, _ := slicesx.FilterFunc(*lspTypeDefinitions, func(typeDefinition languageserver.Location) (bool, error) {
		return !scratchBufferName.Match([]byte(typeDefinition.TargetUri)), nil
	})

	locations, err := slicesx.MapFunc(typeDefinitions, func(lspLocation languageserver.Location) (Location, error) {
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

	if err != nil {
		return nil, err
	}

	if lspDefinitions == nil {
		return nil, nil
	}

	if len(*lspDefinitions) == 0 {
		return nil, nil
	}

	// removing any results pointing to scratch buffers
	definitions, _ := slicesx.FilterFunc(*lspDefinitions, func(definition languageserver.Location) (bool, error) {
		return !scratchBufferName.Match([]byte(definition.TargetUri)), nil
	})

	locations, err := slicesx.MapFunc(definitions, func(definition languageserver.Location) (Location, error) {
		location := Location{
			Location: definition,
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

	if result == nil {
		return []string{}, nil
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

	return slicesx.FilterFunc(completion, func(completionItem string) (bool, error) {
		_, isMetamethod := metamethods[completionItem]
		return !isMetamethod, nil
	})
}

func (n *Nvim) GetValueCompletion(value string) ([]string, error) {
	_, err := n.execLua("_G.Anydev:wait_for_lsp_idle()", []any{})

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", value, err)
	}

	cmd := fmt.Sprintf("lua %s.", value)
	completionResult, err := n.callFunction("getcompletion", []any{cmd, "cmdline"})

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", value, err)
	}

	completion, err := anyx.ToSliceOf[string](completionResult)

	if err != nil {
		return []string{}, fmt.Errorf("Error getting completion for %s: %w", value, err)
	}

	return slicesx.FilterFunc(completion, func(completionItem string) (bool, error) {
		_, isMetamethod := metamethods[completionItem]
		return !isMetamethod, nil
	})
}

func (n *Nvim) GetValueType(value string) (string, error) {
	runtimePath := value
	parts := strings.Split(runtimePath, ".")

	if len(parts) > 1 {
		tail := parts[len(parts)-1]

		if _, isKeyword := keywords[tail]; isKeyword {
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
