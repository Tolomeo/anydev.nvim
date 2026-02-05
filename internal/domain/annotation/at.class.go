package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var ClassAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(class_annotation
		"@class"
		.
		"(exact)"?
		.
		(identifier) @class.name
		.
		(":" 
			. (%s) @class.parent
			("," (%s) @class.parent)*
		)?
	) @class`, AnyTypeQuery, AnyTypeQuery),
}

func NewClass(captures nvim.TsQueryMatch) {


}
