package crawl

import (
	"fmt"
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/domain/symbol"
	"github.com/Tolomeo/anydev.nvim/internal/domain/target"
	"github.com/Tolomeo/anydev.nvim/internal/log"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/nvim/languageserver"
)

type ctx interface {
	Target() *target.Target
	TargetDefinitionOverride() *string
	Nvim() *nvim.Nvim
	Logger() *log.Logger
}

type Crawler struct {
	context ctx
}

func (c *Crawler) GetOriginChain() (origin.OriginChain, error) {
	originOverride, err := c.getDefinitionOverrideOriginChain()

	if err != nil {
		return nil, err
	}

	if originOverride != nil {
		return originOverride, nil
	}

	var locations *[]nvim.Location

	switch c.context.Target().Kind() {
	case target.TargetKindValue:
		locations, err = c.findIdentifierDefinitionLocations(c.context.Target().Identifier())
	case target.TargetKindType:
		locations, err = c.findTypeIdentifierDefinitionLocations(c.context.Target().Name(), c.context.Target().ParentName())
	case target.TargetKindModule:
		locations, err = c.findModuleDefinitionLocations(c.context.Target().Identifier())
	}

	if err != nil {
		return nil, err
	}

	if locations == nil {
		c.context.Logger().Warnf("No locations found for '%s' %s symbol", c.context.Target().Identifier(), c.context.Target().Kind())
		return origin.NewOriginChain(origin.NewUnkownOrigin()), nil
	}

	origins, err := c.getOriginChain(*locations)

	if err != nil {
		return nil, err
	}

	return origins, nil
}

func (c *Crawler) FollowOriginChain() (origin.OriginChain, error) {
	targetOriginChain := c.context.Target().OriginChain()
	targetOrigin := c.context.Target().Origin()

	switch locationOriginType := targetOrigin.(type) {
	case *origin.VariableOrigin:
		c.context.Logger().Verbosef("Following variable origin <%+v>", locationOriginType)

		url, line, character :=
			locationOriginType.Url(),
			uint(locationOriginType.NameRange().End.Line),
			uint(locationOriginType.NameRange().End.Character)
		rightValueLocations, err := c.findDefinitionLocationsAt(url, line, character)

		if err != nil {
			return nil, err
		}

		variableOriginChain, err := c.getOriginChain(*rightValueLocations)

		if err != nil {
			return nil, err
		}

		return targetOriginChain.Concat(variableOriginChain), nil
	case *origin.FunctionCallOrigin:
		c.context.Logger().Verbosef("Following variable origin <%+v>", locationOriginType)

		url, line, character :=
			locationOriginType.Url(),
			uint(locationOriginType.FunctionNameRange().End.Line),
			uint(locationOriginType.FunctionNameRange().End.Character)
		functionLocations, err := c.findDefinitionLocationsAt(url, line, character)

		if err != nil {
			return nil, err
		}

		functionOriginChain, err := c.getOriginChain(*functionLocations)

		if err != nil {
			return nil, err
		}

		return targetOriginChain.Concat(functionOriginChain), nil
	}

	return nil, fmt.Errorf("Cannot follow origin of type <%T>", targetOrigin)
}

func (c *Crawler) GetDocumentation() (symbol.Documentation, error) {

	var markupContent *languageserver.MarkupContent
	var err error

	switch c.context.Target().Kind() {
	case target.TargetKindValue:
		markupContent, err = c.getDefinitionDocumentation(c.context.Target().Identifier())
	case target.TargetKindType:
		markupContent, err = c.getTypeDefinitionDocumentation(c.context.Target().Name(), c.context.Target().ParentName())
	}

	if err != nil {
		return symbol.Documentation{}, err
	}

	if markupContent == nil {
		return symbol.Documentation{}, nil
	}

	return strings.Split(markupContent.Value, "\n"), nil
}

func NewCrawler(context ctx) *Crawler {
	return &Crawler{
		context: context,
	}
}
