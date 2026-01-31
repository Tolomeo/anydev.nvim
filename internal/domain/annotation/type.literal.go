package annotation

import (
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var LiteralStringQuery = treesitter.Query{
	Language: "luadoc",
	Query:    `(literal_type) @stringliteral`,
}

type LiteralString struct {
	value string
}

func (ls *LiteralString) Value() string {
	return ls.value
}

func NewLiteralString(captures nvim.TsQueryMatch) *LiteralString {
	literalString := LiteralString{}

	for _, capture := range captures {
		switch capture.Id {
		case "stringliteral":
			literalString.value = capture.Node.Text
		}
	}

	return &literalString
}

var LiteralBooleanQuery = treesitter.Query{
	Language: "luadoc",
	Query: `(
		(identifier) @literal.value
		(#any-of? @literal.value "true" "false")
	) @literal`,
}

type LiteralBoolean struct {
	value string
}

func (lb *LiteralBoolean) Value() string {
	return lb.value
}

func NewLiteralBoolean(captures nvim.TsQueryMatch) *LiteralBoolean {
	literalBoolean := LiteralBoolean{}

	for _, capture := range captures {
		switch capture.Id {
		case "literal.value":
			literalBoolean.value = capture.Node.Text
		}
	}

	return &literalBoolean
}

var LiteralNumberQuery = treesitter.Query{
	Language: "luadoc",
	Query:    `(numeric_literal_type) @literal.number`,
}

type LiteralNumber struct {
	value string
}

func (ln *LiteralNumber) Value() string {
	return ln.value
}

func NewLiteralNumber(captures nvim.TsQueryMatch) *LiteralNumber {
	literalNumber := LiteralNumber{}

	for _, matchCapture := range captures {
		switch matchCapture.Id {
		case "literal.number":
			literalNumber.value = matchCapture.Node.Text
		}
	}

	return &literalNumber
}
