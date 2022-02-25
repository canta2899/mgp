package main

import (
	"errors"
	"sync"
)

type node struct {
    prev *node
    next *node
    value string
}

type Queue struct {
    mutex sync.Mutex
    cond *sync.Cond
    head *node
    tail *node
    closed bool
}

func (q *Queue) Enqueue(value string) {
    q.mutex.Lock()
    defer q.mutex.Unlock()
    
    if q.head == nil {
        n := &node{value: value}
        q.head = n
        q.tail = n
    } else {
        n := &node{prev: q.tail, value: value}
        q.tail.next = n
        q.tail = n
    }

    q.cond.Signal()
}

func (q *Queue) Dequeue() (string, error) {
    defer q.mutex.Unlock()

    q.mutex.Lock()

    for q.head == nil && !q.closed {
        q.cond.Wait()
    }

    if q.head == nil {
        return "", errors.New("Queue was closed")
    }
 
    n := q.head.value

    if (q.head == q.tail)  {
        q.tail = nil
        q.head = nil
    } else {
        q.head = q.head.next
        q.head.prev = nil
    }

    return n, nil
}

func (q *Queue) Done() {
    q.closed = true 
    q.cond.Broadcast()
}


func NewQueue() *Queue {
    q := &Queue{closed: false}
    q.cond = sync.NewCond(&q.mutex)
    return q
}
