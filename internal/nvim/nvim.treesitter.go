package nvim

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/Tolomeo/anydev.nvim/internal/nvim/internal/scripts"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
	"github.com/Tolomeo/anydev.nvim/internal/utils/mapx"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (n *Nvim) startTS() error {
	script, err := scripts.Read("start-ts")

	if err != nil {
		return err
	}

	_, err = n.execLua(script, []any{30000})

	if err != nil {
		return fmt.Errorf("Error starting treesitter lua: %w", err)
	}

	return nil
}

func (n *Nvim) execTsQuery(query treesitter.Query) (*[]treesitter.Capture, error) {
	err := n.startTS()

	if err != nil {
		return nil, err
	}

	script, err := scripts.Read("exec-ts-query")

	if err != nil {
		return nil, err
	}

	scriptArgs := []any{query.Language, query.Query}

	if query.Range != nil {
		scriptArgs = append(scriptArgs, query.Range.Start, query.Range.End+1)
	}

	result, err := n.execLua(script, scriptArgs)

	// fmt.Printf("\nQuery: \n%v\n%v\n%v\n", config.Query, result, err)

	switch {
	case err != nil:
		return nil, fmt.Errorf("Nvim TSQuery error: %w", err)
	case result == nil:
		return nil, nil
	}

	stringResult, ok := result.(string)

	if !ok {
		return nil, fmt.Errorf("Error reading tsNodes query result as a string: %v", result)
	}

	var captures []treesitter.Capture

	if err := json.Unmarshal([]byte(stringResult), &captures); err != nil {
		return nil, fmt.Errorf("Error decoding tsNodes json response: %w", err)
	}

	return &captures, nil
}

type TsQueryMatch []treesitter.Capture

func (m *TsQueryMatch) LineRange() *treesitter.LineRange {
	startLines, _ := slicesx.MapFunc(*m, func(capture treesitter.Capture) (float64, error) {
		return capture.Node.Range.Start.Line, nil
	})
	endLines, _ := slicesx.MapFunc(*m, func(capture treesitter.Capture) (float64, error) {
		return capture.Node.Range.End.Line, nil
	})

	return &treesitter.LineRange{
		Start: slices.Min(startLines),
		End:   slices.Max(endLines),
	}
}

// TODO: create a custom slicesx.MaxFunc util
func (m *TsQueryMatch) Range() *treesitter.Range {
	lineRange := m.LineRange()

	startRangeCaptures, _ := slicesx.FilterFunc(*m, func(capture treesitter.Capture) (bool, error) {
		return capture.Node.Range.Start.Line == lineRange.Start, nil
	})

	startCharacters, _ := slicesx.MapFunc(startRangeCaptures, func(capture treesitter.Capture) (float64, error) {
		return capture.Node.Range.Start.Character, nil
	})

	endRangeCaptures, _ := slicesx.FilterFunc(*m, func(capture treesitter.Capture) (bool, error) {
		return capture.Node.Range.End.Line == lineRange.End, nil
	})

	endCharacters, _ := slicesx.MapFunc(endRangeCaptures, func(capture treesitter.Capture) (float64, error) {
		return capture.Node.Range.End.Character, nil
	})

	return &treesitter.Range{
		Start: treesitter.Position{
			Line:      lineRange.Start,
			Character: slices.Min(startCharacters),
		},
		End: treesitter.Position{
			Line:      lineRange.End,
			Character: slices.Max(endCharacters),
		},
	}

}

func (n *Nvim) tsQueryAll(config treesitter.Query) (*[]TsQueryMatch, error) {
	queryAllConfig := treesitter.Query{
		Language: config.Language,
		Query:    fmt.Sprintf("(%s) @tsquery.match", config.Query),
		Range:    config.Range,
	}

	captures, err := n.execTsQuery(queryAllConfig)

	switch {
	case err != nil:
		return nil, err
	case captures == nil:
		return nil, nil
	}

	/* fmt.Println("query all captures")
	fmt.Printf("\n\n%+v\n\n", captures) */

	queryCaptures, _ := slicesx.FilterFunc(*captures, func(capture treesitter.Capture) (bool, error) {
		return (capture.Id == "tsquery.match"), nil
	})

	queryMatches, _ := slicesx.MapFunc(queryCaptures, func(queryCapture treesitter.Capture) (TsQueryMatch, error) {
		return slicesx.FilterFunc(*captures, func(capture treesitter.Capture) (bool, error) {
			if capture.Id == queryCapture.Id {
				return false, nil
			}

			return queryCapture.Node.Contains(capture.Node), nil
		})
	})

	/* fmt.Println("query all matches")
	fmt.Printf("\n\n%+v\n\n", queryMatches) */

	return &queryMatches, nil
}

func (n *Nvim) tsQueryOne(query treesitter.Query) (*TsQueryMatch, error) {
	matches, err := n.tsQueryAll(query)

	switch {
	case err != nil:
		return nil, err
	case matches == nil:
		return nil, nil
	case len(*matches) > 1:
		return nil, fmt.Errorf("Error executing TSQuery: too many matches, expected 1 but received %d <%+v>", len(*matches), matches)
	case len(*matches) < 1:
		return nil, nil
	}

	/* fmt.Println("query one")
	fmt.Printf("\n\n%+v\n\n", matches) */

	match := (*matches)[0]

	return &match, nil
}

