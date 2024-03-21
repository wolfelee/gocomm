package localstorage

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewLocalStorage(t *testing.T) {
	nums := genNum(20)
	store := NewLocalStorage(10, false)
	for i, n := range nums {
		exp := time.Duration(n) * time.Second
		key := fmt.Sprintf("aa%d", i)
		t.Log(key, ":", n)
		store.Set(key, n, exp)
	}
	store.Set("aa0", 100, 10*time.Second)
	if !checkHeap(t, store.(*LocalStorage).timeHeap, 0) {
		t.Error("最小堆错误")
	}
	time.Sleep(2 * time.Second)
	t.Log("aa1:", store.Get("aa1"))
	t.Log("aa0:", store.Get("aa0"))

	for i := range nums {
		key := fmt.Sprintf("aa%d", i)
		t.Log(key, store.Get(key))
	}
	time.Sleep(20 * time.Second)
}

func genNum(num int) []int {
	ret := make([]int, 0, num)
	for i := 0; i < num; i++ {
		r := rand.Int() % 100
		if i%2 == 0 {
			r *= -1
		}
		ret = append(ret, r)
	}
	return ret
}
