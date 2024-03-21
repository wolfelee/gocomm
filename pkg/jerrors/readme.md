grpc错误和errno统一

1.http直接返回errno

2.grpc server
  1. grpc可以返回errno.ToGRPC,
  2. 也可以使用中间件通过判断返回值是否是*errno类型,如果是自动调用ToGRPC

2.grpc client

  判断返回值err是否为nil, 如果不是nil,调用FromGRPC来获取errno错误,通过err.Equal(e)来判断是否是同一个错误

老代码迁移(errno.Errno->jerrors.Errno)
在老的errno包里添加如下代码
```go
type Errno = jerrors.Errno  

var NewError = jerrors.NewError
```
删除老errno相关代码即可快速迁移