// var ErrSafeTSQueryNoMatch = errors.New("The parsed language tree contains errors")

type SafeTsQueryResult struct {
	HasError bool
	Captures TsQueryMatch
}

func (n *Nvim) safeTsQueryOne(query treesitter.Query) (*SafeTsQueryResult, error) {
	match, err := n.tsQueryOne(query)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	errorCaptures, err := n.tsQueryAll(treesitter.Query{
		Language: query.Language,
		Query:    `(ERROR) @tsquery.error`,
		Range:    match.LineRange(),
	})

	switch {
	case err != nil:
		return nil, err
	case errorCaptures != nil:
		return &SafeTsQueryResult{
			HasError: true,
			Captures: *match,
		}, nil
	default:
		return &SafeTsQueryResult{
			HasError: false,
			Captures: *match,
		}, nil
	}
}

func (n *Nvim) safeTsQueryAll(query treesitter.Query) (*[]SafeTsQueryResult, error) {
	matches, err := n.tsQueryAll(query)

	switch {
	case err != nil:
		return nil, err
	case matches == nil:
		return nil, nil
	}

	results := []SafeTsQueryResult{}

	for _, match := range *matches {
		errorCaptures, err := n.tsQueryAll(treesitter.Query{
			Language: query.Language,
			Query:    `(ERROR) @tsquery.error`,
			Range:    match.LineRange(),
		})

		// fmt.Printf("\n%+v\n\n", errorCaptures)

		switch {
		case err != nil:
			return nil, err
		case errorCaptures != nil:
			results = append(results, SafeTsQueryResult{
				HasError: true,
				Captures: match,
			})
		default:
			results = append(results, SafeTsQueryResult{
				HasError: false,
				Captures: match,
			})
		}
	}

	return &results, nil
}

func (n *Nvim) getTSNodeAt(nodeTypes []string, line uint, character uint) (*treesitter.TsNode, error) {
	err := n.startTS()

	if err != nil {
		return nil, err
	}

	script, err := scripts.Read("get-ts-node")

	if err != nil {
		return nil, err
	}

	result, err := n.execLua(string(script), []any{nodeTypes, line, character})

	switch {
	case err != nil:
		return nil, err
	case result == nil:
		return nil, nil
	}

	stringResult, ok := result.(string)

	if !ok {
		return nil, fmt.Errorf("Error converting result into string: %v", result)
	}

	tsNode := treesitter.TsNode{}
	err = tsNode.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling tsnode response: %w", err)
	}

	return &tsNode, nil
}

func (n *Nvim) getTsCommentBlockAt(line uint, character uint) (*treesitter.TsNode, error) {
	lines, err := n.getBufferLines(int(line), int(line)+1)

	switch {
	case err != nil:
		return nil, err
	case len(lines) < 1:
		return nil, nil
	}

	// clamping the received character to be inside the line
	character = max(0, min(character, uint(len(lines[0])-1)))

	node, err := n.getTSNodeAt([]string{treesitter.COMMENT}, line, character)

	switch {
	case err != nil:
		return nil, err
	case node == nil:
		return nil, nil
	}

	script, err := scripts.Read("get-ts-commentblock")

	if err != nil {
		return nil, err
	}

	result, err := n.execLua(script, []any{node.Range.Start.Line, node.Range.End.Line})

	switch {
	case err != nil:
		return nil, err
	case result == nil:
		return nil, nil
	}

	stringResult, ok := result.(string)

	if !ok {
		return nil, fmt.Errorf("Error converting result into string: %v", result)
	}

	tsNode := treesitter.TsNode{}
	err = tsNode.UnmarshalJSON([]byte(stringResult))

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling tsnode response: %w", err)
	}

	return &tsNode, nil
}

type TsNodeQueryMatch struct {
	Node  treesitter.TsNode
	Match TsQueryMatch
}

func (n *Nvim) tsNodeQueryAt(queryMap map[string]treesitter.Query, line uint, character uint) (*TsNodeQueryMatch, error) {
	targetNodes := mapx.Keys(queryMap)

	node, err := n.getTSNodeAt(targetNodes, line, character)

	switch {
	case err != nil:
		return nil, err
	case node == nil:
		return nil, nil
	}

	nodeQuery := queryMap[node.Type]
	nodeRange := node.Range.LineRange()
	rangedNodeQuery := treesitter.Query{
		Language: nodeQuery.Language,
		Query:    nodeQuery.Query,
		Range:    &nodeRange,
	}

	match, err := n.tsQueryOne(rangedNodeQuery)

	switch {
	case err != nil:
		return nil, err
	case match == nil:
		return nil, nil
	}

	return &TsNodeQueryMatch{
		Node:  *node,
		Match: *match,
	}, nil
}
