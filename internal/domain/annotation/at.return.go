package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var AtReturnQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
		(return_annotation
			"@return"
			.
			(%s) @return.type
			.
			(comment)? @return.documentation
			.
		) @return`, AnyTypeQuery),
}

type AtReturn struct {
	name  string
	type_ string
}

func (ar *AtReturn) Name() string {
	return ar.name
}

func (ar *AtReturn) Type() string {
	return ar.type_
}

func NewAtReturn(captures nvim.TsQueryMatch) *AtReturn {
	atReturn := AtReturn{}

	for _, capture := range captures {
		switch capture.Id {
		case "return.name":
			atReturn.name = capture.Node.Text
		case "return.type":
			atReturn.type_ = capture.Node.Text
		}
	}

	return &atReturn
}
