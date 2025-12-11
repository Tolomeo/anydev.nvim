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

func IndexFunc[T any](s []T, f func(T) (bool, error)) (int, error) {
	var notFound = -1

	for i, v := range s {
		found, err := f(v)

		if err != nil {
			return notFound, err
		}

		if found {
			return i, nil
		}
	}

	return notFound, nil
}

func MapFunc[T any, U any](ts []T, f func(T) (U, error)) ([]U, error) {
	us := make([]U, len(ts))

	for i, t := range ts {
		v, err := f(t)

		if err != nil {
			return us, err
		}

		us[i] = v
	}

	return us, nil
}

func MapFindFunc[T any, R any](
	slice []T,
	mapper func(T) (R, bool, error),
) (R, bool, error) {
	var r R

	for _, item := range slice {
		result, matched, err := mapper(item)

		switch {
		case matched && err != nil:
			return r, false, err // Stop immediately and return the error
		case matched:
			return result, true, nil // Success!
		}
	}

	return r, false, nil // Not found, no error
}
