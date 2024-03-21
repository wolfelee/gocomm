package xgo

import (
	"fmt"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"runtime"
	"sync"
	"time"

	"github.com/codegangsta/inject"
)

// Serial 串行
func Serial(fns ...func()) func() {
	return func() {
		for _, fn := range fns {
			fn()
		}
	}
}

// Parallel 并发执行
func Parallel(fns ...func()) func() {
	var wg sync.WaitGroup
	return func() {
		wg.Add(len(fns))
		for _, fn := range fns {
			go try2(fn, wg.Done)
		}
		wg.Wait()
	}
}

// RestrictParallel 并发,最大并发量restrict
func RestrictParallel(restrict int, fns ...func()) func() {
	var channel = make(chan struct{}, restrict)
	return func() {
		var wg sync.WaitGroup
		for _, fn := range fns {
			wg.Add(1)
			go func(fn func()) {
				defer wg.Done()
				channel <- struct{}{}
				try2(fn, nil)
				<-channel
			}(fn)
		}
		wg.Wait()
		close(channel)
	}
}

// GoDirect ...
func GoDirect(fn interface{}, args ...interface{}) {
	var inj = inject.New()
	for _, arg := range args {
		inj.Map(arg)
	}

	_, file, line, _ := runtime.Caller(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				jlog.Error("recover", jlog.Any("err", err), jlog.String("line", fmt.Sprintf("%s:%d", file, line)))
			}
		}()
		// 忽略返回值, goroutine执行的返回值通常都会忽略掉
		_, err := inj.Invoke(fn)
		if err != nil {
			jlog.Error(err.Error())
			return
		}
	}()
}

// Go goroutine
func Go(fn func()) {
	go try2(fn, nil)
}

// DelayGo goroutine
func DelayGo(delay time.Duration, fn func()) {
	_, file, line, _ := runtime.Caller(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				jlog.Error("recover", jlog.Any("err", err), jlog.String("line", fmt.Sprintf("%s:%d", file, line)))
			}
		}()
		time.Sleep(delay)
		fn()
	}()
}

// SafeGo safe go
func SafeGo(fn func(), rec func(error)) {
	go func() {
		err := try2(fn, nil)
		if err != nil {
			rec(err)
		}
	}()
}
