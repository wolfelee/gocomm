package concurrence

import (
	"log"
	"sync"
)

//WorkerController 并发控制器
type WorkerController struct {
	max     int //最大数量
	current int //当前数量
	cond *sync.Cond
}

//add 通知制器增加添加一个正在工作的worker, 超过最大数量会等待
func (w *WorkerController) add() {
	w.cond.L.Lock()
	w.current += 1
	for w.current > w.max {
		w.cond.Wait()
	}
	w.cond.L.Unlock()
}

//done 通知控制一个worker已完成
func (w *WorkerController) done() {
	w.cond.L.Lock()
	w.current -= 1
	w.cond.L.Unlock()
	w.cond.Broadcast()
}

//Wait 等待所有worker完成
func (w *WorkerController) Wait() {
	w.cond.L.Lock()
	for w.current > 0 {
		w.cond.Wait()
	}
	w.cond.L.Unlock()
}

//Go 开启一个worker
func (w *WorkerController) Go(f func()) {
	w.add()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic", err)
			}
			w.done()
		}()
		f()
	}()
}

//NewWorkerController 创建一个控制器
func NewWorkerController(max int) *WorkerController {
	return &WorkerController{
		max:  max,
		cond: sync.NewCond(&sync.Mutex{}),
	}
}
