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

type Set struct {
	els []any
}

func NewSet() *Set {
	return &Set{els: make([]any, 0)}
}

func (s *Set) Add(el any) {
	if !s.Contains(el) {
		s.els = append(s.els, el)
	}
}

func (s *Set) Len() int {
	return len(s.els)
}

func (s *Set) Contains(el any) bool {
	ret := false
	for _, a := range s.els {
		if a == el {
			ret = true
		}
	}
	return ret
}

func (s *Set) Get() []any {
	ret := make([]any, 0)
	for _, a := range s.els {
		ret = append(ret, a)
	}
	return ret
}
