package main

/* Thread safe flag for synchronization purposes */

import "sync"

type SyncFlag struct {
    done bool
    mutex sync.Mutex
}

func (f *SyncFlag) Get() bool {
    defer f.mutex.Unlock()
    f.mutex.Lock()
    return f.done 
}

func (f *SyncFlag) Set(val bool) {
    defer f.mutex.Unlock() 
    f.mutex.Lock()
    f.done = val
}

func NewFlag(initval bool) *SyncFlag {
    return &SyncFlag{done: initval}
}
