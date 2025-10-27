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
