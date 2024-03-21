package mqueue

import (
	"fmt"
	"sync"
)

type EmptyWatcher struct {
	isEmpty bool
	cond *sync.Cond
}


func NewEmptyWatcher(empty bool) *EmptyWatcher{
	return &EmptyWatcher{
		isEmpty: empty,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (e *EmptyWatcher) Set(empty bool)  {
	e.cond.L.Lock()
	e.isEmpty = empty
	e.cond.L.Unlock()

	if empty {
		e.cond.Broadcast()
	}
}

func (e *EmptyWatcher) Wait() {
	e.cond.L.Lock()
	defer e.cond.L.Unlock()
	fmt.Println("wait")
	for !e.isEmpty {
		e.cond.Wait()
	}
}
