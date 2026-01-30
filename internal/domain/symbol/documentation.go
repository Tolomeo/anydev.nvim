package symbol

import "encoding/json"

type Documentation []string

func (d Documentation) MarshalJSON() ([]byte, error) {
	if d == nil {
		return []byte("[]"), nil
	}

	return json.Marshal([]string(d))
}
