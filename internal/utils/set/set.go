package set

import "encoding/json"

type StringSet map[string]struct{}

func (s StringSet) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(s))

	for key := range s {
		keys = append(keys, key)
	}

	return json.Marshal(keys)
}
