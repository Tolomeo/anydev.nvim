package anyx

import "fmt"

func ToStringSlice(a any) ([]string, error) {
	sliceOfAny, ok := a.([]any)

	if !ok {
		return []string{}, fmt.Errorf("Error converting %v to array", a)
	}

	sliceOfStrings := make([]string, len(sliceOfAny))

	for index, value := range sliceOfAny {
		str, ok := value.(string)

		if !ok {
			return []string{}, fmt.Errorf("Error converting %v to string", value)
		}

		sliceOfStrings[index] = str
	}

	return sliceOfStrings, nil
}
