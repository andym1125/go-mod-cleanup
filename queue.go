package main

import "container/list"

type Queue struct {
	list.List
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Poll() any {
	return q.Remove(q.Front())
}

func (q *Queue) Push(el any) {
	q.PushBack(el)
}

/* ===== Set ===== */

type Set[T comparable] struct {
	els []T
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{els: make([]T, 0)}
}

//ensures than any duplicates are flattened (ex 3x -> 1x)
func NewSetFromArray[T comparable](arr []T) *Set[T] {
	set := NewSet[T]()
	for _, a := range arr {
		set.Add(a)
	}
	return set
}

func (s *Set[T]) Add(el T) {
	if !s.Contains(el) {
		s.els = append(s.els, el)
	}
}

func (s *Set[T]) Len() int {
	return len(s.els)
}

func (s *Set[T]) Contains(el T) bool {
	ret := false
	for _, a := range s.els {
		if a == el {
			ret = true
		}
	}
	return ret
}

func (s *Set[T]) ToArray() []T {
	ret := make([]T, 0)
	for _, a := range s.els {
		ret = append(ret, a)
	}
	return ret
}

// ========== Nonreceiver funcs ==========

type Collection interface {
	GetAll()
}
