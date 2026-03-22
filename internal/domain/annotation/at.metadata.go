package annotation

import "github.com/Tolomeo/anydev.nvim/internal/nvim/treesitter"

var AtPrivateAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(qualifier_annotation "@private") @private`,
}

var AtProtectedAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(qualifier_annotation "@protected") @protected`,
}

var AtPackageAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(package_annotation) @package`,
}

var AtDeprectedAnnotationQuery = treesitter.Query{
	Language: "luadoc",
	Query: `
	(deprecated_annotation) @deprecated`,
}
