package set

import (
	"iter"
	"maps"
	"sort"
)

type OrderedSet[T comparable] struct {
	list  map[T]int //empty structs occupy 0 memory
	count int
}

func (s *OrderedSet[T]) Has(v T) bool {
	_, ok := s.list[v]
	return ok
}

func (s *OrderedSet[T]) Ord(v T) int {
	ord, ok := s.list[v]

	if !ok {
		return -1
	}

	return ord
}

func (s *OrderedSet[T]) Add(v T) {
	s.list[v] = s.count
	s.count += 1
}

func (s *OrderedSet[T]) Remove(v T) {
	delete(s.list, v)
}

func (s *OrderedSet[T]) Clear() {
	s.list = make(map[T]int)
}

func (s *OrderedSet[T]) Size() int {
	return len(s.list)
}

func (s *OrderedSet[T]) Iterate() iter.Seq[T] {
	return maps.Keys(s.list)
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	s := &OrderedSet[T]{}
	s.list = make(map[T]int)
	s.count = 0
	return s
}

// AddMulti Add multiple values in the set
func (s *OrderedSet[T]) AddMulti(list ...T) {
	for _, v := range list {
		s.Add(v)
	}
}

// Filter returns a subset, that contains only the values that satisfies the given predicate P
func (s *OrderedSet[T]) Filter(P FilterFunc[T]) *OrderedSet[T] {
	res := &OrderedSet[T]{}
	res.list = make(map[T]int)
	for v := range s.list {
		if !P(v) {
			continue
		}
		res.Add(v)
	}
	return res
}

func (s *OrderedSet[T]) Union(s2 *OrderedSet[T]) *OrderedSet[T] {
	res := NewOrderedSet[T]()
	for v := range s.list {
		res.Add(v)
	}

	for v := range s2.list {
		res.Add(v)
	}
	return res
}

func (s *OrderedSet[T]) Intersect(s2 *OrderedSet[T]) *OrderedSet[T] {
	res := NewOrderedSet[T]()
	for v := range s.list {
		if !s2.Has(v) {
			continue
		}
		res.Add(v)
	}
	return res
}

// Difference returns the subset from s, that doesn't exists in s2 (param)
func (s *OrderedSet[T]) Difference(s2 *OrderedSet[T]) *OrderedSet[T] {
	res := NewOrderedSet[T]()
	for v := range s.list {
		if s2.Has(v) {
			continue
		}
		res.Add(v)
	}
	return res
}

func MapTracked[T comparable, S comparable](s *OrderedSet[T], f MapFunc[T, S]) *OrderedSet[S] {
	res := NewOrderedSet[S]()
	for v := range s.list {
		res.Add(f(v))
	}
	return res
}

func (s *OrderedSet[T]) ToList() []T {
	res := []T{}
	for v := range s.list {
		res = append(res, v)
	}

	sort.Slice(res, func(i, j int) bool {
		return s.Ord(res[i]) < s.Ord(res[j])
	})

	return res
}
