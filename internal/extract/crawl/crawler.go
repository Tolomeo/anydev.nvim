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
	Nvim() *nvim.Nvim
	Logger() *log.Logger
}

type Crawler struct {
	context ctx
}

func (c *Crawler) GetOriginChain() (origin.OriginChain, error) {
	var locations *[]nvim.Location
	var err error

	switch c.context.Target().Kind() {
	case target.TargetKindValue:
		locations, err = c.findIdentifierDefinitionLocations(c.context.Target().Identifier())
	case target.TargetKindType:
		locations, err = c.findTypeIdentifierDefinitionLocations(c.context.Target().Name(), c.context.Target().ParentName())
	}

	if err != nil {
		return nil, err
	}

	if locations == nil {
		c.context.Logger().Warnf("No locations found for '%s' %s symbol", c.context.Target().Identifier(), c.context.Target().Kind())
		return nil, nil
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
	case *origin.ModuleRequireOrigin:
		c.context.Logger().Verbosef("Following module origin <%+v>", locationOriginType)

		moduleName := locationOriginType.Name()
		moduleLocations, err := c.findModuleRequireDefinitionLocations(moduleName)

		if err != nil {
			return nil, err
		}

		moduleRequireOriginChain, err := c.getOriginChain(*moduleLocations)

		if err != nil {
			return nil, err
		}

		return targetOriginChain.Append(moduleRequireOriginChain), nil
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

		return targetOriginChain.Append(variableOriginChain), nil
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
