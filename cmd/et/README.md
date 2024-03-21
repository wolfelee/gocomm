## 安装
go >= 1.16
```bash
go install github.com/wolfelee/gocomm/cmd/et@latest
```
go < 1.16
```bash
go get -u github.com/wolfelee/gocomm/cmd/et
```
测试
```bash
et -v
easyTech version v0.0.1
```

## 使用
### 生成新的项目
```bash
et new hello
```

### 升级插件 (et,protoc-gen-go-grpc,cmd/protoc-gen-go)

```bash
et upgrade
```

**注意升级,grpc插件有可能不兼容**

### proto

1. add 

```bash
et proto add internal/grpc/pb/user/user.proto
```

2. build

只编译proto文件

```bash
et proto build internal/grpc/pb/user/user.proto       
```

3. server 

编译proto文件 + 生成基础代码

```bash
et proto server internal/grpc/pb/user/user.proto --dir internal/grpc/svc/user
```

## model
目前仅支持 datasource 

| 命令  | 作用                     | 默认            |
| ----- | ------------------------ | --------------- |
| url   | mysql url 可以放到path中 | 无              |
| table | 表名字 可以用*           | 无              |
| dir   | 生成model的目录          | internal/models |
| cache | 是否使用缓存             | False           |

例子：

```bash
et model datasource --url="root:zhy1996@tcp(localhost:3306)/test" --table="*" --dir="internal/models/user1" --cache=true
```

## 计划
1. 根据SQL生成简单增删改查代码 ✅
2. 生成http api 代码?
3. 解决proto文件管理?
4. 添加复杂db及缓存接口
5. redis连接池? 
6. grpc trace
7. gin trace 
8. 负载均衡
9. 服务注册发现