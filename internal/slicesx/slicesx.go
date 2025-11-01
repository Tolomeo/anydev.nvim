package slicesx

func FindFunc[T any](s []T, f func(T) bool) (T, bool) {
    var zero T
    for _, v := range s {
        if f(v) {
            return v, true
        }
    }
    return zero, false
}

func AnyToString(s []any) ([]string, bool) {
	sliceOfStrings := make([]string, len(s))

	for index, value := range s {
		str, ok := value.(string)

		if !ok {
			return []string{}, false
		}

		sliceOfStrings[index] = str
	}

	return sliceOfStrings, true
}
