package crawl

import (
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

type ctx interface {
	Target() *target.Target
	TargetDefinitionOverride() *[]string
	Nvim() *nvim.Nvim
	Logger() *log.Logger
}

type Crawler struct {
	context ctx
}

func NewCrawler(context ctx) *Crawler {
	return &Crawler{
		context: context,
	}
}
