package lexed

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var eCommentContent = regexp.MustCompile(`^[ \t]*-{2,3}(.*)$`)

func (d *Documentation) MarshalJSON() ([]byte, error) {
	fmt.Println(d)

	out := make([]string, len(*d))

	for index, line := range *d {
		matches := eCommentContent.FindStringSubmatch(line)

		if len(matches) > 1 {
			out[index] = matches[1]
			continue
		}

		out[index] = line
	}

	fmt.Println(out)
	return json.Marshal(out)
}
