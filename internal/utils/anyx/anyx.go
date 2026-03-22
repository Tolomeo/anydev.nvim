package anyx

import "fmt"

func ToSliceOf[T any](a any) ([]T, error) {
	sliceOfAny, ok := a.([]any)

	if !ok {
		return []T{}, fmt.Errorf("Error converting %v to slice", a)
	}

	sliceOf := make([]T, len(sliceOfAny))

	for index, value := range sliceOfAny {
		t, ok := value.(T)

		if !ok {
			return []T{}, fmt.Errorf("Error converting %v to string", value)
		}

		sliceOf[index] = t
	}

	return sliceOf, nil
}
