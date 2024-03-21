package concurrence

import (
	"log"
	"testing"
	"time"
)

func TestNewWorkerController(t *testing.T) {
	con := NewWorkerController(2)
	for i := 0; i < 100; i++ {
		con.Go(Helper(i)) //建议使用高阶函数配合for使用
	}
	con.Wait()
}

func Helper(i int) func() {
	return func() {
		log.Println(i)
		time.Sleep(time.Second)
	}
}