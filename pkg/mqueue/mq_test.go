package mqueue

import (
	"fmt"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"testing"
	"time"
)

func TestMq(t *testing.T) {
	jlog.DefaultConfig().Build()
	err := Init("./testdata/mq.yml")
	if err != nil {
		t.Error(err)
	}
	t.Log("init success")
	p, e := Use("default").GetProducer("testmq")
	if e != nil {
		t.Error(e)
	}

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(2000 * time.Millisecond)
			p.Produce(PublicData{
				Message:  []byte(fmt.Sprintf("msg:%d", 1)),
				RouteKey: "",
				Delay:    0,
			})
		}

	}()
	//p.WaitEmpty()
	select {}

	//time.Sleep(1 * time.Second)
	//p.WaitEmpty()
	//Use("default").SetConsumerHandler("regsendcoupon", func(ctx *Context) {
	//	t.Log(string(ctx.GetBody()))
	//	t.Log(ctx.GetMsgId())
	//	time.Sleep(100 * time.Millisecond)
	//	if e := ctx.Ack(); e != nil {
	//		t.Error(e)
	//	}
	//})
	//
	//select {
	//
	//}
}

func TestClient(t *testing.T) {
	jlog.DefaultConfig().Build()
	err := Init("./testdata/mq.yml")
	if err != nil {
		t.Error(err)
	}

	Use("default").SetConsumerHandler("testmq", func(ctx *Context) {
		time.Sleep(5 * time.Second)
		t.Log(string(ctx.GetBody()))
		if err := ctx.Ack(); err != nil {
			t.Error(err)
		}

	})
	select {}
}

func TestEmptyWatcher(t *testing.T) {
	w := NewEmptyWatcher(false)
	w.Set(true)
	w.Wait()
	t.Log("case 1 finished")

	w.Set(false)
	go func() {
		time.Sleep(5 * time.Second)
		w.Set(true)
		t.Log("set true")
	}()
	for i := 0; i < 4; i++ {
		go func() {
			fmt.Println("wait before")
			w.Wait()
			fmt.Println("wait after")
		}()
	}
	w.Wait()
	t.Log("success")
}
