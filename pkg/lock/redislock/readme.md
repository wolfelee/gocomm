使用方法

//全局设置SetGlobalPrefix 锁前缀
//先在全局设置redis客户端 SetRedisClient
//NewRedisLock 获取锁对象

//NewRedisLockV2 添加本地锁优化的分布式锁,抢锁的任务多的情况下优势明显性能提升明显
//NewRedisLock   20线程同时抢锁最长时间TestNewRedisLock2/g:18 (1.80s)
//NewRedisLockV2 20线程同时抢锁最长时间TestNewRedisLockV2/g:12 (0.09s)