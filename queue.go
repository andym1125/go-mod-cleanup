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
