package set

import (
	"iter"
	"maps"
)

type Seter[T comparable] interface {
	Add(T)
}

type Set[T comparable] struct {
	list map[T]struct{} //empty structs occupy 0 memory
}

func (s *Set[T]) Has(v T) bool {
	_, ok := s.list[v]
	return ok
}

func (s *Set[T]) Add(v T) {
	s.list[v] = struct{}{}
}

func (s *Set[T]) Remove(v T) {
	delete(s.list, v)
}

func (s *Set[T]) Clear() {
	s.list = make(map[T]struct{})
}

func (s *Set[T]) Size() int {
	return len(s.list)
}

func (s *Set[T]) Iterate() iter.Seq[T] {
	return maps.Keys(s.list)
}

func NewSet[T comparable]() *Set[T] {
	s := &Set[T]{}
	s.list = make(map[T]struct{})
	return s
}

// AddMulti Add multiple values in the set
func (s *Set[T]) AddMulti(list ...T) {
	for _, v := range list {
		s.Add(v)
	}
}

type FilterFunc[T comparable] func(v T) bool

// Filter returns a subset, that contains only the values that satisfies the given predicate P
func (s *Set[T]) Filter(P FilterFunc[T]) *Set[T] {
	res := &Set[T]{}
	res.list = make(map[T]struct{})
	for v := range s.list {
		if !P(v) {
			continue
		}
		res.Add(v)
	}
	return res
}

func (s *Set[T]) Union(s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	for v := range s.list {
		res.Add(v)
	}

	for v := range s2.list {
		res.Add(v)
	}
	return res
}

func (s *Set[T]) UnionWithTracked(s2 *OrderedSet[T]) *Set[T] {
	res := NewSet[T]()
	for v := range s.list {
		res.Add(v)
	}

	for v := range s2.list {
		res.Add(v)
	}
	return res
}

func (s *Set[T]) Intersect(s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	for v := range s.list {
		if !s2.Has(v) {
			continue
		}
		res.Add(v)
	}
	return res
}

// Difference returns the subset from s, that doesn't exists in s2 (param)
func (s *Set[T]) Difference(s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	for v := range s.list {
		if s2.Has(v) {
			continue
		}
		res.Add(v)
	}
	return res
}

type MapFunc[T comparable, S comparable] func(t T) S

func Map[T comparable, S comparable](s *Set[T], f MapFunc[T, S]) *Set[S] {
	res := NewSet[S]()
	for v := range s.list {
		res.Add(f(v))
	}
	return res
}

func (s *Set[T]) ToList() []T {
	res := []T{}
	for v := range s.list {
		res = append(res, v)
	}
	return res
}
