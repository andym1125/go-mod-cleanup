package main

import "container/list"

/* ========== Queue ========== */

/* Queue is a loose implementation of a FIFO ordered list. */
type Queue struct {
	list.List
}

/* NewQueue creates a new Queue. */
func NewQueue() *Queue {
	return &Queue{}
}

/* Pop removes and returns the first element in this Queue. */
func (q *Queue) Pop() any {
	return q.Remove(q.Front())
}

/* Push pushes an element to the back of the Queue. */
func (q *Queue) Push(el any) {
	q.PushBack(el)
}

/* ===== Set ===== */

/* Set is an implementation of a non-ordered data structure with non-repeating entries. It is a
valid Collection. */
type Set[T comparable] struct {
	els []T
}

/* NewSet creates a new Set */
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{els: make([]T, 0)}
}

/* NewSetFromArray creates a new Set from the given array, preventing any duplicate entries from
creating duplicate elements in the Set */
func NewSetFromArray[T comparable](arr []T) *Set[T] {
	set := NewSet[T]()
	for _, a := range arr {
		set.Add(a)
	}
	return set
}

/* Add adds an element to the Set. This does no operation if the Set already contains this element,
according to Set.Contains(). */
func (s *Set[T]) Add(el T) {
	if !s.Contains(el) {
		s.els = append(s.els, el)
	}
}

/* Len returns the number of elements in this Set. */
func (s *Set[T]) Len() int {
	return len(s.els)
}

/* Contains returns true if this Set already contains the given element, according to the
comparable interface. */
func (s *Set[T]) Contains(el T) bool {
	ret := false
	for _, a := range s.els {
		if a == el {
			ret = true
		}
	}
	return ret
}

/* ToArray returns the elements of this Set, represented as a slice. Multiple calls to add an
element to this Set results in only one copy of that element in the returned slice. This function
ensures compliance with the Collection interface. */
func (s *Set[T]) ToArray() []T {
	ret := make([]T, 0)
	for _, a := range s.els {
		ret = append(ret, a)
	}
	return ret
}

/* ========== General ========== */

/* Collection is an interface to detail iterable data structures */
type Collection[T comparable] interface {
	ToArray() []T
}
