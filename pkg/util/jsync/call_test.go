package jsync

import (
	"testing"
	"time"
)

func TestShareGroup_Do(t *testing.T) {
	sharescall := NewSharesCall()

	for i := 0; i < 50; i++ {
		go func(i int) {
			v, err := sharescall.Do("shares_call_test", func() (interface{}, error) {
				time.Sleep(time.Second)
				t.Log("do func execute!")
				return i, nil
			})
			if err != nil {
				t.Error("err != nil,err = ", err)
				return
			}
			t.Log("value:-->", v)
		}(i)
	}
	time.Sleep(time.Second * 2)
}
