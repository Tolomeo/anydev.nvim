package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AtParamQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
		(param_annotation
			"@param"
			.
			([("...") (identifier)]) @name
			.
			"?"? @optional
			.
			(%s) @type
			.
			(comment)? @documentation
			.
		) @param`, AnyTypeQuery),
}

type AtParam struct {
	name     string
	type_    string
	optional bool
}

func (ap *AtParam) Name() string {
	return ap.name
}

func (ap *AtParam) Type() string {
	return ap.type_
}

func (ap *AtParam) Optional() bool {
	return ap.optional
}

func NewAtParam(captures nvim.TsQueryMatch) *AtParam {
	atParam := AtParam{}

	for _, capture := range captures {
		switch capture.Id {
		case "name":
			atParam.name = capture.Node.Text
		case "optional":
			atParam.optional = true
		case "type":
			atParam.type_ = capture.Node.Text
		}
	}

	return &atParam
}
