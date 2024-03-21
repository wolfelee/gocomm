## 作用
 同一时刻对相同的`key`，只执行一次
 
 用于解决redis缓存，缓存击穿（多个请求相同数据，cache miss全部落到db）

## 使用
```go
package main

import (
	"fmt"
	"github.com/wolfelee/gocomm/pkg/util/jsync"
	"time"
)

var sharesCall = jsync.NewSharesCall()

func SelectDB() (models []string,err error) {
	v,err := sharesCall.Do("SelectDB", func() (interface{}, error) {
		// 处理数据逻辑 查询DB Add cache等
		fmt.Println("SelectDB 执行")
		var res = []string{"1","2","3","4"}
		return res,nil
	})
	if err != nil {
		return models,err
	}

	if models,ok := v.([]string); ok {
		return models,nil
	}

	return models,fmt.Errorf("类型断言错误")
}

func main() {
	for i := 0; i < 10; i++ {
		go func() {
			res,err := SelectDB()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res)
		}()
	}
	time.Sleep(time.Second)
}



```

## Do 和 DoEx
 DoEx 返回值比DO多个bool值 判断是否是本函数的执行

