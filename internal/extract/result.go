package extract

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/extract/crawl"
)

type parser struct {
	crawler *crawl.Crawler
}

func (p *parser) Parse(source crawl.Source) {
	switch v := source.(type) {
	case *crawl.TableSource:
		fmt.Println(v, "table")
	case *crawl.FunctionSource:
		fmt.Println("function")
	case *crawl.VariableSource:
		fmt.Println("variable")
	}
	// fmt.Printf("%+v", source)
}

func NewParser(crawler *crawl.Crawler) *parser {
	prsr := parser{
		crawler: crawler,
	}

	return &prsr
}
