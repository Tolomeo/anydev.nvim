package mapx

func Values[K string, V any](mapValue map[K]V) []V {
	values := []V{}

	for _, value := range mapValue {
		values = append(values, value)
	}

	return values
}

func Keys[K string, V any](mapValue map[K]V) []K {
	keys := []K{}

	for key := range mapValue {
		keys = append(keys, key)
	}

	return keys
}
