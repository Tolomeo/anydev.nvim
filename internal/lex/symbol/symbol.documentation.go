package symbol

import (
	"encoding/json"
	"regexp"
)

var eCommentContent = regexp.MustCompile(`^[ \t]*-{2,3}(.*)$`)

func (d *Documentation) MarshalJSON() ([]byte, error) {
	out := make([]string, len(*d))

	for index, line := range *d {
		matches := eCommentContent.FindStringSubmatch(line)

		if len(matches) > 1 {
			out[index] = matches[1]
			continue
		}

		out[index] = line
	}

	return json.Marshal(out)
}
