package crawl

import (
	"strings"

	"github.com/Tolomeo/anydev.nvim/internal/nvim"
)

func (c *Crawler) getAnnotations(location nvim.Location) ([]string, error) {
	buffer, err := c.context.Nvim().OpenBuffer(location.Url)

	if err != nil {
		return nil, err
	}

	defer buffer.Close()

	annotations, err := buffer.GetTsCommentBlockAt(location.StartLine()-1, location.StartCharacter())

	if err != nil {
		return nil, err
	}

	if annotations == nil {
		c.context.Logger().Warnf("No documentation found for location <%v>", location)
		return []string{}, nil
	}

	return strings.Split(annotations.Text, "\n"), nil
}
