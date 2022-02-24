package main

import (
    "sync"
    "errors"
)

type node struct {
    prev *node
    next *node
    value string
}

type Queue struct {
    mutex sync.Mutex
    head *node
    tail *node
}

func (q *Queue) Enqueue(value string) {
    defer q.mutex.Unlock()

    q.mutex.Lock()
    
    if q.head == nil {
        n := &node{value: value}
        q.head = n
        q.tail = n
    } else {
        n := &node{prev: q.tail, value: value}
        q.tail.next = n
        q.tail = n
    }

}

func (q *Queue) Dequeue() (string, error) {
    defer q.mutex.Unlock()

    q.mutex.Lock()

    if q.head == nil {
        return "", errors.New("Empty list")
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


func NewQueue() *Queue {
    return &Queue{}
}
