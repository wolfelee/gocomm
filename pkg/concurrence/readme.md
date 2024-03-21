## 并发控制工具

```go
    con := NewWorkerController(2)
	for i := 0; i < 100; i++ {
		con.Go(Helper(i)) //如何配合for 建议使用高阶函数绑定参数
	}
	con.Wait()


func Helper(i int) func() {
        return func() {
                log.Println(i)
                time.Sleep(time.Second)
				
        }
}
```