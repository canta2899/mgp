package main

/* Thread safe flag for synchronization purposes */

import "sync"

type Flag struct {
    done bool
    mutex sync.Mutex
}

func (f *Flag) Get() bool {
    defer f.mutex.Unlock()
    f.mutex.Lock()
    return f.done 
}

func (f *Flag) Set(val bool) {
    defer f.mutex.Unlock() 
    f.mutex.Lock()
    f.done = val
}

func NewFlag(initval bool) *Flag {
    return &Flag{done: initval}
}
