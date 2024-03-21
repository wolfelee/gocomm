package redislock

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestNewRedisLock(t *testing.T) {
	rl, err := NewRedisLock("key")
	if err != nil {
		t.Error(err)
		return
	}
	err = rl.Lock()
	if err != nil {
		t.Error(err)
	}
	t.Log("lock success")
	//err = rl.Unlock()
	//if err != nil {
	//	t.Error(err)
	//}
	time.Sleep(21 * time.Second)
	rl.Unlock()
	time.Sleep(10 * time.Second)
}

func TestNewRedisLock2(t *testing.T) {
	var a = 0
	for i := 0; i < 20; i++ {
		t.Run(fmt.Sprintf("g:%d", i), func(t *testing.T) {
			t.Parallel()
			rl, err := NewRedisLock("key1")
			if err != nil {
				t.Error(err)
				return
			}
			//ctx, cannel := context.WithTimeout(context.Background(), time.Second*10)
			//defer cannel()
			err = rl.Lock()
			//err = rl.LockWithCtx(ctx)
			if err != nil {

				t.Error("lock err", err)
				return
			}
			// time.Sleep(time.Second)
			a++
			t.Log("lock a=", a)
			rl.Unlock()
			t.Log("unlock")
		})
	}
}

func TestNewRedisLock3(t *testing.T) {
	wg := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	var a = 0
	num := 100
	wg.Add(1)
	wg2.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			wg.Wait()
			defer wg2.Done()
			rl, err := NewRedisLock("key2")
			if err != nil {
				t.Error(err)
				return
			}
			err = rl.Lock()
			if err != nil {
				t.Error(err)
				return
			}
			defer rl.Unlock()
			a++
		}()
	}
	wg.Done()
	wg2.Wait()
	if a != num {
		t.Error("not", num)
	}
	t.Log(a)
}

func TestNewRedisLockV2(t *testing.T) {
	var a = 0
	for i := 0; i < 20; i++ {
		t.Run(fmt.Sprintf("g:%d", i), func(t *testing.T) {
			t.Parallel()
			rl, err := NewRedisLockV2("key1")
			if err != nil {
				t.Error(err)
				return
			}
			//ctx, cannel := context.WithTimeout(context.Background(), time.Second*10)
			//defer cannel()
			err = rl.Lock()
			//err = rl.LockWithCtx(ctx)
			if err != nil {

				t.Error("lock err", err)
				return
			}
			// time.Sleep(time.Second)
			a++
			t.Log("lock a=", a)
			rl.Unlock()
			t.Log("unlock")
		})
	}
	// time.Sleep(5 * time.Second)
	// t.Log("----------------a=", a)
}

func TestMain(m *testing.M) {
	client := redis.NewClient(&redis.Options{
		Addr:     "mm-dev-redis4-0-14.jd100.com:7430",
		Password: "e5ebb00c9f8914dc4c6db09a30cb8877",
		DB:       0,
	})
	SetRedisClient(client)
	log.Println("test main")
	m.Run()
}
