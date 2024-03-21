package jsync

import "sync"

type SharesCall interface {
	Do(key string, fn func() (interface{}, error)) (interface{}, error)
	DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error)
}

type call struct {
	wg    sync.WaitGroup
	value interface{}
	err   error
}

type ShareGroup struct {
	lock  sync.Mutex
	calls map[string]*call
}

func NewSharesCall() SharesCall {
	return &ShareGroup{
		calls: make(map[string]*call),
	}
}

func (sg *ShareGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	c, done := sg.createCall(key)
	if done {
		return c.value, c.err
	}

	sg.makeCall(c, key, fn)
	return c.value, c.err
}

func (sg *ShareGroup) DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error) {
	c, done := sg.createCall(key)
	if done {
		return c.value, false, c.err
	}

	sg.makeCall(c, key, fn)
	return c.value, true, c.err
}

func (sg *ShareGroup) createCall(key string) (*call, bool) {
	sg.lock.Lock()
	if call, ok := sg.calls[key]; ok {
		sg.lock.Unlock()
		call.wg.Wait()
		return call, true
	}

	call := new(call)
	call.wg.Add(1)
	sg.calls[key] = call
	sg.lock.Unlock()
	return call, false
}

func (sg *ShareGroup) makeCall(c *call, key string, fn func() (interface{}, error)) {
	defer func() {
		sg.lock.Lock()
		delete(sg.calls, key)
		sg.lock.Unlock()
		c.wg.Done()
	}()

	c.value, c.err = fn()
}
