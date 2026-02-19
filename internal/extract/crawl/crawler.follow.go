package crawl

import (
	"fmt"

	"github.com/Tolomeo/anydev.nvim/internal/domain/origin"
	"github.com/Tolomeo/anydev.nvim/internal/nvim"
	"github.com/Tolomeo/anydev.nvim/internal/utils/slicesx"
)

func (c *Crawler) FollowOriginChain() (origin.OriginChain, error) {
	switch targetOrigin := c.context.Target().Origin().(type) {
	case *origin.VariableOrigin:
		c.context.Logger().Verbosef("Following variable origin <%+v>", targetOrigin)
		return c.followVariableOriginChain(targetOrigin)

	case *origin.FunctionCallOrigin:
		c.context.Logger().Verbosef("Following variable origin <%+v>", targetOrigin)
		return c.followFunctionCallOrigin(targetOrigin)
	}

	return nil, fmt.Errorf("Cannot follow origin of type <%T>", c.context.Target().Origin())
}

func (c *Crawler) followVariableOriginChain(variableOrigin *origin.VariableOrigin) (origin.OriginChain, error) {
	targetOriginChain := c.context.Target().OriginChain()
	url, line, character :=
		variableOrigin.Url(),
		uint(variableOrigin.NameRange().End.Line),
		uint(variableOrigin.NameRange().End.Character)
	rightValueLocations, err := c.findDefinitionLocationsAt(url, line, character)

	if err != nil {
		return nil, err
	}

	unvisitedLocations, _ := slicesx.FilterFunc(*rightValueLocations, func(rhsLocation nvim.Location) (bool, error) {
		return !targetOriginChain.ContainsLocation(rhsLocation), nil
	})

	if len(unvisitedLocations) == 0 {
		c.context.Logger().Warn("No unvisited locations found")
		return targetOriginChain.Append(origin.NewUnkownOrigin()), nil
	}

	followedOriginChain, err := c.getOriginChain(unvisitedLocations)

	if err != nil {
		return nil, err
	}

	return targetOriginChain.Concat(followedOriginChain), nil
}

func (c *Crawler) followFunctionCallOrigin(functionCallOrigin *origin.FunctionCallOrigin) (origin.OriginChain, error) {
	targetOriginChain := c.context.Target().OriginChain()
	url, line, character :=
		functionCallOrigin.Url(),
		uint(functionCallOrigin.FunctionNameRange().End.Line),
		uint(functionCallOrigin.FunctionNameRange().End.Character)
	functionLocations, err := c.findDefinitionLocationsAt(url, line, character)

	if err != nil {
		return nil, err
	}

	unvisitedLocations, _ := slicesx.FilterFunc(*functionLocations, func(rhsLocation nvim.Location) (bool, error) {
		return !targetOriginChain.ContainsLocation(rhsLocation), nil
	})

	if len(unvisitedLocations) == 0 {
		c.context.Logger().Warn("No unvisited locations found")
		return targetOriginChain.Append(origin.NewUnkownOrigin()), nil
	}

	functionOriginChain, err := c.getOriginChain(*functionLocations)

	if err != nil {
		return nil, err
	}

	return targetOriginChain.Concat(functionOriginChain), nil
}
