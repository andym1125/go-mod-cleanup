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
	els []int
}

func NewSet() *Set {
	return &Set{els: make([]int, 0)}
}

func (s *Set) Add(el int) {
	if !s.Contains(el) {
		s.els = append(s.els, el)
	}
}

func (s *Set) Contains(el int) bool {
	ret := false
	for _, a := range s.els {
		if a == el {
			ret = true
		}
	}
	return ret
}

func (s *Set) Get() []int {
	ret := make([]int, 0)
	for _, a := range s.els {
		ret = append(ret, a)
	}
	return ret
}
