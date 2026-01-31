package annotation

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"
)

var FunctionQuery = treesitter.Query{
	Language: "luadoc",
	Query: fmt.Sprintf(`
	(function_type
		(parameter
			(identifier) @parameter.name
			":"
			(%s) @parameter.type
		)? @parameter
		("," (parameter
			(identifier) @parameter.name
			":"
			(%s) @parameter.type
		) @parameter)*
		("," (parameter
			"..." @parameter.name
			":"
			(%s) @parameter.type
		) @parameter)?
		(parameter
			"..." @parameter.name
			":"
			(%s) @parameter.type
		)? @parameter
		(":"
			(%s) @return.type
			("," (%s) @return.type)*
		)? @return
	)
`, AnyTypeQuery, AnyTypeQuery, AnyTypeQuery, AnyTypeQuery, AnyTypeQuery, AnyTypeQuery),
}

type functionArgument struct {
	name  string
	type_ string
}

func (fa functionArgument) Name() string {
	return fa.name
}

func (fa functionArgument) Type() string {
	return fa.type_
}

type functionReturn struct {
	type_ string
}

func (fr functionReturn) Type() string {
	return fr.type_
}

type Function struct {
	arguments []functionArgument
	returns   []functionReturn
}

func (f *Function) Arguments() []functionArgument {
	return f.arguments
}

func (f *Function) Returns() []functionReturn {
	return f.returns
}

func NewFunction(captures nvim.TsQueryMatch) *Function {
	function := Function{
		arguments: []functionArgument{},
		returns:   []functionReturn{},
	}

	for _, capture := range captures {
		switch capture.Id {
		case "parameter":
			function.arguments = append(function.arguments, functionArgument{})
		case "parameter.name":
			function.arguments[len(function.arguments)-1].name = capture.Node.Text
		case "parameter.type":
			function.arguments[len(function.arguments)-1].type_ = capture.Node.Text
		case "return.type":
			function.returns = append(function.returns, functionReturn{type_: capture.Node.Text})
		}
	}

	return &function
}
