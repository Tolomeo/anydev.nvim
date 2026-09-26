package set

import "encoding/json"

type Set[T comparable] struct {
	list []T
	bag  map[T]struct{}
}

func (s *Set[T]) Add(item T) {
	if _, ok := s.bag[item]; ok {
		return
	}

	s.list = append(s.list, item)
	s.bag[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) {
	if _, ok := s.bag[item]; !ok {
		return
	}

	delete(s.bag, item)

	for i, v := range s.list {
		if v == item {
			s.list = append(s.list[:i], s.list[i+1:]...)
			return
		}
	}
}

func (s *Set[T]) Clear() {
	s.list = []T{}
	s.bag = make(map[T]struct{})
}

func (s Set[T]) Len() int {
	return len(s.list)
}

func (s Set[T]) Has(item T) bool {
	_, ok := s.bag[item]
	return ok
}

func (s Set[T]) Items() []T {
	result := make([]T, len(s.list))
	copy(result, s.list)
	return result
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Items())
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		list: []T{},
		bag:  make(map[T]struct{}),
	}
}
