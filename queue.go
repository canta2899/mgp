package main

/* Implements a thread safe queue using channels */

import "errors"

// Interface is used in order to use different types
type Queue struct {
	maxc  int
	queue chan interface{}
}

func (q *Queue) Enqueue(entry interface{}) error {
	if len(q.queue) >= q.maxc {
		return errors.New("Queue is full")
	}

	q.queue <- entry
	return nil
}

func (q *Queue) Dequeue() (interface{}, error) {
	if len(q.queue) == 0 {
		return nil, errors.New("Queue is empty")
	}

	next := <-q.queue
	return next, nil
}

func NewQueue(maxc int) *Queue {
	return &Queue{maxc, make(chan interface{}, maxc)}
}
