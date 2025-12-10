package mapx

func Values[K string, V any](mapValue map[K]V) []V {
	values := []V{}

	for _, value := range mapValue {
		values = append(values, value)

	}

	return values
}
